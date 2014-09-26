package cloud

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/SyntropyDev/milli"
	"github.com/SyntropyDev/sqlutil"
	"github.com/SyntropyDev/val"
	"github.com/coopernurse/gorp"
	"github.com/huandu/facebook"
	"github.com/jteeuwen/go-pkg-rss"
	"github.com/jteeuwen/go-pkg-xmlx"

	"appengine"
	"appengine/urlfetch"
)

const (
	ObjectNameStory = "Story"
	TableNameStory  = "stories"
)

type FeedType string

const (
	FeedTypeTwitter  FeedType = "twitter"
	FeedTypeFacebook FeedType = "facebook"
	FeedTypeRSS      FeedType = "rss"
)

func (ft FeedType) GetStories(c appengine.Context, s gorp.SqlExecutor, m *Member, f *Feed) error {
	switch ft {
	case FeedTypeRSS:
		itemHandler := func(
			fe *feeder.Feed,
			ch *feeder.Channel,
			newitems []*feeder.Item) {

			for _, item := range newitems {
				story := NewStoryRSS(m, f, item)
				s.Insert(story)
			}
		}
		feed := feeder.New(1, true, nil, itemHandler)
		if err := feed.Fetch(f.Identifier, nil); err != nil {
			return err
		}
	case FeedTypeTwitter:
		v := url.Values{}
		v.Set("screen_name", f.Identifier)
		v.Set("include_rts", "false")

		anaconda.SetConsumerKey(os.Getenv("twitterApiKey"))
		anaconda.SetConsumerSecret(os.Getenv("twitterApiSecret"))
		api := anaconda.NewTwitterApi("", "")
		api.HttpClient = urlfetch.Client(c)

		tweets, err := api.GetUserTimeline(v)
		if err != nil {
			return err
		}

		for _, t := range tweets {
			story := NewStoryTwitter(m, f, t)
			s.Insert(story)
		}
	case FeedTypeFacebook:
		app := facebook.New(os.Getenv("facebookApiID"), os.Getenv("facebookAppSecret"))
		app.RedirectUri = "http://syntropy.io"
		session := app.Session(app.AppAccessToken())
		session.HttpClient = urlfetch.Client(c)

		route := fmt.Sprintf("/%s/posts", f.Identifier)
		result, err := session.Api(route, facebook.GET, nil)
		if err != nil {
			return err
		}

		posts := &FacebookPosts{}
		if err := result.Decode(posts); err != nil {
			return err
		}

		for _, post := range posts.Data {
			story := NewFacebookStory(m, f, post)
			s.Insert(story)
		}
	}
	return nil
}

func FeedTypes() []FeedType {
	return []FeedType{FeedTypeTwitter, FeedTypeFacebook, FeedTypeRSS}
}

type FacebookPosts struct {
	Data []*FacebookPost
}

type FacebookPost struct {
	CreatedAt string
	Id        string
	Picture   string
	Link      string
	Message   string
	Story     string
	Title     string
	Likes     FacebookLikes
}

type FacebookLikes struct {
	Data []interface{}
}

type Story struct {
	ID      int64  `json:"id"`
	Created int64  `json:"created" val:"nonzero"`
	Updated int64  `json:"updated" val:"nonzero"`
	Deleted bool   `json:"deleted" merge:"true"`
	Object  string `db:"-" json:"object"`

	MemberID       int64     `json:"memberId" val:"nonzero"`
	MemberName     string    `json:"memberName"`
	FeedID         int64     `json:"memberId" val:"nonzero"`
	FeedIdentifier string    `json:"feedIdentifier"`
	Timestamp      int64     `json:"timestamp"`
	Body           string    `json:"body"`
	SourceType     string    `json:"sourceType"`
	SourceURL      string    `json:"sourceUrl"`
	SourceID       string    `json:"sourceID"`
	Score          int64     `json:"score"`
	Latitude       float64   `json:"-"`
	Longitude      float64   `json:"-"`
	LinksRaw       string    `json:"-"`
	ImagesRaw      string    `json:"-"`
	HashtagsRaw    string    `json:"-"`
	Links          []string  `db:"-" json:"links"`
	Images         []string  `db:"-" json:"images"`
	Hashtags       []string  `db:"-" json:"hashTags"`
	Location       []float64 `db:"-" json:"location"`
}

func (story *Story) LinksSlice() []string {
	return sliceFromString(story.LinksRaw)
}

func (story *Story) ImagesSlice() []string {
	return sliceFromString(story.ImagesRaw)
}

func (story *Story) HashtagsSlice() []string {
	return sliceFromString(story.HashtagsRaw)
}

func (story *Story) LocationCoords() []float64 {
	return []float64{story.Latitude, story.Longitude}
}

func sliceFromString(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}

func (story *Story) CalculateScore(s gorp.SqlExecutor) error {
	f := &Feed{}
	if err := sqlutil.SelectOneRelation(s, TableNameFeed, story.FeedID, f); err != nil {
		return err
	}

	score := 0
	if len(story.ImagesSlice()) > 0 {
		score += 10
	}
	if len(story.LinksSlice()) > 0 {
		score += 2
	}
	if len(story.HashtagsSlice()) > 0 {
		score += 4
	}
	if story.Latitude != 0.0 {
		score += 10
	}

	switch FeedType(f.Type) {
	case FeedTypeFacebook:
		score += 3
	case FeedTypeTwitter:
		score += 3
	}

	story.Score += int64(score)
	return nil
}

func NewFacebookStory(member *Member, feed *Feed, post *FacebookPost) *Story {
	t, err := time.Parse(time.RFC3339Nano, post.CreatedAt)
	if err != nil {
		t = time.Now()
	}

	// try message, then story, then title
	text := post.Message
	if text == "" {
		text = post.Story
	}
	if text == "" {
		text = post.Title
	}

	sourceURL := fmt.Sprintf("https://www.facebook.com/%s/posts/%s", feed.Identifier, post.Id)

	return &Story{
		MemberID:       member.ID,
		MemberName:     member.Name,
		FeedID:         feed.ID,
		FeedIdentifier: feed.Identifier,
		Timestamp:      milli.Timestamp(t),
		Body:           text,
		SourceType:     string(FeedTypeFacebook),
		SourceURL:      sourceURL,
		SourceID:       post.Id,
		Latitude:       0.0,
		Longitude:      0.0,
		Score:          int64(len(post.Likes.Data)),
		LinksRaw:       post.Link,
		HashtagsRaw:    "",
		ImagesRaw:      post.Picture,
	}
}

func NewStoryTwitter(member *Member, feed *Feed, tweet anaconda.Tweet) *Story {
	hashtags := []string{}
	for _, hashtag := range tweet.Entities.Hashtags {
		hashtags = append(hashtags, hashtag.Text)
	}

	urls := []string{}
	for _, url := range tweet.Entities.Urls {
		urls = append(urls, url.Url)
	}

	t, err := tweet.CreatedAtTime()
	if err != nil {
		t = time.Now()
	}

	images := []string{}
	for _, media := range tweet.Entities.Media {
		if media.Type == "image" {
			images = append(images, media.Media_url)
		}
	}
	score := tweet.FavoriteCount + (2 * tweet.RetweetCount)
	sourceURL := fmt.Sprintf("http://twitter.com/%s/status/%s", tweet.User.ScreenName, tweet.IdStr)
	return &Story{
		MemberID:       member.ID,
		MemberName:     member.Name,
		FeedID:         feed.ID,
		FeedIdentifier: feed.Identifier,
		Timestamp:      milli.Timestamp(t),
		Body:           tweet.Text,
		SourceType:     string(FeedTypeTwitter),
		SourceURL:      sourceURL,
		SourceID:       tweet.IdStr,
		Latitude:       0.0,
		Longitude:      0.0,
		Score:          int64(score),
		LinksRaw:       strings.Join(urls, ","),
		HashtagsRaw:    strings.Join(hashtags, ","),
		ImagesRaw:      strings.Join(images, ","),
	}
}

func NewStoryRSS(member *Member, feed *Feed, item *feeder.Item) *Story {
	// parse pub date
	itemTime, err := item.ParsedPubDate()
	if err != nil {
		itemTime = time.Now()
	}
	// form links
	links := []string{}
	for _, link := range item.Links {
		links = append(links, link.Href)
	}

	isAtom := func() bool {
		return item.Id != ""
	}

	sourceID := ""
	body := ""
	if isAtom() {
		sourceID = item.Id
		body = item.Content.Text
	} else { // is RSS
		sourceID = *item.Guid
		// use description or title for body
		body = item.Description
		if body == "" {
			body = item.Title
		}
	}

	// parse html for images
	images := []string{}
	doc := xmlx.New()
	doc.LoadString(strings.ToLower(body), nil)
	imgNodes := doc.SelectNodesRecursive("", "img")
	for _, img := range imgNodes {
		images = append(images, img.As("", "src"))
	}

	return &Story{
		MemberID:       member.ID,
		MemberName:     member.Name,
		FeedID:         feed.ID,
		FeedIdentifier: feed.Identifier,
		Timestamp:      milli.Timestamp(itemTime),
		Body:           body,
		SourceType:     string(FeedTypeRSS),
		SourceURL:      "",
		SourceID:       sourceID,
		Latitude:       0.0,
		Longitude:      0.0,
		Score:          0.0,
		LinksRaw:       strings.Join(links, ","),
		ImagesRaw:      strings.Join(images, ","),
	}
}

func (story *Story) Validate() error {
	if valid, errMap := val.Struct(story); !valid {
		return ErrorFromMap(errMap)
	}
	return nil
}

func (story *Story) PreInsert(s gorp.SqlExecutor) error {
	story.Created = milli.Timestamp(time.Now())
	story.Updated = milli.Timestamp(time.Now())
	return story.Validate()
}

func (story *Story) PreUpdate(s gorp.SqlExecutor) error {
	story.Updated = milli.Timestamp(time.Now())
	return story.Validate()
}

func (story *Story) PostInsert(s gorp.SqlExecutor) error {
	m := &Member{}
	if err := sqlutil.SelectOneRelation(s, TableNameMember, story.MemberID, m); err != nil {
		return err
	}

	images := append(story.ImagesSlice(), m.ImagesSlice()...)
	m.SetImages(images)

	hashtags := append(story.HashtagsSlice(), m.HashtagsSlice()...)
	m.SetHashtags(hashtags)

	s.Update(m)
	return nil
}

func (story *Story) PostGet(s gorp.SqlExecutor) error {
	story.Images = story.ImagesSlice()
	story.Hashtags = story.HashtagsSlice()
	story.Location = story.LocationCoords()
	story.Object = ObjectNameMember
	return nil
}

// CrudResource interface

func (story *Story) TableName() string {
	return TableNameFeed
}

func (story *Story) TableId() int64 {
	return story.ID
}

func (story *Story) Delete() {
	story.Deleted = true
}

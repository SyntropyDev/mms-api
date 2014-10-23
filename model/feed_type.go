package model

import (
	"fmt"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/SyntropyDev/milli"
	"github.com/coopernurse/gorp"
	"github.com/huandu/facebook"
	"github.com/jteeuwen/go-pkg-rss"
)

type FeedType string

const (
	FeedTypeTwitter  FeedType = "twitter"
	FeedTypeFacebook FeedType = "facebook"
	FeedTypeRSS      FeedType = "rss"
)

func facebookSession() *facebook.Session {
	app := facebook.New(os.Getenv("facebookApiID"), os.Getenv("facebookAppSecret"))
	app.RedirectUri = "http://syntropy.io"
	return app.Session(app.AppAccessToken())
}

func twitterAPI() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(os.Getenv("twitterApiKey"))
	anaconda.SetConsumerSecret(os.Getenv("twitterApiSecret"))
	return anaconda.NewTwitterApi("", "")
}

func (ft FeedType) GetStories(s gorp.SqlExecutor, m *Member, f *Feed) error {
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
			fmt.Println(err)
		}
	case FeedTypeTwitter:
		v := url.Values{}
		v.Set("screen_name", f.Identifier)
		v.Set("include_rts", "false")

		anaconda.SetConsumerKey(os.Getenv("twitterApiKey"))
		anaconda.SetConsumerSecret(os.Getenv("twitterApiSecret"))
		api := anaconda.NewTwitterApi("", "")

		tweets, err := api.GetUserTimeline(v)
		if err != nil {
			fmt.Println(err)
		}

		for _, t := range tweets {
			story := NewStoryTwitter(m, f, t)
			if err := s.Insert(story); err == nil {
				fmt.Printf("Added Twitter story for %s.  Date: %s Score: %f\n", m.Name, milli.Time(story.Timestamp).String(), story.Score)
			} else {
				fmt.Printf("Failed to add Twitter story for %s.  Error: %s\n", m.Name, err)
			}
		}
	case FeedTypeFacebook:
		session := facebookSession()
		route := fmt.Sprintf("/%s/posts", f.Identifier)
		result, err := session.Api(route, facebook.GET, nil)
		if err != nil {
			fmt.Println(err)
		}

		posts := &FacebookPosts{}
		if err := result.Decode(posts); err != nil {
			fmt.Println(err)
		}

		for _, post := range posts.Data {
			story := NewFacebookStory(m, f, post)
			if story != nil {
				if err := s.Insert(story); err == nil {
					fmt.Printf("Added Facebook story for %s.  Date: %s  Score: %f\n", m.Name, milli.Time(story.Timestamp).String(), story.Score)
				} else {
					fmt.Printf("Failed to add Facebook story for %s.  Error: %s\n", m.Name, err)
				}
			}
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
	ObjectId  string
	Type      string
	Likes     FacebookLikes
}

type FacebookPhoto struct {
	CreatedTime string `json:"created_time"`
	Id          string `json:"id"`
	Images      []struct {
		Height int    `json:"height"`
		Source string `json:"source"`
		Width  int    `json:"width"`
	} `json:"images"`
}

type FacebookLikes struct {
	Data []interface{}
}

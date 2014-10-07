package model

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/SyntropyDev/httperr"
	"github.com/SyntropyDev/milli"
	"github.com/SyntropyDev/sqlutil"
	"github.com/SyntropyDev/val"
	"github.com/coopernurse/gorp"
	"github.com/huandu/facebook"
	"github.com/lann/squirrel"
)

const (
	ObjectNameFeed = "Feed"
	TableNameFeed  = "feeds"
)

type fbookUser struct {
	Cover struct {
		Source string
	}
}

type Feed struct {
	ID      int64  `json:"id"`
	Created int64  `json:"created" val:"nonzero"`
	Updated int64  `json:"updated" val:"nonzero"`
	Deleted bool   `json:"deleted" merge:"true"`
	Object  string `db:"-" json:"object"`

	MemberID      int64  `json:"memberId" val:"nonzero" merge:"true"`
	Type          string `json:"type" val:"in(twitter,facebook,rss)" merge:"true"`
	Identifier    string `json:"identifier" val:"nonzero" merge:"true"`
	LastRetrieved int64  `json:"-"`
}

func ListenToFeeds(s gorp.SqlExecutor) error {
	feeds := []*Feed{}
	query := squirrel.Select("*").From(TableNameFeed)
	if err := sqlutil.Select(s, query, &feeds); err != nil {
		return err
	}

	for _, feed := range feeds {
		if err := feed.UpdateStories(s); err != nil {
			return err
		}
	}
	return nil
}

func (f *Feed) UpdateStories(s gorp.SqlExecutor) error {
	m := &Member{}
	if err := sqlutil.SelectOneRelation(s, TableNameMember, f.MemberID, m); err != nil {
		return err
	}
	return FeedType(f.Type).GetStories(s, m, f)
}

func (f *Feed) Validate() error {
	if valid, errMap := val.Struct(f); !valid {
		return ErrorFromMap(errMap)
	}
	return nil
}

func (f *Feed) PreInsert(s gorp.SqlExecutor) error {
	f.Created = milli.Timestamp(time.Now())
	f.Updated = milli.Timestamp(time.Now())

	icon := ""
	switch FeedType(f.Type) {
	case FeedTypeFacebook:
		// get user from facebook api
		session := facebookSession()
		route := fmt.Sprintf("/%s", f.Identifier)
		result, err := session.Api(route, facebook.GET, nil)
		if err != nil {
			return httperr.New(http.StatusBadRequest, "invalid facebook id", err)
		}
		// decode response
		user := &fbookUser{}
		if err := result.Decode(user); err != nil {
			return err
		}
		icon = user.Cover.Source
	case FeedTypeTwitter:
		api := twitterAPI()
		user, err := api.GetUsersShow(f.Identifier, url.Values{})
		if err != nil {
			return err
		}
		icon = user.ProfileImageURL
	}

	if icon != "" {
		member := &Member{}
		if err := sqlutil.SelectOneRelation(s, TableNameMember, f.MemberID, member); err != nil {
			return err
		}
		member.Icon = icon

		if _, err := s.Update(member); err != nil {
			return err
		}
	}

	return f.Validate()
}

func (f *Feed) PreUpdate(s gorp.SqlExecutor) error {
	f.Updated = milli.Timestamp(time.Now())
	return f.Validate()
}

func (f *Feed) PostGet(s gorp.SqlExecutor) error {
	f.Object = ObjectNameFeed
	return nil
}

// CrudResource interface

func (f *Feed) TableName() string {
	return TableNameFeed
}

func (f *Feed) TableId() int64 {
	return f.ID
}

func (f *Feed) Delete() {
	f.Deleted = true
}

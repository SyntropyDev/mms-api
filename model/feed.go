package model

import (
	"time"

	"github.com/SyntropyDev/milli"
	"github.com/SyntropyDev/sqlutil"
	"github.com/SyntropyDev/val"
	"github.com/coopernurse/gorp"
)

const (
	ObjectNameFeed = "Feed"
	TableNameFeed  = "feeds"
)

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

func (u *Feed) PreInsert(s gorp.SqlExecutor) error {
	u.Created = milli.Timestamp(time.Now())
	u.Updated = milli.Timestamp(time.Now())
	return u.Validate()
}

func (u *Feed) PreUpdate(s gorp.SqlExecutor) error {
	u.Updated = milli.Timestamp(time.Now())
	return u.Validate()
}

func (m *Feed) PostGet(s gorp.SqlExecutor) error {
	m.Object = ObjectNameFeed
	return nil
}

// CrudResource interface

func (u *Feed) TableName() string {
	return TableNameFeed
}

func (u *Feed) TableId() int64 {
	return u.ID
}

func (u *Feed) Delete() {
	u.Deleted = true
}

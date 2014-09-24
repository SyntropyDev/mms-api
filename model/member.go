package model

import (
	"strings"
	"time"

	"github.com/SyntropyDev/milli"
	"github.com/SyntropyDev/val"
	"github.com/coopernurse/gorp"
)

const (
	ObjectNameMember = "Member"
	TableNameMember  = "members"
)

type Member struct {
	ID      int64  `json:"id"`
	Created int64  `json:"created" val:"nonzero"`
	Updated int64  `json:"updated" val:"nonzero"`
	Deleted bool   `json:"deleted" merge:"true"`
	Object  string `db:"-" json:"object"`

	Name        string   `json:"name" val:"nonzero" merge:"true"`
	Address     string   `json:"address" val:"nonzero" merge:"true"`
	Phone       string   `json:"phone" merge:"true"`
	Description string   `json:"description" merge:"true"`
	Icon        string   `json:"icon" merge:"true"`
	Website     string   `json:"website" merge:"true"`
	Latitude    float64  `json:"-" merge:"true"`
	Longitude   float64  `json:"-" merge:"true"`
	ImagesRaw   string   `json:"-"`
	HashtagsRaw string   `json:"-"`
	Images      []string `db:"-" json:"images"`
	Hashtags    []string `db:"-" json:"hashTags"`
	Location    []int64  `db:"-" json:"location"`
}

func (m *Member) ImagesSlice() []string {
	return strings.Split(m.ImagesRaw, ",")
}

func (m *Member) SetImages(s []string) {
	if len(s) > 5 {
		s = s[:5]
	}
	m.ImagesRaw = strings.Join(s, ",")
}

func (m *Member) HashtagsSlice() []string {
	return strings.Split(m.HashtagsRaw, ",")
}

func (m *Member) SetHashtags(s []string) {
	if len(s) > 5 {
		s = s[:5]
	}
	m.HashtagsRaw = strings.Join(s, ",")
}

// func (u *User) SetPassword(p *Password) {
// 	u.Password = p.String()
// 	u.PasswordHash = p.Hash()
// }

// func (u *User) HasPassword(password string) bool {
// 	bHash := []byte(u.PasswordHash)
// 	bPass := []byte(password)
// 	err := bcrypt.CompareHashAndPassword(bHash, bPass)
// 	return err == nil
// }

func (m *Member) Validate() error {
	if valid, errMap := val.Struct(m); !valid {
		return val.ErrorFromMap(errMap)
	}
	return nil
}

func (m *Member) PreInsert(s gorp.SqlExecutor) error {
	m.Created = milli.Timestamp(time.Now())
	m.Updated = milli.Timestamp(time.Now())
	return m.Validate()
}

func (m *Member) PreUpdate(s gorp.SqlExecutor) error {
	m.Updated = milli.Timestamp(time.Now())
	return m.Validate()
}

// CrudResource interface

func (m *Member) TableName() string {
	return TableNameMember
}

func (m *Member) TableId() int64 {
	return m.ID
}

func (m *Member) Delete() {
	m.Deleted = true
}

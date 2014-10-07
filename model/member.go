package model

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"code.google.com/p/go.crypto/bcrypt"

	"github.com/SyntropyDev/httperr"
	"github.com/SyntropyDev/milli"
	"github.com/SyntropyDev/sqlutil"
	"github.com/SyntropyDev/val"
	"github.com/coopernurse/gorp"
	"github.com/lann/squirrel"

	"github.com/mailgun/mailgun-go"
)

const (
	ObjectNameMember = "Member"
	TableNameMember  = "members"

	inviteEmailTemplate   = "You have been invited to Mobile Main Street!  Here is your temporary password: %s"
	passwordResetTemplate = "Here is your temporary password: %s"
)

type Member struct {
	ID      int64  `json:"id"`
	Created int64  `json:"created" val:"nonzero"`
	Updated int64  `json:"updated" val:"nonzero"`
	Deleted bool   `json:"deleted" merge:"true"`
	Object  string `db:"-" json:"object"`

	// auth user
	Email        string `json:"email" merge:"true"`
	Organizer    bool   `json:"-"`
	Token        string `db:"-" json:"token,omitempty"`
	Password     string `db:"-" json:"password,omitempty"`
	PasswordHash string `json:"-"`

	// member
	Name        string  `json:"name" val:"nonzero" merge:"true"`
	Address     string  `json:"address" merge:"true"`
	Phone       string  `json:"phone" merge:"true"`
	Description string  `json:"description" merge:"true"`
	Icon        string  `json:"icon" merge:"true"`
	Website     string  `json:"website" merge:"true"`
	Latitude    float64 `json:"-" merge:"true"`
	Longitude   float64 `json:"-" merge:"true"`
	ImagesRaw   string  `json:"-"`
	HashtagsRaw string  `json:"-"`

	CategoryIds []int64   `db:"-" json:"categoryIds"`
	Images      []string  `db:"-" json:"images"`
	Hashtags    []string  `db:"-" json:"hashTags"`
	Location    []float64 `db:"-" json:"location"`
}

func FindMember(s gorp.SqlExecutor, email string) (*Member, error) {
	query := squirrel.Select("*").From(TableNameMember).
		Where(squirrel.Eq{"email": email})
	members := []*Member{}
	if err := sqlutil.Select(s, query, &members); err != nil {
		return nil, err
	}

	err := fmt.Errorf("no user for email")
	if len(members) == 0 {
		return nil, err
	}

	return members[0], nil
}

func AuthenticateMember(s gorp.SqlExecutor, email, password string) (*Member, error) {
	err := fmt.Errorf("email / password invalid")
	respErr := httperr.New(http.StatusUnauthorized, err.Error(), err)

	member, err := FindMember(s, email)
	if err != nil || !member.HasPassword(password) {
		return nil, respErr
	}

	return member, nil
}

func (m *Member) ResetPassword() error {
	m.SetPassword(NewAutoPassword())

	// form invite email
	body := fmt.Sprintf(passwordResetTemplate, m.Password)
	recipient := fmt.Sprintf("%s <%s>", m.Name, m.Email)
	message := mailgun.NewMessage(
		"organizer@mobilemainst.com",
		"Mobile Main Street Password Reset",
		body, recipient)
	return sendEmail(message)
}

func (m *Member) Invite(email string) error {
	// set email and temp password
	m.Email = email
	m.SetPassword(NewAutoPassword())

	// form invite email
	body := fmt.Sprintf(inviteEmailTemplate, m.Password)
	recipient := fmt.Sprintf("%s <%s>", m.Name, m.Email)
	message := mailgun.NewMessage(
		"organizer@mobilemainst.com",
		"Mobile Main Street Invite",
		body, recipient)
	return sendEmail(message)
}

func (m *Member) ImagesSlice() []string {
	return sliceFromString(m.ImagesRaw)
}

func (m *Member) SetImages(s []string) {
	if len(s) > 5 {
		s = s[:5]
	}
	m.ImagesRaw = strings.Join(s, ",")
}

func (m *Member) HashtagsSlice() []string {
	return sliceFromString(m.HashtagsRaw)
}

func (m *Member) SetHashtags(s []string) {
	if len(s) > 5 {
		s = s[:5]
	}
	m.HashtagsRaw = strings.Join(s, ",")
}

func (m *Member) LocationCoords() []float64 {
	return []float64{m.Latitude, m.Longitude}
}

func (m *Member) SetPassword(p *Password) {
	m.Password = p.String()
	m.PasswordHash = p.Hash()
}

func (m *Member) HasPassword(password string) bool {
	bHash := []byte(m.PasswordHash)
	bPass := []byte(password)
	err := bcrypt.CompareHashAndPassword(bHash, bPass)
	return err == nil
}

func (m *Member) Validate() error {
	if valid, errMap := val.Struct(m); !valid {
		return ErrorFromMap(errMap)
	}
	return nil
}

func (m *Member) PreInsert(s gorp.SqlExecutor) error {
	m.Created = milli.Timestamp(time.Now())
	m.Updated = milli.Timestamp(time.Now())
	if err := m.updateCategoires(s); err != nil {
		return err
	}
	return m.Validate()
}

func (m *Member) PreUpdate(s gorp.SqlExecutor) error {
	m.Updated = milli.Timestamp(time.Now())
	if err := m.updateCategoires(s); err != nil {
		return err
	}
	return m.Validate()
}

func (m *Member) PostGet(s gorp.SqlExecutor) error {
	m.Images = m.ImagesSlice()
	m.Hashtags = m.HashtagsSlice()
	m.Location = m.LocationCoords()
	m.Object = ObjectNameMember

	catIds := []int64{}
	catMems := []*CategoryMember{}
	query := squirrel.Select("*").From(TableNameCategoryMember).
		Where(squirrel.Eq{"memberID": m.ID})
	sqlutil.Select(s, query, &catMems)
	for _, catMem := range catMems {
		catIds = append(catIds, catMem.CategoryID)
	}
	m.CategoryIds = catIds

	return nil
}

func (m *Member) updateCategoires(s gorp.SqlExecutor) error {
	// delete existing categories
	format := "delete from " + TableNameCategoryMember + " where memberid = ?"
	if _, err := s.Exec(format, m.ID); err != nil {
		return err
	}

	// create new categories
	for _, catId := range m.CategoryIds {
		catMem := NewCategoryMember(catId, m.ID)
		if err := s.Insert(catMem); err != nil {
			return err
		}
	}
	return nil
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

// util

func sendEmail(message *mailgun.Message) error {
	gun := mailgun.NewMailgun(
		os.Getenv("mailgunDomain"),
		os.Getenv("mailgunPublicApiKey"),
		os.Getenv("mailgunPrivateApiKey"))
	_, _, err := gun.Send(message)
	return err
}

package model

import (
	"fmt"
	"time"

	"github.com/SyntropyDev/milli"
	"github.com/SyntropyDev/sqlutil"
	"github.com/SyntropyDev/val"
	"github.com/coopernurse/gorp"
	"github.com/dchest/uniuri"
	"github.com/lann/squirrel"
)

const (
	ModelNameToken = "Token"
	TableNameToken = "tokens"
)

type Token struct {
	ID        int64  `json:"id"`
	AccountID int64  `json:"accountId" val:"nonzero"`
	Created   int64  `json:"created" val:"nonzero"`
	Updated   int64  `json:"updated" val:"nonzero"`
	Deleted   bool   `json:"deleted" merge:"true"`
	ModelName string `db:"-" json:"modelName"`

	MemberID   int64  `json:"memberId" val:"nonzero"`
	Value      string `json:"value" val:"nonzero"`
	Expiration int64  `json:"expirationTimestamp" val:"nonzero"`
}

func ValidateToken(s gorp.SqlExecutor, memberID int64, token string) error {
	query := squirrel.Select("*").From(TableNameToken).
		Where(squirrel.Eq{"memberID": memberID, "value": token})
	tokens := []*Token{}
	sqlutil.Select(s, query, &tokens)
	if len(tokens) > 0 {
		return fmt.Errorf("token not found")
	}
	return nil
}

func (t *Token) IsExpired() bool {
	exp := milli.Time(t.Expiration)
	return time.Now().After(exp)
}

func (t *Token) Validate() error {
	if valid, errMap := val.Struct(t); !valid {
		return val.ErrorFromMap(errMap)
	}
	return nil
}

func (t *Token) PreInsert(s gorp.SqlExecutor) error {
	t.Created = milli.Timestamp(time.Now())
	t.Updated = milli.Timestamp(time.Now())
	t.Value = uniuri.NewLen(30)
	ex := time.Now().AddDate(0, 0, 14)
	t.Expiration = milli.Timestamp(ex)
	return t.Validate()
}

func (t *Token) PreUpdate(s gorp.SqlExecutor) error {
	t.Updated = milli.Timestamp(time.Now())
	return t.Validate()
}

func (t *Token) PostGet(s gorp.SqlExecutor) error {
	t.ModelName = ModelNameToken
	return nil
}

package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/SyntropyDev/milli"
	"github.com/SyntropyDev/val"
	"github.com/coopernurse/gorp"
)

const (
	ObjectNameCategory = "Category"
	TableNameCategory  = "categories"
)

type Category struct {
	ID      int64  `json:"id"`
	Created int64  `json:"created" val:"nonzero"`
	Updated int64  `json:"updated" val:"nonzero"`
	Deleted bool   `json:"deleted" merge:"true"`
	Object  string `db:"-" json:"object"`

	Name string `json:"name" val:"nonzero" merge:"true"`
}

func (c *Category) Validate() error {
	if valid, errMap := val.Struct(c); !valid {
		return ErrorFromMap(errMap)
	}
	return nil
}

func (c *Category) PreInsert(s gorp.SqlExecutor) error {
	c.Created = milli.Timestamp(time.Now())
	c.Updated = milli.Timestamp(time.Now())
	return c.Validate()
}

func (c *Category) PreUpdate(s gorp.SqlExecutor) error {
	c.Updated = milli.Timestamp(time.Now())
	return c.Validate()
}

// CrudResource interface

func (c *Category) TableName() string {
	return TableNameCategory
}

func (c *Category) TableId() int64 {
	return c.ID
}

func (c *Category) Delete() {
	c.Deleted = true
}

func ErrorFromMap(errMap map[string]error) error {
	errs := []string{}
	for key, err := range errMap {
		nErr := fmt.Sprintf("%s - %s", key, err)
		errs = append(errs, nErr)
	}
	return errors.New(strings.Join(errs, ","))
}

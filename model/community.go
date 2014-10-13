package model

import (
	"errors"
	"net/http"
	"time"

	"github.com/SyntropyDev/httperr"
	"github.com/SyntropyDev/milli"
	"github.com/SyntropyDev/val"
	"github.com/coopernurse/gorp"
)

const (
	ObjectNameCommunity = "Community"
	TableNameCommunity  = "communities"
)

type Community struct {
	ID      int64  `json:"id"`
	Created int64  `json:"created" val:"nonzero"`
	Updated int64  `json:"updated" val:"nonzero"`
	Deleted bool   `json:"deleted" merge:"true"`
	Object  string `db:"-" json:"object"`

	Name        string    `json:"name" val:"nonzero" merge:"true"`
	Latitude    float64   `json:"-" val:"lat"`
	Longitude   float64   `json:"-" val:"lon"`
	Description string    `json:"description" val:"nonzero"`
	Location    []float64 `db:"-" json:"location" merge:"true"`
}

func (c *Community) Validate() error {
	if err := c.updateLatLng(); err != nil {
		return err
	}
	if valid, errMap := val.Struct(c); !valid {
		return ErrorFromMap(errMap)
	}
	return nil
}

func (c *Community) PreInsert(s gorp.SqlExecutor) error {
	c.Created = milli.Timestamp(time.Now())
	c.Updated = milli.Timestamp(time.Now())
	return c.Validate()
}

func (c *Community) PreUpdate(s gorp.SqlExecutor) error {
	c.Updated = milli.Timestamp(time.Now())
	return c.Validate()
}

func (c *Community) PostGet(s gorp.SqlExecutor) error {
	c.Object = ObjectNameCommunity
	c.Location = []float64{c.Latitude, c.Longitude}
	return nil
}

func (c *Community) updateLatLng() error {
	if len(c.Location) != 2 {
		err := errors.New("Location should be [lat, lon]")
		return httperr.New(http.StatusBadRequest, err.Error(), err)
	}
	c.Latitude = c.Location[0]
	c.Longitude = c.Location[1]
	return nil
}

// CrudResource interface

func (c *Community) TableName() string {
	return TableNameCommunity
}

func (c *Community) TableId() int64 {
	return c.ID
}

func (c *Community) Delete() {
	c.Deleted = true
}

package milli_test

import (
	"testing"
	"time"

	"github.com/SyntropyDev/milli"
)

func TestTime(t *testing.T) {
	now := time.Now()
	timestamp := milli.Timestamp(now)
	final := milli.Time(timestamp)
	if now.Unix() != final.Unix() {
		t.Error("incorrect conversion")
	}
}

func TestDuration(t *testing.T) {
	timestamp := milli.DurationMilli(time.Hour)
	final := milli.Duration(timestamp)
	if time.Hour != final {
		t.Error("incorrect conversion")
	}
}

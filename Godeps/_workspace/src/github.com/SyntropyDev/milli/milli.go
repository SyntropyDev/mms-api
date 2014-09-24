/*
Package milli provide convience methods for converting between milliseconds, time.Duration, and time.Time.
*/
package milli

import "time"

// Duration converts the millisecond duration into a duration
func Duration(timestamp int64) time.Duration {
	return time.Duration(timestamp * int64(time.Millisecond))
}

// DurationMilli converts d into a millisecond duration
func DurationMilli(d time.Duration) int64 {
	return int64(d) / int64(time.Millisecond)
}

// Time converts the millisecond timestamp into a time
func Time(timestamp int64) time.Time {
	nsecs := timestamp * int64(time.Millisecond)
	return time.Unix(0, nsecs)
}

// Timestamp converts t into a millisecond timestamp
func Timestamp(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

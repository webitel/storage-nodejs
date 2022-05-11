package model

import "time"

const msInSecond int64 = 1e3
const nsInMillisecond int64 = 1e6

func NewBool(b bool) *bool       { return &b }
func NewInt(n int) *int          { return &n }
func NewInt64(n int64) *int64    { return &n }
func NewString(s string) *string { return &s }

// GetMillis is a convience method to get milliseconds since epoch.
func GetMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// UnixToMS Converts Unix Epoch from milliseconds to time.Time
func UnixToTime(ms int64) time.Time {
	return time.Unix(ms/msInSecond, (ms%msInSecond)*nsInMillisecond)
}

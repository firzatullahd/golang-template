package timeutils

import "time"

type ITime interface {
	Now() time.Time
}

type Time struct{}

func (Time) Now() time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	return time.Now().In(loc)
}

func (Time) UntilMidnight() time.Duration {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	return time.Duration(time.Until(now.Truncate(24 * time.Hour).Add(24 * time.Hour)).Seconds())
}

package models

import "time"

type Session struct {
	Id          string     `pg:"type:uuid" pg:",pk" json:"id"`
	AnonymousID string     `pg:"type:uuid" json:"anonymous_id"`
	Timestamp   time.Time  `json:"timestamp"`
	Referrer    string     `json:"referrer"`
	ScreenWidth int        `json:"screen_width"`
	IP          string     `json:"ip"`
	UserAgent   string     `json:"user_agent"`
	Country     string     `json:"country"`
	CountryCode string     `json:"country_code"`
	OS          string     `json:"os"`
	Browser     string     `json:"browser"`
	Events      []Event    `json:"events"`
	PageViews   []PageView `json:"pageviews"`
}

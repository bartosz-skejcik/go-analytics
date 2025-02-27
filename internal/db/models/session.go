package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	Id          uint       `gorm:"primaryKey" json:"id"`
	AnonymousID string     `json:"anonymous_id" gorm:"default:''"`
	Timestamp   time.Time  `json:"timestamp"`
	Referrer    string     `json:"referrer" gorm:"default:''"`
	ScreenWidth int        `json:"screen_width" gorm:"default:0"`
	IP          string     `json:"ip" gorm:"default:''"`
	UserAgent   string     `json:"user_agent" gorm:"default:''"`
	Country     string     `json:"country" gorm:"default:''"`
	CountryCode string     `json:"country_code" gorm:"default:''"`
	OS          string     `json:"os" gorm:"default:''"`
	Browser     string     `json:"browser" gorm:"default:''"`
	Events      []Event    `gorm:"foreignKey:session_id;references:ID" json:"events"`
	PageViews   []PageView `gorm:"foreignKey:session_id;references:ID" json:"pageviews"`
}

type SessionCreateDto struct {
	AnonymousID string              `json:"anonymous_id"`
	Timestamp   time.Time           `json:"timestamp"`
	Referrer    string              `json:"referrer"`
	ScreenWidth int                 `json:"screen_width"`
	IP          string              `json:"ip"`
	UserAgent   string              `json:"user_agent"`
	Country     string              `json:"country"`
	CountryCode string              `json:"country_code"`
	OS          string              `json:"os"`
	Browser     string              `json:"browser"`
	Events      []EventCreateDto    `json:"events"`
	PageViews   []PageViewCreateDto `json:"pageviews"`
}

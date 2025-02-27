package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Id        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
	Data      []byte    `json:"data" gorm:"type:jsonb"` // Changed to jsonb type
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	SessionId *uint     `json:"session_id" gorm:"index"`
}

type EventCreateDto struct {
	Name      string      `json:"name"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
	IP        string      `json:"ip"`
	UserAgent string      `json:"user_agent"`
}

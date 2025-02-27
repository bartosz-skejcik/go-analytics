package models

import "gorm.io/gorm"

type PageView struct {
	gorm.Model
	Id        uint   `gorm:"primaryKey" json:"id"`
	Route     string `json:"route" gorm:"unique"`
	Data      []byte `json:"data" gorm:"type:jsonb"` // Changed to jsonb type
	SessionId *uint  `json:"session_id" gorm:"index"`
}

type PageViewCreateDto struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

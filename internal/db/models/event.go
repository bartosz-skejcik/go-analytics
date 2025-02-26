package models

import "time"

type Event struct {
	Id        string                 `pg:"type:uuid" pg:",pk" json:"id"`
	Name      string                 `db:"name" json:"name"`
	Timestamp time.Time              `db:"timestamp" json:"timestamp"`
	Data      map[string]interface{} `db:"data" json:"data"`
	IP        string                 `db:"ip" json:"ip"`
	UserAgent string                 `db:"user_agent" json:"user_agent"`
}

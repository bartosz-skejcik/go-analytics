package service

import "github.com/bartosz-skejcik/go-analytics/internal/db"

type Service struct {
	db *db.Database
}

func New(db *db.Database) *Service {
	return &Service{
		db: db,
	}
}

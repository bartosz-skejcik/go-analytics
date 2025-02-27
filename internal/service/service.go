package service

import (
	"github.com/bartosz-skejcik/go-analytics/internal/db"
	"github.com/bartosz-skejcik/go-analytics/internal/db/models"
)

type Service struct {
	db      *db.Database
	session *models.Session
}

func New(db *db.Database) *Service {
	return &Service{
		db: db,
	}
}

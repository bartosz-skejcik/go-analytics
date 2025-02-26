package handlers

import (
	"github.com/bartosz-skejcik/go-analytics/internal/service"
)

type Handlers struct {
	service *service.Service
}

func New(s *service.Service) *Handlers {
	return &Handlers{
		service: s,
	}
}

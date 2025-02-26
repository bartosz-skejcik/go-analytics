package server

import (
	"os"

	"github.com/bartosz-skejcik/go-analytics/internal/db"
	"github.com/bartosz-skejcik/go-analytics/internal/server/router"
	"github.com/bartosz-skejcik/go-analytics/internal/service"
	"github.com/gin-gonic/gin"
)

type Server struct {
	r  *router.Router
	db *db.Database
}

func New(db *db.Database) *Server {
	return &Server{
		db: db,
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS_ORIGIN"))
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		c.Writer.Header().Set("Access-Control-Max-Age", "300")
	}
}

func (s *Server) Run() error {
	service := service.New(s.db)

	s.r = router.New(corsMiddleware, service)

	s.r.RegisterRoutes()
	return s.r.App.Run(":42068")
}

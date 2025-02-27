package router

import (
	"github.com/bartosz-skejcik/go-analytics/internal/server/handlers"
	"github.com/bartosz-skejcik/go-analytics/internal/service"
	"github.com/gin-gonic/gin"
)

type Router struct {
	*handlers.Handlers
	middleware func() gin.HandlerFunc
	service    *service.Service
	App        *gin.Engine
	Router     *gin.RouterGroup
}

func New(mw func() gin.HandlerFunc, s *service.Service) *Router {
	h := handlers.New(s)
	app := gin.Default()
	r := app.Group("/api")
	r.Use(mw())

	return &Router{
		Handlers:   h,
		middleware: mw,
		service:    s,
		App:        app,
		Router:     r,
	}
}

func (r *Router) RegisterRoutes() {
	r.Router.GET("/", r.HandleHome)
	r.Router.GET("/sessions", r.GetSessions)
	r.Router.POST("/sessions", r.CreateSession)
}

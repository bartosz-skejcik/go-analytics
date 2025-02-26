package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
)

type api struct {
	Version string    `json:"version"`
	Time    time.Time `json:"time"`
}

func (h *Handlers) HandleHome(c *gin.Context) {
	apiInfo := api{
		Version: "v0.0.1",
		Time:    time.Now(),
	}

	h.SuccessResponse(c, "", apiInfo)
}

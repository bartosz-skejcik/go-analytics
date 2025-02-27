package handlers

import (
	"github.com/bartosz-skejcik/go-analytics/internal/db/models"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) GetSessions(c *gin.Context) {
	sessions, err := h.service.GetAllSessions()
	if err != nil {
		h.ErrorResponse(c, err.Error())
		return
	}

	h.SuccessResponse(c, "Successfuly fetched all sessions", sessions)
}

func (h *Handlers) CreateSession(c *gin.Context) {
	var body models.SessionCreateDto
	err := c.Bind(&body)
	if err != nil {
		h.ErrorResponse(c, err.Error())
	}

}

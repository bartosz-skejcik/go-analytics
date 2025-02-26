package handlers

import "github.com/gin-gonic/gin"

func (h *Handlers) HandleGetSessions(c *gin.Context) {
	sessions, err := h.service.GetAllSessions()
	if err != nil {
		h.ErrorResponse(c, err.Error())
		return
	}

	h.SuccessResponse(c, "Successfuly fetched all sessions", sessions)
}

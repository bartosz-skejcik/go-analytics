package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Error   bool        `json:"error"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func (h *Handlers) SuccessResponse(c *gin.Context, message string, data interface{}) {
	resp := Response{
		Error:   false,
		Data:    data,
		Message: message,
	}

	responseJson, err := json.Marshal(resp)
	if err != nil {
		h.ErrorResponse(c, err.Error())
		return
	}

	c.Data(http.StatusOK, "application/json", responseJson)
}

func (h *Handlers) ErrorResponse(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   true,
		"message": message,
		"data":    nil,
	})
}

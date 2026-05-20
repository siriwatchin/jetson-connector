package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RawDataHandler struct{}

func (h *RawDataHandler) Create(c *gin.Context) {
	var payload json.RawMessage
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

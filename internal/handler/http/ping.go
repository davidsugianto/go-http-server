package http

import (
	"github.com/davidsugianto/go-pkgs/response"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Ping(c *gin.Context) {
	response.GinSuccess(c, gin.H{"status": "ok"})
}

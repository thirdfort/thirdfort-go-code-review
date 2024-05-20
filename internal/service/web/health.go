package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *WebService) setHealth() {
	s.Gin.GET("/healthz", s.getHealth)
}

func (s *WebService) getHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (obj *serviceSt) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "I am healthy!",
	})
}

package preflight

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Preflight function returns status OK for all preflight requests and NotFound for others
func Preflight(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusOK)
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

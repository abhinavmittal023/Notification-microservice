package preflight

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//Preflight function responds to the preflight options request
func Preflight(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, origin")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
	c.JSON(http.StatusOK, struct{}{})
}

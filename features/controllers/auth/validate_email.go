package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//ValidateEmail Controller for /signup route
func ValidateEmail(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{"token": c.Param("token")})
}

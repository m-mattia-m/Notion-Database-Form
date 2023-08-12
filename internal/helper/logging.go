package helper

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type httpError struct {
	Message string `json:"message"`
}

func SetNotFoundResponse(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusNotFound, httpError{
		Message: message,
	})
	return
}

func SetBadRequestResponse(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, httpError{
		Message: message,
	})
	return
}

func SetUnauthorizedResponse(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, httpError{
		Message: message,
	})
	return
}

func SetForbiddenResponse(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusForbidden, httpError{
		Message: message,
	})
	return
}

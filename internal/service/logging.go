package service

import (
	"Notion-Forms/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type httpError struct {
	Message string `json:"message"`
	LogId   string `json:"log_id"`
}

func (svc Clients) SetAbortResponse(c *gin.Context, targetServiceName string, method string, message string, err error) string {
	logId := uuid.New().String()
	svc.logger.Error(targetServiceName, method, logging.Message{
		Description: message,
		Detail:      err,
		LogId:       logId,
	})
	c.AbortWithStatusJSON(http.StatusInternalServerError, httpError{
		Message: message,
		LogId:   logId,
	})
	return logId
}

func (svc Clients) SetAbortWithoutResponse(targetServiceName string, method string, message string, err error) string {
	logId := uuid.New().String()
	svc.logger.Error(targetServiceName, method, logging.Message{
		Description: message,
		Detail:      err,
		LogId:       logId,
	})
	return logId
}

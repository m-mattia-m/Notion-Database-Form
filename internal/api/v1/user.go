package v1

import (
	"Notion-Forms/internal/helper"
	"Notion-Forms/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthenticateNotion(c *gin.Context) {
	svc, config, oidcUser, err := getConfigAndService(c)
	if err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("failed to get config and service, contact the administrator"))
		return
	}

	var code model.OAuthCodeRequest
	err = c.BindJSON(&code)
	if err != nil {
		svc.SetAbortResponse(c, "c", "BindJSON", fmt.Sprintf("failed to bind oauth-code-request to object"), err)
		return
	}

	err = svc.ConnectIamUserWithNotionUser(oidcUser.Sub, config.FrontendHost, code.Code)
	if err != nil {
		svc.SetAbortResponse(c, "svc", "ConnectIamUserWithNotionUser", fmt.Sprintf("failed to authorzise notion-user and save their id"), err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

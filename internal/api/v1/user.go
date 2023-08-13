package v1

import (
	"Notion-Forms/internal/helper"
	"Notion-Forms/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthenticateNotion 		godoc
// @title           		AuthenticateNotion
// @description     		Authenticate in Notion with the OAuth code and store the secret together with the IAM and Notion users
// @Tags 					Notion-Authenticate
// @Router  				/notion/authenticate [post]
// @Accept 					json
// @Produce					json
// @Security				Bearer
// @Param					OAuthCodeRequest	body 		model.OAuthCodeRequest 	true 	"RecordRequest"
// @Success      			200  				{object} 	model.HttpError
// @Failure      			400  				{object} 	model.HttpError
// @Failure      			404  				{object} 	model.HttpError
// @Failure      			500  				{object} 	model.HttpError
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

package v1

import (
	"Notion-Forms/internal/helper"
	"Notion-Forms/internal/model"
	"Notion-Forms/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ApiConfig struct {
	RunMode         string
	FrontendHost    string
	DefaultPageSize int
	MaxPageSize     int
	Port            int
	Host            string
	Domain          string
	Schemas         string
	OidcAuthority   string
	OidcClientId    string
}

func SetService(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("service", svc)
		//c.Next()
	}
}

func SetApiConfig(apiConfig ApiConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("config", apiConfig)
		//c.Next()
	}
}

func getConfigAndService(c *gin.Context) (*service.Clients, *ApiConfig, *model.OidcUser, error) {
	config, found := c.Get("config")
	if !found {
		return nil, nil, nil, fmt.Errorf("failed to get config from request")
	}
	configObject := config.(ApiConfig)

	svc, found := c.Get("service")
	if !found {
		return nil, nil, nil, fmt.Errorf("failed to get services from request")
	}
	svcObject := svc.(service.Clients)

	user, err := helper.GetUser(c, model.OidcConfig{
		AppEnv:        configObject.RunMode,
		OidcAuthority: configObject.OidcAuthority,
		OidcClientId:  configObject.OidcClientId,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	return &svcObject, &configObject, &user, nil
}

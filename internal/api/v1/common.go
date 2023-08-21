package v1

import (
	"Notion-Forms/internal/helper"
	"Notion-Forms/internal/model"
	"Notion-Forms/internal/service"
	"Notion-Forms/pkg/notion"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type ApiConfig struct {
	RunMode            string
	FrontendHost       string
	DefaultPageSize    int
	MaxPageSize        int
	Port               int
	Host               string
	Domain             string
	Schemas            string
	OidcAuthority      string
	OidcClientId       string
	NotionClientId     string
	NotionClientSecret string
	NotionRedirectUri  string
}

func SetService(svc service.Service, cfg ApiConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		notionClient, err := notion.New(cfg.NotionClientSecret, cfg.NotionClientSecret, cfg.NotionClientId, cfg.NotionRedirectUri)
		if err != nil {
			svc.SetAbortResponse(c, "notion", "New", fmt.Sprintf("faield to create notion client"), err)
			return
		}
		if !strings.Contains(c.FullPath(), "/notion/authenticate") {
			oidcUser, err := helper.GetUser(c, model.OidcConfig{
				AppEnv:        cfg.RunMode,
				OidcAuthority: cfg.OidcAuthority,
				OidcClientId:  cfg.OidcClientId,
			})
			if err != nil {
				svc.SetAbortResponse(c, "helper", "GetUser", fmt.Sprintf("failed to get user from context"), err)
				return
			}

			iamUserData, err := svc.GetOwnUser(oidcUser)
			if err != nil {
				svc.SetAbortResponse(c, "svc", "GetOwnUser", fmt.Sprintf("faield to get own user"), err)
				return
			}

			notionClient, err = notion.New(iamUserData.NotionCredentials.AccessToken, iamUserData.NotionCredentials.AccessToken, cfg.NotionClientId, cfg.NotionRedirectUri)
			if err != nil {
				svc.SetAbortResponse(c, "notion", "New", fmt.Sprintf("faield to create notion client"), err)
				return
			}
		}

		svc = svc.SetNotionClient(notionClient)
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

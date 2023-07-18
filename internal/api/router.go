package api

import (
	apiV1 "Notion-Forms/internal/api/v1"
	"Notion-Forms/internal/service"
	"github.com/gin-gonic/gin"
)

type Client struct {
	Api *apiV1.ApiClient
}

func New(svc service.Service, apiConfig apiV1.ApiConfig) *Client {
	return &Client{
		Api: &apiV1.ApiClient{
			ApiConfig: apiConfig,
			Service:   &svc,
		},
	}
}

func Router(client *Client) error {
	r := gin.New()
	v1 := r.Group("/api/v1")
	{
		notionGroup := v1.Group("/notion")
		{
			notionPages := notionGroup.Group("/pages")
			{
				notionPages.GET("/", client.Api.ListPages)
				notionPages.GET("/:type/:id", client.Api.GetPage)
			}
			notionRecords := notionGroup.Group("/records")
			{
				notionRecords.POST("/:databaseId", client.Api.CreateRecord)
				notionRecords.GET("/:databaseId/options", client.Api.ListAllSelectOptions)
				notionRecords.GET("/:databaseId/options/:select", client.Api.ListSelectOptions)
			}
		}
	}
	err := r.Run(":3000")
	return err
}

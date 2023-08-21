package api

import (
	apiV1 "Notion-Forms/internal/api/v1"
	"Notion-Forms/internal/service"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

// TODO: consider if the individual routing groups should not be called directly under `apiV1`, but e.g. `apiV1.notion.List()`

func Router(svc service.Service, apiConfig apiV1.ApiConfig) error {
	r := gin.Default()
	r.Use(sentrygin.New(sentrygin.Options{}))
	r.Use(apiV1.SetService(svc, apiConfig))
	r.Use(apiV1.SetApiConfig(apiConfig))
	v1 := r.Group("/api/v1")
	{
		notionGroup := v1.Group("/notion", Authenticate)
		{
			notionAuth := notionGroup.Group("/authenticate")
			{
				notionAuth.POST("/", apiV1.AuthenticateNotion)
			}
			notionPages := notionGroup.Group("/page")
			{
				notionPages.GET("/", apiV1.ListPages)
				notionPages.GET("/:pageId", apiV1.GetPage)
				//notionPages.DELETE("/cache", apiV1.DeletePageListFromCache)
				//notionPages.DELETE("/:pageId/cache", apiV1.DeletePageFromCache)
			}
			notionDatabases := notionGroup.Group("/database")
			{
				notionDatabases.GET("/", apiV1.ListDatabases)
				notionDatabases.GET("/:databaseId", apiV1.GetDatabase)
				notionDatabases.POST("/:databaseId", apiV1.CreateRecord)
				notionDatabases.GET("/:databaseId/properties", apiV1.GetDatabasePropertiesById)
				notionDatabases.GET("/:databaseId/properties/options", apiV1.ListAllSelectOptions)
				notionDatabases.GET("/:databaseId/properties/options/:notionSelectId", apiV1.ListSelectOptions)
				//notionDatabases.DELETE("/cache", apiV1.DeleteDatabaseListFromCache)
				//notionDatabases.DELETE("/:id/cache", apiV1.DeleteDatabaseFromCache)
			}
		}
		storageGroup := v1.Group("/storage", Authenticate)
		{
			storageGroup.POST("/authenticate/google", apiV1.AuthenticateGoogleDrive)
			storageGroup.POST("/provider/:databaseId", apiV1.SetStorageProvider)
			storageGroup.POST("/location/:databaseId", apiV1.SetBaseStorageLocation)
			storageGroup.POST("/upload/:databaseId", apiV1.UploadFile)
		}
	}
	err := r.Run(":3001")
	return err
}

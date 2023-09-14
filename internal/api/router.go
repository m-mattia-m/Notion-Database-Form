package api

import (
	"Notion-Forms/docs"
	apiV1 "Notion-Forms/internal/api/v1"
	"Notion-Forms/internal/service"
	"fmt"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

// TODO: consider if the individual routing groups should not be called directly under `apiV1`, but e.g. `apiV1.notion.List()`

func Router(svc service.Service, apiConfig apiV1.ApiConfig) error {
	if apiConfig.RunMode != "PROD" {
		fmt.Println(fmt.Sprintf("[HOST]: %s", apiConfig.Host))
		fmt.Println(fmt.Sprintf("[Swagger]: %s/api/v1/swagger/index.html", apiConfig.Host))
	}

	docs.SwaggerInfo.Schemes = []string{apiConfig.Schemas}
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", apiConfig.Domain, apiConfig.Port)
	docs.SwaggerInfo.BasePath = "/api/v1"

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // , "HEAD", "CONNECT", "TRACE", "PATCH"
		AllowHeaders: []string{"Authorization", "Origin", "Content-Type", "Accept", "Referer"},
		//ExposeHeaders:    []string{"Content-Length"},
		//AllowCredentials: true,
	}))

	r.Use(sentrygin.New(sentrygin.Options{}))
	r.Use(apiV1.SetService(svc, apiConfig))
	r.Use(apiV1.SetApiConfig(apiConfig))

	//r.GET("/api/v1/notion", func(c *gin.Context) {
	//	c.JSON(http.StatusOK, "👋 OK")
	//})
	//r.GET("/api/v1/notion/database", func(c *gin.Context) {
	//	c.JSON(http.StatusOK, "👋 OK")
	//})

	r.GET("/readiness", func(c *gin.Context) {
		// TODO: Check if clients are available and if db available
		c.JSON(http.StatusOK, "👋 OK")
	})
	r.GET("/liveliness", func(c *gin.Context) {
		c.JSON(http.StatusOK, "👋 OK")
	})

	v1 := r.Group("/api/v1")
	{
		notionGroup := v1.Group("/notion", Authenticate)
		{
			notionAuth := notionGroup.Group("/authenticate")
			{
				notionAuth.POST("", apiV1.AuthenticateNotion)
			}
			notionPages := notionGroup.Group("/page")
			{
				notionPages.GET("", apiV1.ListPages)
				notionPages.GET("/:pageId", apiV1.GetPage)
				//notionPages.DELETE("/cache", apiV1.DeletePageListFromCache)
				//notionPages.DELETE("/:pageId/cache", apiV1.DeletePageFromCache)
			}
			notionDatabases := notionGroup.Group("/database")
			{
				notionDatabases.GET("", apiV1.ListDatabases)
				notionDatabases.GET("/:databaseId", apiV1.GetDatabase)
				notionDatabases.POST("/:databaseId", apiV1.CreateRecord)
				notionDatabases.GET("/:databaseId/properties", apiV1.GetDatabasePropertiesById)
				notionDatabases.GET("/:databaseId/properties/options", apiV1.ListAllSelectOptions)
				notionDatabases.GET("/:databaseId/properties/options/:notionSelectId", apiV1.ListSelectOptions)
				//notionDatabases.DELETE("/cache", apiV1.DeleteDatabaseListFromCache)
				//notionDatabases.DELETE("/:id/cache", apiV1.DeleteDatabaseFromCache)
			}
		}
		formGroup := v1.Group("/form", Authenticate)
		{
			// TODO: add form-config and save this in the db -> form-url, other form-settings/config
			formGroup.POST("", apiV1.CreateForm)
			formGroup.GET("", apiV1.ListForm)
			formGroup.GET(":databaseId", apiV1.GetForm)
			formGroup.PUT(":databaseId", apiV1.UpdateForm)
			formGroup.DELETE(":databaseId", apiV1.DeleteForm)
		}
		storageGroup := v1.Group("/storage", Authenticate)
		{
			storageGroup.POST("/authenticate/google", apiV1.AuthenticateGoogleDrive)
			storageGroup.POST("/provider/:databaseId", apiV1.SetStorageProvider)
			storageGroup.POST("/location/:databaseId", apiV1.SetBaseStorageLocation)
			storageGroup.POST("/upload/:databaseId", apiV1.UploadFile)
		}
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	err := r.Run(fmt.Sprintf(":%d", apiConfig.Port))
	return err
}

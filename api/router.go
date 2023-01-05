package api

import (
	apiV1 "Notion-Forms/api/v1"
	"github.com/gin-gonic/gin"
)

func Router() {
	r := gin.New()
	v1 := r.Group("/api/v1")
	{
		notionPages := v1.Group("/pages")
		{
			notionPages.GET("/", apiV1.ListPages)
			notionPages.GET("/:type/:id", apiV1.ListPage)
		}
		notionRecords := v1.Group("/records")
		{
			notionRecords.POST("/:databaseId", apiV1.CreateRecord)
			notionRecords.GET("/:databaseId/options", apiV1.GetSelectAllOptions)
			notionRecords.GET("/:databaseId/options/:select", apiV1.GetSelectOptions)
		}
	}
	r.Run(":3000")
}

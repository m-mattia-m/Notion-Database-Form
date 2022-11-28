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
			notionPages.GET("/", apiV1.ListProjects)
		}
		notionDatabase := v1.Group("/databases")
		{
			notionDatabase.GET("/:id", apiV1.ListDatabase)
		}
	}
	r.Run(":3000")
}

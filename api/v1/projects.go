package v1

import (
	"Notion-Forms/pkg/notion"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func ListProjects(c *gin.Context) {
	showPagesQueryBool := false
	var err error
	showPagesQuery := c.Query("showpages")
	if showPagesQuery != "" {
		showPagesQueryBool, err = strconv.ParseBool(showPagesQuery)
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	pages, err := notion.ListAllPages(showPagesQueryBool)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, pages)
}

func ListDatabase(c *gin.Context) {
	id := c.Param("id")
	project, err := notion.ListDatabase(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, project)
}

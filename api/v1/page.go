package v1

import (
	"Notion-Forms/pkg/notion"
	"Notion-Forms/pkg/notion/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func ListPages(c *gin.Context) {
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

func ListPage(c *gin.Context) {
	id := c.Param("id")
	pageType := c.Param("type")

	if pageType == "database" {
		database, err := notion.ListDatabase(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, database)
		return
	}

	page, err := notion.ListPage(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, page)
}

func CreateRecord(c *gin.Context) {
	databaseId := c.Param("databaseId")
	var recordRequest []models.RecordRequest
	if error := c.BindJSON(&recordRequest); error != nil {
		fmt.Println(error)
		return
	}
	project, err := notion.CreateRecord(databaseId, recordRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, project)
}

func GetSelectAllOptions(c *gin.Context) {
	databaseId := c.Param("databaseId")
	options, err := notion.GetSelectAllOptions(databaseId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, options)
}

func GetSelectOptions(c *gin.Context) {
	databaseId := c.Param("databaseId")
	selectColumn := c.Param("select")
	options, err := notion.GetSelectOptions(databaseId, selectColumn)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusBadRequest, options)
}

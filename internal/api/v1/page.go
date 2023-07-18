package v1

import (
	"Notion-Forms/pkg/notion/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (api ApiClient) ListPages(c *gin.Context) {
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
	pages, err := api.Service.ListAllPages(showPagesQueryBool)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, pages)
}

func (api ApiClient) GetPage(c *gin.Context) {
	id := c.Param("id")
	pageType := c.Param("type")

	if pageType == "database" {
		database, err := api.Service.GetDatabase(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, database)
		return
	}

	page, err := api.Service.GetPage(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, page)
}

func (api ApiClient) CreateRecord(c *gin.Context) {
	databaseId := c.Param("databaseId")
	var recordRequest []model.RecordRequest
	if error := c.BindJSON(&recordRequest); error != nil {
		fmt.Println(error)
		return
	}
	project, err := api.Service.CreateRecord(databaseId, recordRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, project)
}

func (api ApiClient) ListAllSelectOptions(c *gin.Context) {
	databaseId := c.Param("databaseId")
	options, err := api.Service.ListAllSelectOptions(databaseId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, options)
}

func (api ApiClient) ListSelectOptions(c *gin.Context) {
	databaseId := c.Param("databaseId")
	selectColumn := c.Param("select")
	options, err := api.Service.ListSelectOptions(databaseId, selectColumn)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusBadRequest, options)
}

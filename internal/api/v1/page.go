package v1

import (
	"Notion-Forms/internal/helper"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListPages(c *gin.Context) {
	svc, _, _, err := getConfigAndService(c)
	if err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("failed to get config and service, contact the administrator"))
		return
	}

	pages, err := svc.ListPages()
	if err != nil {
		svc.SetAbortResponse(c, "svc", "ListPages", fmt.Sprintf("failed to list notion pages"), err)
		return
	}

	c.JSON(http.StatusOK, pages)
}

func GetPage(c *gin.Context) {
	svc, _, _, err := getConfigAndService(c)
	if err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("failed to get config and service, contact the administrator"))
		return
	}

	pageId := c.Param("pageId")
	if pageId == "" {
		helper.SetBadRequestResponse(c, fmt.Sprintf("page-id is required"))
		return
	}

	page, err := svc.GetPage(pageId)
	if err != nil {
		svc.SetAbortResponse(c, "svc", "GetPage", fmt.Sprintf("failed to get notion page by their id"), err)
		return
	}

	c.JSON(http.StatusOK, page)
}

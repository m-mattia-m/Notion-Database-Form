package v1

import (
	"Notion-Forms/internal/helper"
	"Notion-Forms/pkg/notion/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ListDatabases 			godoc
// @title           		ListDatabases
// @description     		Return a list of all own databases, where you have given access to it
// @Tags 					Notion
// @Router  				/notion/database [get]
// @Accept 					json
// @Produce					json
// @Security				Bearer
// @Success      			200  {object} []model.MinimalistDatabase
// @Failure      			400  {object} model.HttpError
// @Failure      			404  {object} model.HttpError
// @Failure      			500  {object} model.HttpError
func ListDatabases(c *gin.Context) {
	svc, _, _, err := getConfigAndService(c)
	// TODO: create an emergency log client which is controlled by eventbus in case the config or the service cannot be loaded correctly.
	if err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("failed to get config and service, contact the administrator"))
		return
	}

	databases, err := svc.ListDatabases()
	if err != nil {
		svc.SetAbortResponse(c, "svc", "ListDatabases", fmt.Sprintf("failed to list notion databases"), err)
		return
	}

	// TODO: Map notion-api-model to own rebuild model

	minimalistDatabaseList := svc.ConvertToMinimalistDatabaseList(databases)

	c.JSON(http.StatusOK, minimalistDatabaseList)
}

// GetDatabase 				godoc
// @title           		GetDatabase
// @description     		Return a own databases, where you have given access to it
// @Tags 					Notion
// @Router  				/notion/database/{databaseId} [get]
// @Accept 					json
// @Produce					json
// @Security				Bearer
// @Param        			databaseId    path     string  true  "databaseId"
// @Success      			200  {object} notion.Database
// @Failure      			400  {object} model.HttpError
// @Failure      			404  {object} model.HttpError
// @Failure      			500  {object} model.HttpError
func GetDatabase(c *gin.Context) {
	svc, _, _, err := getConfigAndService(c)
	if err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("failed to get config and service, contact the administrator"))
		return
	}

	databaseId := c.Param("databaseId")
	if databaseId == "" {
		helper.SetBadRequestResponse(c, fmt.Sprintf("database-id is required"))
	}

	// TODO: proof which person this database belongs to

	database, err := svc.GetDatabase(databaseId)
	if err != nil {
		svc.SetAbortResponse(c, "svc", "GetDatabase", fmt.Sprintf("failed to get notion databases by their id"), err)
		return
	}

	c.JSON(http.StatusOK, database)
}

// CreateRecord 			godoc
// @title           		CreateRecord
// @description     		Return a own databases, where you have given access to it
// @Tags 					Notion
// @Router  				/notion/database/{databaseId} [post]
// @Accept 					json
// @Produce					json
// @Security				Bearer
// @Param        			databaseId    	path     	string  				true  	"databaseId"
// @Param					RecordRequest 	body 		[]model.RecordRequest 	true 	"RecordRequest"
// @Success      			200  			{object} 	notion.Page
// @Failure      			400  			{object} 	model.HttpError
// @Failure      			404  			{object} 	model.HttpError
// @Failure      			500  			{object} 	model.HttpError
func CreateRecord(c *gin.Context) {
	svc, _, oidcUser, err := getConfigAndService(c)
	if err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("failed to get config and service, contact the administrator"))
		return
	}

	databaseId := c.Param("databaseId")
	if databaseId == "" {
		helper.SetBadRequestResponse(c, fmt.Sprintf("database-id is required"))
		return
	}

	userData, err := svc.GetOwnUser(*oidcUser)
	if err != nil {
		svc.SetAbortResponse(c, "svc", "ListDatabases", fmt.Sprintf("failed to list notion databases"), err)
		return
	}

	var recordRequest []model.RecordRequest
	if err := c.BindJSON(&recordRequest); err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("record-request-body can't bind to a object and is required"))
		return
	}

	project, err := svc.CreateRecord(databaseId, userData.NotionCredentials.UserId, recordRequest)
	if err != nil {
		svc.SetAbortResponse(c, "svc", "CreateRecord", fmt.Sprintf("failed to create notion-database record"), err)
		return
	}

	// TODO: Add success-response

	c.JSON(http.StatusOK, project)
}

// GetDatabasePropertiesById 	godoc
// @title           			GetDatabasePropertiesById
// @description     			Return the properties of a database by id
// @Tags 						Notion
// @Router  					/notion/database/{databaseId}/properties [get]
// @Accept 						json
// @Produce						json
// @Security					Bearer
// @Param        				databaseId    path     string  true  "databaseId"
// @Success      				200  {object} []model.DatabaseProperty
// @Failure      				400  {object} model.HttpError
// @Failure      				404  {object} model.HttpError
// @Failure      				500  {object} model.HttpError
func GetDatabasePropertiesById(c *gin.Context) {
	svc, _, _, err := getConfigAndService(c)
	if err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("failed to get config and service, contact the administrator"))
		return
	}

	databaseId := c.Param("databaseId")
	if databaseId == "" {
		helper.SetBadRequestResponse(c, fmt.Sprintf("database-id is required"))
	}

	// TODO: proof which person this database belongs to

	database, err := svc.GetDatabase(databaseId)
	if err != nil {
		svc.SetAbortResponse(c, "svc", "GetDatabase", fmt.Sprintf("failed to get notion databases by their id"), err)
		return
	}

	properties, err := svc.ConvertDatabaseToPropertyList(database)

	c.JSON(http.StatusOK, properties)
}

// ListAllSelectOptions 	godoc
// @title           		ListAllSelectOptions
// @description     		Return a list of all select options from a database by their id, where you have given access to it
// @Tags 					Notion
// @Router  				/notion/database/{databaseId}/properties/options [get]
// @Accept 					json
// @Produce					json
// @Security				Bearer
// @Param        			databaseId    	path     	string  				true  	"databaseId"
// @Success      			200  			{object} 	[]notionModel.Select
// @Failure      			400  			{object} 	model.HttpError
// @Failure      			404  			{object} 	model.HttpError
// @Failure      			500  			{object} 	model.HttpError
func ListAllSelectOptions(c *gin.Context) {
	svc, _, _, err := getConfigAndService(c)
	if err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("failed to get config and service, contact the administrator"))
		return
	}

	databaseId := c.Param("databaseId")
	if databaseId == "" {
		helper.SetBadRequestResponse(c, fmt.Sprintf("database-id is required"))
		return
	}

	options, err := svc.ListAllSelectOptions(databaseId)
	if err != nil {
		svc.SetAbortResponse(c, "svc", "ListAllSelectOptions", fmt.Sprintf("failed to list all notion-database-select-options"), err)
		return
	}

	c.JSON(http.StatusOK, options)
}

// ListSelectOptions	 	godoc
// @title           		ListSelectOptions
// @description     		Return a list of all select options from a database by their id, where you have given access to it
// @Tags 					Notion
// @Router  				/notion/database/{databaseId}/properties/options/{notionSelectId} [get]
// @Accept 					json
// @Produce					json
// @Security				Bearer
// @Param        			databaseId    	path     	string  				true  	"databaseId"
// @Param        			notionSelectId  path     	string  				true  	"notionSelectId"
// @Success      			200  			{object} 	[]notionModel.Select
// @Failure      			400  			{object} 	model.HttpError
// @Failure      			404  			{object} 	model.HttpError
// @Failure      			500  			{object} 	model.HttpError
func ListSelectOptions(c *gin.Context) {
	svc, _, _, err := getConfigAndService(c)
	if err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("failed to get config and service, contact the administrator"))
		return
	}

	databaseId := c.Param("databaseId")
	if databaseId == "" {
		helper.SetBadRequestResponse(c, fmt.Sprintf("database-id is required"))
		return
	}

	notionSelectId := c.Param("notionSelectId")
	if notionSelectId == "" {
		helper.SetBadRequestResponse(c, fmt.Sprintf("select-id is required"))
		return
	}

	options, err := svc.ListSelectOptions(databaseId, notionSelectId)
	if err != nil {
		svc.SetAbortResponse(c, "svc", "ListSelectOptions", fmt.Sprintf("failed to list notion-database-select-options"), err)
		return
	}

	c.JSON(http.StatusOK, options)
}

package v1

import (
	"Notion-Forms/internal/helper"
	"Notion-Forms/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateForm 				godoc
// @title           		CreateForm
// @description     		Create a new form
// @Tags 					Form
// @Router  				/form [post]
// @Param					FormRequest			body 		model.FormRequest 	true 		"FormRequest"
// @Accept 					json
// @Produce					json
// @Security				Bearer
// @Success      			200  				{object} 	model.FormResponse
// @Failure      			400  				{object} 	model.HttpError
// @Failure      			404  				{object} 	model.HttpError
// @Failure      			500  				{object} 	model.HttpError
func CreateForm(c *gin.Context) {
	svc, apiConfig, oidcUser, err := getConfigAndService(c)
	if err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("failed to get config and service, contact the administrator"))
		return
	}

	var formRequestBody model.FormRequest
	err = c.BindJSON(&formRequestBody)
	if err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("failed to convert request-body to object"))
		return
	}

	form, err := svc.CreateFormToDatabase(*oidcUser, formRequestBody)
	if err != nil {
		svc.SetAbortResponse(c, "svc", "SaveFormToDatabase", "failed to save form into the database", err)
		return
	}

	formResponse := helper.ConvertFormToFormResponse(*form, fmt.Sprintf("%s://%s", apiConfig.Schemas, apiConfig.Domain))

	c.JSON(http.StatusOK, formResponse)
}

// ListForm 				godoc
// @title           		ListForm
// @description     		List all own forms
// @Tags 					Form
// @Router  				/form [get]
// @Accept 					json
// @Produce					json
// @Security				Bearer
// @Success      			200  				{object} 	[]model.FormResponse
// @Failure      			400  				{object} 	model.HttpError
// @Failure      			404  				{object} 	model.HttpError
// @Failure      			500  				{object} 	model.HttpError
func ListForm(c *gin.Context) {
	svc, apiConfig, oidcUser, err := getConfigAndService(c)
	if err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("failed to get config and service, contact the administrator"))
		return
	}

	forms, err := svc.ListForms(*oidcUser)
	if err != nil {
		svc.SetAbortResponse(c, "svc", "ListForms", "failed to list forms from database", err)
		return
	}

	formResponse := helper.ConvertFormsListToFormResponse(forms, fmt.Sprintf("%s://%s", apiConfig.Schemas, apiConfig.Domain))

	c.JSON(http.StatusOK, formResponse)
}

// GetForm 					godoc
// @title           		GetForm
// @description     		Get a specific form by id
// @Tags 					Form
// @Router  				/form/{databaseId} [get]
// @Param        			databaseId    		path     	string  			true  		"databaseId"
// @Accept 					json
// @Produce					json
// @Security				Bearer
// @Success      			200  				{object} 	model.FormResponse
// @Failure      			400  				{object} 	model.HttpError
// @Failure      			404  				{object} 	model.HttpError
// @Failure      			500  				{object} 	model.HttpError
func GetForm(c *gin.Context) {
	svc, apiConfig, oidcUser, err := getConfigAndService(c)
	if err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("failed to get config and service, contact the administrator"))
		return
	}

	databaseId := c.Param("databaseId")
	if databaseId == "" {
		helper.SetBadRequestResponse(c, fmt.Sprintf("database-id is required"))
		return
	}

	form, err := svc.GetFormById(databaseId, *oidcUser)
	if err != nil {
		svc.SetAbortResponse(c, "svc", "GetFormById", "failed to get form from database by id", err)
		return
	}

	formResponse := helper.ConvertFormToFormResponse(*form, fmt.Sprintf("%s://%s", apiConfig.Schemas, apiConfig.Domain))

	c.JSON(http.StatusOK, formResponse)
}

// UpdateForm 				godoc
// @title           		UpdateForm
// @description     		Update a specific form by id
// @Tags 					Form
// @Router  				/form/{formId} [put]
// @Param        			formId    			path     	string  			true  		"formId"
// @Param					FormRequest			body 		model.FormRequest 	true 		"FormRequest"
// @Accept 					json
// @Produce					json
// @Security				Bearer
// @Success      			200  				{object} 	model.FormResponse
// @Failure      			400  				{object} 	model.HttpError
// @Failure      			404  				{object} 	model.HttpError
// @Failure      			500  				{object} 	model.HttpError
func UpdateForm(c *gin.Context) {
	svc, apiConfig, oidcUser, err := getConfigAndService(c)
	if err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("failed to get config and service, contact the administrator"))
		return
	}

	databaseId := c.Param("databaseId")
	if databaseId == "" {
		helper.SetBadRequestResponse(c, fmt.Sprintf("database-id is required"))
		return
	}

	var formRequestBody model.FormRequest
	err = c.BindJSON(&formRequestBody)
	if err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("failed to convert request-body to object"))
		return
	}

	formRequestBody.DatabaseId = databaseId

	form, err := svc.UpdateFormById(formRequestBody, *oidcUser)
	if err != nil {
		svc.SetAbortResponse(c, "svc", "UpdateFormById", "failed to update form from database by id", err)
		return
	}

	formResponse := helper.ConvertFormToFormResponse(*form, fmt.Sprintf("%s://%s", apiConfig.Schemas, apiConfig.Domain))

	c.JSON(http.StatusOK, formResponse)
}

// DeleteForm 				godoc
// @title           		DeleteForm
// @description     		Delete a specific form by id
// @Tags 					Form
// @Router  				/form/{formId} [delete]
// @Param        			formId    			path     	string  			true  		"formId"
// @Accept 					json
// @Produce					json
// @Security				Bearer
// @Success      			200  				{object} 	model.FormResponse
// @Failure      			400  				{object} 	model.HttpError
// @Failure      			404  				{object} 	model.HttpError
// @Failure      			500  				{object} 	model.HttpError
func DeleteForm(c *gin.Context) {
	svc, apiConfig, oidcUser, err := getConfigAndService(c)
	if err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("failed to get config and service, contact the administrator"))
		return
	}

	databaseId := c.Param("databaseId")
	if databaseId == "" {
		helper.SetBadRequestResponse(c, fmt.Sprintf("database-id is required"))
		return
	}

	form, err := svc.GetFormById(databaseId, *oidcUser)
	if err != nil {
		svc.SetAbortResponse(c, "svc", "GetFormById", "failed to get form from database by id", err)
		return
	}

	err = svc.DeleteFormById(databaseId, *oidcUser)
	if err != nil {
		svc.SetAbortResponse(c, "svc", "DeleteFormById", "failed to delete form from database by id", err)
		return
	}

	formResponse := helper.ConvertFormToFormResponse(*form, fmt.Sprintf("%s://%s", apiConfig.Schemas, apiConfig.Domain))

	c.JSON(http.StatusOK, formResponse)
}

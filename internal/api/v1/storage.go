package v1

import (
	"Notion-Forms/internal/helper"
	"Notion-Forms/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthenticateGoogleDrive 	godoc
// @title           		AuthenticateGoogleDrive
// @description     		Authenticate in Google Drive with the OAuth code and store the secret together with the IAM, Notion user and Google user
// @Tags 					Storage
// @Router  				/storage//authenticate/google [post]
// @Accept 					json
// @Produce					json
// @Security				Bearer
// @Param					OAuthCodeRequest	body 		model.OAuthCodeRequest 	true 	"RecordRequest"
// @Success      			200  				{object} 	nil
// @Failure      			400  				{object} 	model.HttpError
// @Failure      			404  				{object} 	model.HttpError
// @Failure      			500  				{object} 	model.HttpError
func AuthenticateGoogleDrive(c *gin.Context) {
	svc, _, oidcUser, err := getConfigAndService(c)
	if err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("failed to get config and service, contact the administrator"))
		return
	}

	var code model.OAuthCodeRequest
	err = c.BindJSON(&code)
	if err != nil {
		svc.SetAbortResponse(c, "c", "BindJSON", fmt.Sprintf("failed to bind oauth-code-request to object"), err)
		return
	}

	err = svc.ConnectIamUserWithGoogleUser(oidcUser.Sub, code.Code)
	if err != nil {
		svc.SetAbortResponse(c, "svc", "ConnectIamUserWithNotionUser", fmt.Sprintf("failed to authorzise google-user and save their id"), err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// SetStorageProvider 		godoc
// @title           		SetStorageProvider
// @description     		Set the provider (googleDrive, MicrosoftOneDrive, Dropbox) where the uploaded files should be stored
// @Tags 					Storage
// @Router  				/storage/provider/{databaseId} [post]
// @Accept 					json
// @Produce					json
// @Security				Bearer
// @Param        			databaseId    			path     		string  						true  	"databaseId"
// @Param					StorageProviderRequest 	body 			[]model.StorageProviderRequest 	true 	"StorageProviderRequest"
// @Success      			200  					{object} 		nil
// @Failure      			400  					{object} 		model.HttpError
// @Failure      			404  					{object} 		model.HttpError
// @Failure      			500  					{object} 		model.HttpError
func SetStorageProvider(c *gin.Context) {
	svc, _, _, err := getConfigAndService(c)
	if err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("failed to get config and service, contact the administrator"))
		return
	}

	databaseId := c.Param("databaseId")
	if databaseId == "" {
		helper.SetBadRequestResponse(c, fmt.Sprintf("database-id is required"))
	}

	var storageProviderRequest model.StorageProviderRequest
	if err := c.BindJSON(&storageProviderRequest); err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("storage-provider-request-body can't bind to a object and is required"))
		return
	}

	validProvider := helper.IsValidStorageProvider(storageProviderRequest.StorageProvider)
	if !validProvider {
		helper.SetBadRequestResponse(c, fmt.Sprintf("invalid storage-provider"))
		return
	}

	err = svc.SetStorageProvider(databaseId, model.StorageProvider(storageProviderRequest.StorageProvider))
	if err != nil {
		svc.SetAbortResponse(c, "svc", "SetStorageProvider", fmt.Sprintf("failed to save the provider"), err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// SetBaseStorageLocation 	godoc
// @title           		SetBaseStorageLocation
// @description     		Set the folderId where the uploaded files should be stored
// @Tags 					Storage
// @Router  				/storage/location/{databaseId} [post]
// @Accept 					json
// @Produce					json
// @Security				Bearer
// @Param        			databaseId    			path     		string  						true  	"databaseId"
// @Param					StorageLocationRequest	body 			model.StorageLocationRequest 	true 	"StorageLocationRequest"
// @Success      			200  					{object} 		nil
// @Failure      			400  					{object} 		model.HttpError
// @Failure      			404  					{object} 		model.HttpError
// @Failure      			500  					{object} 		model.HttpError
func SetBaseStorageLocation(c *gin.Context) {
	svc, _, _, err := getConfigAndService(c)
	if err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("failed to get config and service, contact the administrator"))
		return
	}

	databaseId := c.Param("databaseId")
	if databaseId == "" {
		helper.SetBadRequestResponse(c, fmt.Sprintf("database-id is required"))
	}

	var storageLocationRequest model.StorageLocationRequest
	if err := c.BindJSON(&storageLocationRequest); err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("storage-location-request-body can't bind to a object and is required"))
		return
	}

	err = svc.SetStorageLocation(databaseId, storageLocationRequest.ParentFolderId)
	if err != nil {
		svc.SetAbortResponse(c, "svc", "SetStorageLocation", fmt.Sprintf("failed to save the folder-id"), err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// UploadFile			 	godoc
// @title           		UploadFile
// @description     		Upload a file
// @Tags 					Storage
// @Router  				/storage/upload/{databaseId} [post]
// @Accept 					json
// @Produce					json
// @Security				Bearer
// @Param        			databaseId    			path     		string  					true  	"databaseId"
// @Param        			file    				formData 		file 						true  	"File"
// @Success      			200  					{object} 		model.FileUploadResponse
// @Failure      			400  					{object} 		model.HttpError
// @Failure      			404  					{object} 		model.HttpError
// @Failure      			500  					{object} 		model.HttpError
func UploadFile(c *gin.Context) {
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

	file, err := c.FormFile("file")
	if err != nil {
		helper.SetBadRequestResponse(c, fmt.Sprintf("upload file failed"))
		return
	}

	fileUrl, err := svc.UploadFile(oidcUser.Sub, databaseId, file)
	if err != nil {
		svc.SetAbortResponse(c, "svc", "UploadFile", fmt.Sprintf("failed to store the file"), err)
		return
	}

	c.JSON(http.StatusOK, model.FileUploadResponse{
		Url: *fileUrl,
	})
}

package service

import (
	"Notion-Forms/internal/config"
	"Notion-Forms/internal/model"
	googleDrive "Notion-Forms/pkg/storage/google-drive"
	"fmt"
	"mime/multipart"
)

func (svc Clients) ConnectIamUserWithGoogleUser(iamUserId, code string) error {
	oidcGoogleUser, err := svc.authenticateGoogleUser(code)
	if err != nil {
		return err
	}

	err = svc.db.dao.ConnectIamUserWithGoogleUser(iamUserId, oidcGoogleUser.AccessToken)
	if err != nil {
		return err
	}

	return nil
}

func (svc Clients) authenticateGoogleUser(code string) (*googleDrive.GoogleOauthToken, error) {
	googleOauthToken, err := svc.storage.googleDrive.Authenticate(code)
	if err != nil {
		return nil, err
	}
	return googleOauthToken, nil
}

func (svc Clients) SetStorageProvider(databaseId string, provider model.StorageProvider) error {
	return svc.db.dao.SetProvider(databaseId, provider)
}

func (svc Clients) SetStoragePath(databaseId, path string) error {
	return svc.db.dao.SetPath(databaseId, path)
}

func (svc Clients) GetStorageProvider(databaseId string) (*string, error) {
	return svc.db.dao.GetProvider(databaseId)
}

func (svc Clients) GetStoragePath(databaseId string) (*string, error) {
	return svc.db.dao.GetPath(databaseId)
}

func (svc Clients) UploadFile(databaseId string, file *multipart.FileHeader) error {
	provider, err := svc.GetStorageProvider(databaseId)
	if err != nil {
		return err
	}
	path, err := svc.GetStoragePath(databaseId)
	if err != nil {
		return err
	}

	switch model.StorageProvider(*provider) {
	case config.GoogleDrive:
		return svc.storage.googleDrive.UploadFile(*path, file)
	default:
		return fmt.Errorf("invalid storage provider: %v", model.StorageProvider(*provider))
	}
}

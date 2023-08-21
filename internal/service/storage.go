package service

import (
	"Notion-Forms/internal/config"
	"Notion-Forms/internal/model"
	googleDrive "Notion-Forms/pkg/storage/google-drive"
	"fmt"
	"mime/multipart"
)

func (svc Clients) ConnectIamUserWithGoogleUser(iamUserId, code string) error {
	oAuthGoogleUser, err := svc.authenticateGoogleUser(code)
	if err != nil {
		return err
	}

	err = svc.db.dao.ConnectIamUserWithGoogleUser(
		iamUserId,
		oAuthGoogleUser.AccessToken,
		oAuthGoogleUser.RefreshToken,
		oAuthGoogleUser.TokenType,
		oAuthGoogleUser.ExpiresIn,
		oAuthGoogleUser.Config,
	)
	if err != nil {
		return err
	}

	return nil
}

func (svc Clients) authenticateGoogleUser(code string) (*googleDrive.GoogleOauthTokenConfig, error) {
	googleOauthToken, err := svc.storage.googleDrive.Authenticate(code)
	if err != nil {
		return nil, err
	}
	return googleOauthToken, nil
}

func (svc Clients) SetStorageProvider(databaseId string, provider model.StorageProvider) error {
	return svc.db.dao.SetProvider(databaseId, provider)
}

func (svc Clients) SetStorageLocation(databaseId, folderId string) error {
	return svc.db.dao.SetLocation(databaseId, folderId)
}

func (svc Clients) GetStorageProvider(databaseId string) (*string, error) {
	return svc.db.dao.GetProvider(databaseId)
}

func (svc Clients) GetStorageFolderId(databaseId string) (*string, error) {
	return svc.db.dao.GetStorageFolderId(databaseId)
}

func (svc Clients) UploadFile(oidcUserId, databaseId string, file *multipart.FileHeader) (*string, error) {
	provider, err := svc.GetStorageProvider(databaseId)
	if err != nil {
		return nil, err
	}
	parentFolderId, err := svc.GetStorageFolderId(databaseId)
	if err != nil {
		return nil, err
	}

	credentials, err := svc.db.dao.GetUserDataByIamUserId(oidcUserId)
	if err != nil {
		return nil, err
	}

	switch model.StorageProvider(*provider) {
	case config.GoogleDrive:
		return svc.storage.googleDrive.UploadFile(googleDrive.GoogleOauthTokenConfig{
			Config:       credentials.GoogleCredentials.Config,
			AccessToken:  credentials.GoogleCredentials.AccessToken,
			RefreshToken: credentials.GoogleCredentials.RefreshToken,
			ExpiresIn:    credentials.GoogleCredentials.ExpiresIn,
			TokenType:    credentials.GoogleCredentials.TokenType,
		}, *parentFolderId, file)
	default:
		return nil, fmt.Errorf("invalid storage provider: %v", model.StorageProvider(*provider))
	}
}

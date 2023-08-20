package helper

import (
	"Notion-Forms/internal/config"
	"Notion-Forms/internal/model"
)

func IsValidStorageProvider(provider string) bool {
	switch model.StorageProvider(provider) {
	case config.GoogleDrive:
		return true
	case config.MicrosoftOneDrive:
		return true
	case config.Dropbox:
		return true
	default:
		return false
	}
}

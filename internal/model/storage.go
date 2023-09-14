package model

import "fmt"

type StorageProvider string

func (sp *StorageProvider) ToString() string {
	return fmt.Sprintf("%v", sp)
}

type Storage struct {
	StorageProvider StorageProvider `json:"storage_provider"`
	ParentFolderId  string          `json:"parent_folder_id"`
}

type StorageProviderRequest struct {
	StorageProvider string `json:"store_provider"`
}

type StorageLocationRequest struct {
	ParentFolderId string `json:"parent_folder_id"`
}

type FileUploadResponse struct {
	Url string `json:"url"`
}

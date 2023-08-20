package model

import "fmt"

type StorageProvider string

func (sp *StorageProvider) ToString() string {
	return fmt.Sprintf("%v", sp)
}

type Storage struct {
	DatabaseId      string          `json:"database_id"`
	StorageProvider StorageProvider `json:"storage_provider"`
	StorePath       string          `json:"store_path"`
}

type StorageProviderRequest struct {
	StorageProvider string `json:"store_provider"`
}

type StoragePathRequest struct {
	StoragePath string `json:"storage_path"`
}

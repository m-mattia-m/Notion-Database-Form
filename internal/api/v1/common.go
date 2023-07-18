package v1

import "Notion-Forms/internal/service"

type ApiClient struct {
	ApiConfig ApiConfig
	Service   *service.Service
}

type ApiConfig struct {
	RunMode         string
	FrontendHost    string
	DefaultPageSize int
	MaxPageSize     int
	Port            int
	Host            string
	Domain          string
	Schemas         string
}

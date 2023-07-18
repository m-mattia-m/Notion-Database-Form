package service

import (
	"Notion-Forms/pkg/notion/model"
	notion "github.com/jomei/notionapi"
)

func (svc Service) ListAllPages(showPagesQuery bool) ([]model.SiteResponse, error) {
	return svc.notion.ListAllPages(showPagesQuery)
}

func (svc Service) GetPage(id string) (notion.Page, error) {
	return svc.notion.GetPage(id)
}

func (svc Service) CreateRecord(databaseId string, requests []model.RecordRequest) (notion.Page, error) {
	return svc.notion.CreateRecord(databaseId, requests)
}

func (svc Service) ListSelectOptions(databaseId string, selectName string) ([]model.Select, error) {
	return svc.notion.ListSelectOptions(databaseId, selectName)
}

func (svc Service) ListAllSelectOptions(databaseId string) ([]model.Select, error) {
	return svc.notion.ListAllSelectOptions(databaseId)
}

func (svc Service) GetMe() (*model.User, error) {
	return svc.notion.GetMe()
}

func (svc Service) GetDatabase(id string) (notion.Database, error) {
	return svc.notion.GetDatabase(id)
}

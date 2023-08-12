package service

import (
	notion "github.com/jomei/notionapi"
)

func (svc Clients) ListPages() ([]*notion.Page, error) {
	return svc.notion.ListAllPages()
}

func (svc Clients) GetPage(id string) (notion.Page, error) {
	return svc.notion.GetPage(id)
}

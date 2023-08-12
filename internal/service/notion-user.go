package service

import "Notion-Forms/pkg/notion/model"

func (svc Clients) GetMe() (*model.User, error) {
	return svc.notion.GetMe()
}

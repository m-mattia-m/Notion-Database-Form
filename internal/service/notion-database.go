package service

import (
	"Notion-Forms/internal/listener"
	"Notion-Forms/internal/model"
	notionModel "Notion-Forms/pkg/notion/model"
	"encoding/json"
	notion "github.com/jomei/notionapi"
	"time"
)

// ListDatabases TODO: add user-logic here and then cache list by user-id -> GetDatabasesByUserId
func (svc Clients) ListDatabases() ([]*notion.Database, error) {
	return svc.notion.ListDatabases()
}

func (svc Clients) GetDatabase(id string) (notion.Database, error) {
	cachedString, err := svc.cache.Get(id)
	if err != nil {
		return notion.Database{}, err
	}
	if cachedString != nil {
		var databaseObject model.StoreDatabaseObject
		err = json.Unmarshal([]byte(*cachedString), &databaseObject)
		if err != nil {
			return notion.Database{}, err
		}

		listener.Bus.Publish("notion:update-database", id)
		return databaseObject.Object, nil
	}

	database, err := svc.notion.GetDatabase(id)
	if err != nil {
		return notion.Database{}, err
	}

	databaseString, err := json.Marshal(model.StoreDatabaseObject{
		Expiration:     time.Now(),
		RelevanceScore: 1,
		Object:         database,
	})
	err = svc.cache.Set(string(database.ID), string(databaseString))
	if err != nil {
		return notion.Database{}, err
	}

	return database, nil
}

func (svc Clients) CreateRecord(databaseId, userId string, requests []notionModel.RecordRequest) (notion.Page, error) {
	return svc.notion.CreateRecord(databaseId, userId, requests)
}

func (svc Clients) ListSelectOptions(databaseId string, selectName string) ([]notionModel.Select, error) {
	return svc.notion.ListSelectOptions(databaseId, selectName)
}

func (svc Clients) ListAllSelectOptions(databaseId string) ([]notionModel.Select, error) {
	return svc.notion.ListAllSelectOptions(databaseId)
}

package service

import (
	"Notion-Forms/internal/model"
	notionModel "Notion-Forms/pkg/notion/model"
	"encoding/json"
	"fmt"
	notion "github.com/jomei/notionapi"
	"time"
)

// ListDatabases TODO: add user-logic here and then cache list by user-id -> GetDatabasesByUserId
func (svc Clients) ListDatabases() ([]*notion.Database, error) {
	return svc.notion.ListDatabases()
}

func (svc Clients) GetDatabase(id string) (notion.Database, error) {
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

func (svc Clients) ListSelectOptions(databaseId string, selectId string) ([]notionModel.Select, error) {
	return svc.notion.ListSelectOptions(databaseId, selectId)
}

func (svc Clients) ListAllSelectOptions(databaseId string) ([]notionModel.Select, error) {
	return svc.notion.ListAllSelectOptions(databaseId)
}

func (svc Clients) ConvertToMinimalistDatabaseList(databaseList []*notion.Database) []model.MinimalistDatabase {
	var minimalistDatabaseList []model.MinimalistDatabase
	for _, database := range databaseList {
		minimalistDatabaseList = append(minimalistDatabaseList, model.MinimalistDatabase{
			Id:          database.ID.String(),
			Title:       svc.getNotionTitle(database.Title),
			Description: svc.getNotionTitle(database.Description),
			CreatedTime: database.CreatedTime.String(),
			Url:         database.URL,
		})
	}
	return minimalistDatabaseList
}

func (svc Clients) ConvertDatabaseToPropertyList(database notion.Database) ([]model.DatabasePropertyResponse, error) {
	var propertiesMap map[string]interface{}
	propertiesJson, err := json.Marshal(database.Properties)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(propertiesJson, &propertiesMap)
	if err != nil {
		return nil, err
	}

	var properties []model.DatabasePropertyResponse
	for key, value := range propertiesMap {
		var propertyMap map[string]interface{}
		propertyJson, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(propertyJson, &propertyMap)
		if err != nil {
			return nil, err
		}

		properties = append(properties, model.DatabasePropertyResponse{
			Id:   fmt.Sprintf("%v", propertyMap["id"]),
			Name: key,
			Type: fmt.Sprintf("%v", propertyMap["type"]),
		})
	}

	return properties, nil

}

func (svc Clients) getNotionTitle(title []notion.RichText) string {
	for _, titleElement := range title {
		return titleElement.PlainText
	}
	return ""
}

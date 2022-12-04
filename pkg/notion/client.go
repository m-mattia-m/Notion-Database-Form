package notion

import (
	"Notion-Forms/pkg/helper"
	"Notion-Forms/pkg/notion/models"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"github.com/dstotijn/go-notion"
)

var (
	ctx    context.Context
	client *notion.Client
)

func Client() error {
	key, err := helper.GetEnv("NOTION_SECRET_KEY")
	if err != nil {
		return err
	}
	client = notion.NewClient(key)
	ctx = context.Background()
	return nil
}

func ListAllPages(showPagesQuery bool) ([]models.SiteResponse, error) {
	resp, err := client.Search(ctx, &notion.SearchOpts{})
	if err != nil {
		return nil, err
	}

	var sites []models.SiteResponse
	for _, result := range resp.Results {
		database, status := result.(notion.Database)
		if status {
			databaseResp, err := returnDatabase(database)
			if err != nil {
				return nil, err
			}
			sites = append(sites, databaseResp)
			continue
		}
		if showPagesQuery {
			continue
		}
		page, status := result.(notion.Page)
		if status {
			pageResp, err := returnPage(page)
			if err != nil {
				return nil, err
			}
			sites = append(sites, pageResp)
			continue
		}
		return nil, errors.New("an object is listed which is neither a database nor a page, please check the response")
	}
	return sites, nil
}
func ListDatabase(id string) (notion.Database, error) {
	resp, err := client.FindDatabaseByID(ctx, id)
	if err != nil {
		return notion.Database{}, err
	}
	return resp, nil
}
func CreateRecord(databaseId string, requests []models.RecordRequest) (notion.Page, error) {

	var databaseProperties notion.DatabasePageProperties
	var databaseProperty map[string]interface{}
	for _, request := range requests {
		switch request.Type {
		case "title":
			databaseProperty = map[string]interface{}{request.Name: notion.DatabaseProperty{
				Name: request.Value.(string),
			}}
			break
		case "rich_text":
			propertyType := flag.String("notion.DatabasePropertyType", request.Type, "")
			//propertyValue := flag.String("notion.DatabasePropertyValue", request.Value.(string), "")
			databaseProperty = map[string]interface{}{request.Name: notion.DatabaseProperty{
				Type:     notion.DatabasePropertyType(*propertyType),
				RichText: &notion.EmptyMetadata{},
			}}
			break
		case "number":

			break
		case "select":
			break
		case "multi_select":
			break
		case "date":
			break
		case "people":
			break
		case "files":
			break
		case "checkbox":
			break
		case "url":
			break
		case "email":
			break
		case "phone_number":
			break
		case "formula":
			break
		case "relation":
			break
		case "rollup":
			break
		case "created_time":
			break
		case "created_by":
			break
		case "last_edited_time":
			break
		case "last_edited_by":
			break
		case "status":
			break
		default:
			return notion.Page{}, errors.New("there is no type in notion that is called this: " + request.Type)
		}

		//databasePropertiesMap := map[string]interface{}{request.Name: notion.DatabaseProperty{
		//	Name: request.Name,
		//}}

		databasePropertiesString, err := json.Marshal(databaseProperty)
		if err != nil {
			return notion.Page{}, err
		}
		err = json.Unmarshal(databasePropertiesString, &databaseProperties)
		if err != nil {
			return notion.Page{}, err
		}
	}

	resp, err := client.CreatePage(ctx, notion.CreatePageParams{
		ParentType:             "database",
		ParentID:               databaseId,
		DatabasePageProperties: &databaseProperties,
		Children:               nil,
		Icon:                   nil,
		Cover:                  nil,
	})
	if err != nil {
		return notion.Page{}, err
	}
	_ = resp
	return notion.Page{}, nil

}

func returnPage(page notion.Page) (models.SiteResponse, error) {
	// The problem is, this is a database entry, which has no fixed structure. There does not have to be a name, because this column can also be renamed. The same with the description.
	return models.SiteResponse{
		Id:          page.ID,
		Name:        "",
		Description: "",
		Author:      page.CreatedBy.ID,
		Type:        "page",
	}, nil

}
func returnDatabase(database notion.Database) (models.SiteResponse, error) {
	description := ""
	if len(database.Description) != 0 && database.Description[0].PlainText != "" {
		description = database.Description[0].PlainText
	}
	return models.SiteResponse{
		Id:          database.ID,
		Name:        database.Title[0].PlainText,
		Description: description,
		Author:      database.CreatedBy.ID,
		Type:        "database",
	}, nil

}

package notion

import (
	"Notion-Forms/pkg/helper"
	"Notion-Forms/pkg/notion/models"
	"context"
	"errors"
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

	sites := []models.SiteResponse{}
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

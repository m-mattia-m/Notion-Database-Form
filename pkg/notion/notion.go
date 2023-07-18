package notion

import (
	"Notion-Forms/pkg/notion/model"
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	notion "github.com/jomei/notionapi"
	"os"
	"strconv"
	"time"
)

func init() {
	godotenv.Load()
}

type Client struct {
	ctx    context.Context
	client *notion.Client
	user   *model.User
}

func New() (*Client, error) {
	secretKey := os.Getenv("NOTION_SECRET_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("failed to get env-var 'NOTION_SECRET_KEY'")
	}

	notionClient := notion.NewClient(notion.Token(secretKey))
	ctx := context.Background()

	user, err := notionClient.User.Me(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get notion-user: %s", err)
	}

	return &Client{
		ctx:    ctx,
		client: notionClient,
		user: &model.User{
			Id: user.ID.String(),
		},
	}, nil
}

func (c *Client) ListAllPages(showPagesQuery bool) ([]model.SiteResponse, error) {
	resp, err := c.client.Search.Do(c.ctx, &notion.SearchRequest{})

	if err != nil {
		return nil, err
	}

	var sites []model.SiteResponse
	for _, result := range resp.Results {
		object := result.GetObject()
		if object == "database" {
			database, status := result.(*notion.Database)
			if !status {
				return nil, errors.New("can't cast the notion-search-response to database")
			}

			databaseResp, err := returnDatabase(*database)
			if err != nil {
				return nil, err
			}
			sites = append(sites, databaseResp)
			continue
		}
		if showPagesQuery {
			continue
		}

		page, status := result.(*notion.Page)
		if !status {
			return nil, errors.New("can't cast the notion-search-response to page")
		}

		pageResp, err := returnPage(*page)
		if err != nil {
			return nil, err
		}
		sites = append(sites, pageResp)
		continue
	}

	return sites, err
}

func (c *Client) GetDatabase(id string) (notion.Database, error) {
	resp, err := c.client.Database.Get(c.ctx, notion.DatabaseID(id))
	if err != nil {
		return notion.Database{}, err
	}
	return *resp, nil
}

func (c *Client) GetPage(id string) (notion.Page, error) {
	resp, err := c.client.Page.Get(c.ctx, notion.PageID(id))
	if err != nil {
		return notion.Page{}, err
	}
	return *resp, nil
}

func (c *Client) CreateRecord(databaseId string, requests []model.RecordRequest) (notion.Page, error) {
	properties := map[string]notion.Property{}

	for _, request := range requests {
		fmt.Println("")
		switch request.Type {
		case "title":
			value, status := request.Value.(string)
			if !status {
				return notion.Page{}, errors.New("can't cast request.Value to string")
			}
			var richTextArray []notion.RichText
			richText := notion.RichText{
				Type: "text",
				Text: &notion.Text{
					Content: value,
				},
			}
			richTextArray = append(richTextArray, richText)
			properties[request.Name] = notion.TitleProperty{
				Title: richTextArray,
			}
			break
		case "rich_text":
			value, status := request.Value.(string)
			if !status {
				return notion.Page{}, errors.New("can't cast request.Value to string")
			}
			var richTextArray []notion.RichText
			richTextArray = append(richTextArray, notion.RichText{
				Type: "text",
				Text: &notion.Text{
					Content: value,
				},
				PlainText: value,
			})
			properties[request.Name] = notion.RichTextProperty{
				RichText: richTextArray,
			}
			break
		case "number":
			value, status := request.Value.(string)
			if !status {
				return notion.Page{}, errors.New("can't cast request.Value to string")
			}
			number, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return notion.Page{}, err
			}
			properties[request.Name] = notion.NumberProperty{
				Number: number,
			}
			break
		case "select":
			value, status := request.Value.(string)
			if !status {
				return notion.Page{}, errors.New("can't cast request.Value to string")
			}
			properties[request.Name] = notion.SelectProperty{
				Select: notion.Option{
					ID:   notion.PropertyID(value),
					Name: "undefined",
				},
			}
			break
		case "multi_select":
			optionsArray, err := toStringArray(request.Value)
			if err != nil {
				return notion.Page{}, err
			}
			var options []notion.Option
			for _, option := range optionsArray {
				options = append(options, notion.Option{
					ID:   notion.PropertyID(option),
					Name: "undefined",
				})
			}
			properties[request.Name] = notion.MultiSelectProperty{
				MultiSelect: options,
			}
			break
		case "date":
			var notionDateProperty = notion.DateProperty{
				Date: &notion.DateObject{
					Start: nil,
					End:   nil,
				},
			}
			dates, err := toStringArray(request.Value)
			if err != nil {
				return notion.Page{}, err
			}
			startDate, err := time.Parse(time.RFC3339, dates[0])
			if err != nil {
				return notion.Page{}, err
			}
			notionDateProperty.Date.Start = (*notion.Date)(&startDate)
			var endDate time.Time
			if dates[1] != "" {
				endDate, err = time.Parse(time.RFC3339, dates[1])
				if err != nil {
					return notion.Page{}, err
				}
				notionDateProperty.Date.End = (*notion.Date)(&endDate)
			}
			properties[request.Name] = notionDateProperty
			break
		case "people":
			var users []notion.User
			usersArray, err := toStringArray(request.Value)
			if err != nil {
				return notion.Page{}, err
			}
			for _, user := range usersArray {
				users = append(users, notion.User{
					ID: notion.UserID(user),
				})
			}
			properties[request.Name] = notion.PeopleProperty{
				People: users,
			}
			break
		case "files":
			var files []notion.File
			filesArray, err := toStringArray(request.Value)
			if err != nil {
				return notion.Page{}, err
			}
			for _, user := range filesArray {
				url, err := uploadFile(user)
				if err != nil {
					return notion.Page{}, err
				}
				files = append(files, notion.File{
					File: &notion.FileObject{
						URL: url,
					},
					External: &notion.FileObject{
						URL: url,
					},
				})
				_ = user
			}
			properties[request.Name] = notion.FilesProperty{
				Files: files,
			}
			break
		case "checkbox":
			value, status := request.Value.(bool)
			if !status {
				return notion.Page{}, errors.New("can't cast request.Value to bool")
			}
			properties[request.Name] = notion.CheckboxProperty{
				Checkbox: value,
			}
			break
		case "url":
			value, status := request.Value.(string)
			if !status {
				return notion.Page{}, errors.New("can't cast request.Value to string")
			}
			properties[request.Name] = notion.URLProperty{
				URL: value,
			}
			break
		case "email":
			value, status := request.Value.(string)
			if !status {
				return notion.Page{}, errors.New("can't cast request.Value to string")
			}
			properties[request.Name] = notion.EmailProperty{
				Email: value,
			}
			break
		case "phone_number":
			value, status := request.Value.(string)
			if !status {
				return notion.Page{}, errors.New("can't cast request.Value to string")
			}
			properties[request.Name] = notion.PhoneNumberProperty{
				PhoneNumber: value,
			}
			break
		case "formula":
			// 1:1-string of the formula-value -> something like this: prop("Text") -> make no sense to add this, because you can't create this plain in the frontend
			break
		case "relation":
			// make no sense to add this, because you can't create this plain in the frontend
			break
		case "rollup":
			// make no sense to add this, because you can't create this plain in the frontend
			break
		case "created_time":
			date := time.Now()
			properties[request.Name] = notion.DateProperty{
				Date: &notion.DateObject{
					Start: (*notion.Date)(&date),
					End:   nil,
				},
			}
			break
		case "created_by":
			var users []notion.User
			users = append(users, notion.User{ID: notion.UserID(c.user.Id)})
			properties[request.Name] = notion.PeopleProperty{
				People: users,
			}
			break
		case "last_edited_time":
			date := time.Now()
			properties[request.Name] = notion.DateProperty{
				Date: &notion.DateObject{
					Start: (*notion.Date)(&date),
					End:   nil,
				},
			}
			break
		case "last_edited_by":
			var users []notion.User
			users = append(users, notion.User{ID: notion.UserID(c.user.Id)})
			properties[request.Name] = notion.PeopleProperty{
				People: users,
			}
			break
		case "status":
			value, status := request.Value.(string)
			if !status {
				return notion.Page{}, errors.New("can't cast request.Value to string")
			}
			properties[request.Name] = notion.StatusProperty{
				Status: notion.Status{
					ID:   notion.PropertyID(value),
					Name: "undefined",
				},
			}
			break
		default:
			return notion.Page{}, errors.New("there is no type in notion that is called this: " + request.Type)
		}
	}

	resp, err := c.client.Page.Create(c.ctx, &notion.PageCreateRequest{
		Parent: notion.Parent{
			Type:       "database_id",
			DatabaseID: notion.DatabaseID(databaseId),
			BlockID:    "",
			Workspace:  false,
		},
		Properties: properties,
		Children:   nil,
		Icon:       nil,
		Cover:      nil,
	})
	if err != nil {
		return notion.Page{}, err
	}
	_ = resp
	return notion.Page{}, nil
}

func (c *Client) ListSelectOptions(databaseId string, selectName string) ([]model.Select, error) {
	database, err := c.client.Database.Get(c.ctx, notion.DatabaseID(databaseId))
	if err != nil {
		return nil, err
	}
	var selects []model.Select
	for key, property := range database.Properties {
		if key != selectName {
			continue
		}
		if property.GetType() == "multi_select" {
			propertyObject, status := property.(*notion.MultiSelectPropertyConfig)
			if !status {
				return nil, errors.New("can't cast interface to MultiSelectPropertyConfig")
			}
			var options []model.Option
			for _, option := range propertyObject.MultiSelect.Options {
				options = append(options, model.Option{
					Id:   string(option.ID),
					Name: option.Name,
				})
			}
			selects = append(selects, model.Select{
				Id:      string(propertyObject.ID),
				Name:    key,
				Options: options,
			})
		}
		if property.GetType() == "select" {
			propertyObject, status := property.(*notion.SelectPropertyConfig)
			if !status {
				return nil, errors.New("can't cast interface to SelectPropertyConfig")
			}
			var options []model.Option
			for _, option := range propertyObject.Select.Options {
				options = append(options, model.Option{
					Id:   string(option.ID),
					Name: option.Name,
				})
			}
			selects = append(selects, model.Select{
				Id:      string(propertyObject.ID),
				Name:    key,
				Options: options,
			})

		}
	}
	if len(selects) == 0 {
		return nil, errors.New("there is no select or multiselect in this database")
	}
	return selects, nil
}

func (c *Client) ListAllSelectOptions(databaseId string) ([]model.Select, error) {
	database, err := c.client.Database.Get(c.ctx, notion.DatabaseID(databaseId))
	if err != nil {
		return nil, err
	}
	var selects []model.Select
	for key, property := range database.Properties {
		if property.GetType() == "multi_select" {
			propertyObject, status := property.(*notion.MultiSelectPropertyConfig)
			if !status {
				return nil, errors.New("can't cast interface to MultiSelectPropertyConfig")
			}
			var options []model.Option
			for _, option := range propertyObject.MultiSelect.Options {
				options = append(options, model.Option{
					Id:   string(option.ID),
					Name: option.Name,
				})
			}
			selects = append(selects, model.Select{
				Id:      string(propertyObject.ID),
				Name:    key,
				Options: options,
			})
		}
		if property.GetType() == "select" {
			propertyObject, status := property.(*notion.SelectPropertyConfig)
			if !status {
				return nil, errors.New("can't cast interface to SelectPropertyConfig")
			}
			var options []model.Option
			for _, option := range propertyObject.Select.Options {
				options = append(options, model.Option{
					Id:   string(option.ID),
					Name: option.Name,
				})
			}
			selects = append(selects, model.Select{
				Id:      string(propertyObject.ID),
				Name:    key,
				Options: options,
			})

		}
	}
	if len(selects) == 0 {
		return nil, errors.New("there is no select or multiselect in this database")
	}
	return selects, nil
}

func (c *Client) GetMe() (*model.User, error) {
	user, _ := c.client.User.Me(c.ctx)
	return &model.User{
		Id: user.ID.String(), // 7403beb1-603d-4f32-9014-b15a8e4212c7
	}, nil
}

func returnPage(page notion.Page) (model.SiteResponse, error) {
	// The problem is, this is a database entry, which has no fixed structure. There does not have to be a name, because this column can also be renamed. The same with the description.
	return model.SiteResponse{
		Id:          page.ID.String(),
		Name:        "",
		Description: "",
		Author:      page.CreatedBy.ID.String(),
		Type:        "page",
	}, nil

}

func returnDatabase(database notion.Database) (model.SiteResponse, error) {
	description := ""
	if len(database.Description) != 0 && database.Description[0].PlainText != "" {
		description = database.Description[0].PlainText
	}
	return model.SiteResponse{
		Id:          database.ID.String(),
		Name:        database.Title[0].PlainText,
		Description: description,
		Author:      database.CreatedBy.ID.String(),
		Type:        "database",
	}, nil

}

func toStringArray(object interface{}) ([]string, error) {
	interfaceArray, status := object.([]interface{})
	if !status {
		return nil, errors.New("can't cast interface to interface-array")
	}
	stringArray := make([]string, len(interfaceArray))
	for i, v := range interfaceArray {
		castedString, status := v.(string)
		if !status {
			return nil, errors.New("can't cast interface to string")
		}
		stringArray[i] = castedString
	}
	return stringArray, nil
}

func uploadFile(fileString string) (string, error) {
	return "", nil
}

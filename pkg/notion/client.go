package notion

import (
	"Notion-Forms/pkg/helper"
	"Notion-Forms/pkg/notion/models"
	"context"
	"errors"
	"fmt"
	notion "github.com/jomei/notionapi"
	"strconv"
	"time"
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
	//client = notion.NewClient(key)
	client = notion.NewClient(notion.Token(key))
	ctx = context.Background()
	return nil
}

func ListAllPages(showPagesQuery bool) ([]models.SiteResponse, error) {
	resp, err := client.Search.Do(ctx, &notion.SearchRequest{})

	if err != nil {
		return nil, err
	}

	var sites []models.SiteResponse
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

func ListDatabase(id string) (notion.Database, error) {
	resp, err := client.Database.Get(ctx, notion.DatabaseID(id))
	if err != nil {
		return notion.Database{}, err
	}
	return *resp, nil
}
func CreateRecord(databaseId string, requests []models.RecordRequest) (notion.Page, error) {
	properties := map[string]notion.Property{}

	//properties["name"] = notion.TitleProperty{
	//	Title: nil,
	//}

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
			value, status := request.Value.(string)
			if !status {
				return notion.Page{}, errors.New("can't cast request.Value to string")
			}
			properties[request.Name] = notion.StatusProperty{
				Status: notion.Status{
					ID:   notion.ObjectID(value),
					Name: "undefined",
				},
			}
			break
		default:
			return notion.Page{}, errors.New("there is no type in notion that is called this: " + request.Type)
		}
	}

	resp, err := client.Page.Create(ctx, &notion.PageCreateRequest{
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

func returnPage(page notion.Page) (models.SiteResponse, error) {
	// The problem is, this is a database entry, which has no fixed structure. There does not have to be a name, because this column can also be renamed. The same with the description.
	return models.SiteResponse{
		Id:          page.ID.String(),
		Name:        "",
		Description: "",
		Author:      page.CreatedBy.ID.String(),
		Type:        "page",
	}, nil

}
func returnDatabase(database notion.Database) (models.SiteResponse, error) {
	description := ""
	if len(database.Description) != 0 && database.Description[0].PlainText != "" {
		description = database.Description[0].PlainText
	}
	return models.SiteResponse{
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

//
//func ListAllPages(showPagesQuery bool) ([]models.SiteResponse, error) {
//	resp, err := client.Search.Do(ctx, &notion.SearchRequest{})
//
//	if err != nil {
//		return nil, err
//	}
//
//	var sites []models.SiteResponse
//	for _, result := range resp.Results {
//		database := result.GetObject()
//
//			databaseResp, err := returnDatabase(database)
//			if err != nil {
//				return nil, err
//			}
//			sites = append(sites, databaseResp)
//			continue
//
//		if showPagesQuery {
//			continue
//		}
//		page, status := result.(notion.Page)
//		if status {
//			pageResp, err := returnPage(page)
//			if err != nil {
//				return nil, err
//			}
//			sites = append(sites, pageResp)
//			continue
//		}
//		return nil, errors.New("an object is listed which is neither a database nor a page, please check the response")
//	}
//	return sites, nil
//}
//func ListDatabase(id string) (notion.Database, error) {
//	resp, err := client.FindDatabaseByID(ctx, id)
//	if err != nil {
//		return notion.Database{}, err
//	}
//	return resp, nil
//}
//func CreateRecord(databaseId string, requests []models.RecordRequest) (notion.Page, error) {
//
//	var databaseProperties notion.DatabasePageProperties
//	var databaseProperty map[string]interface{}
//	for _, request := range requests {
//		switch request.Type {
//		case "title":
//			databaseProperty = map[string]interface{}{request.Name: notion.DatabaseProperty{
//				Name: request.Value.(string),
//			}}
//			break
//		case "rich_text":
//			propertyType := flag.String("notion.DatabasePropertyType", request.Type, "")
//			//propertyValue := flag.String("notion.DatabasePropertyValue", request.Value.(string), "")
//			databaseProperty = map[string]interface{}{request.Name: notion.DatabaseProperty{
//				Type:     notion.DatabasePropertyType(*propertyType),
//				RichText: &notion.EmptyMetadata{},
//			}}
//			break
//		case "number":
//
//			break
//		case "select":
//			break
//		case "multi_select":
//			break
//		case "date":
//			break
//		case "people":
//			break
//		case "files":
//			break
//		case "checkbox":
//			break
//		case "url":
//			break
//		case "email":
//			break
//		case "phone_number":
//			break
//		case "formula":
//			break
//		case "relation":
//			break
//		case "rollup":
//			break
//		case "created_time":
//			break
//		case "created_by":
//			break
//		case "last_edited_time":
//			break
//		case "last_edited_by":
//			break
//		case "status":
//			break
//		default:
//			return notion.Page{}, errors.New("there is no type in notion that is called this: " + request.Type)
//		}
//
//		//databasePropertiesMap := map[string]interface{}{request.Name: notion.DatabaseProperty{
//		//	Name: request.Name,
//		//}}
//
//		databasePropertiesString, err := json.Marshal(databaseProperty)
//		if err != nil {
//			return notion.Page{}, err
//		}
//		err = json.Unmarshal(databasePropertiesString, &databaseProperties)
//		if err != nil {
//			return notion.Page{}, err
//		}
//	}
//
//	resp, err := client.CreatePage(ctx, notion.CreatePageParams{
//		ParentType:             "database",
//		ParentID:               databaseId,
//		DatabasePageProperties: &databaseProperties,
//		Children:               nil,
//		Icon:                   nil,
//		Cover:                  nil,
//	})
//	if err != nil {
//		return notion.Page{}, err
//	}
//	_ = resp
//	return notion.Page{}, nil
//
//}
//
//func returnPage(page notion.Page) (models.SiteResponse, error) {
//	// The problem is, this is a database entry, which has no fixed structure. There does not have to be a name, because this column can also be renamed. The same with the description.
//	return models.SiteResponse{
//		Id:          page.ID,
//		Name:        "",
//		Description: "",
//		Author:      page.CreatedBy.ID,
//		Type:        "page",
//	}, nil
//
//}
//func returnDatabase(database notion.ObjectType) (models.SiteResponse, error) {
//	description := ""
//	if len(database.Description) != 0 && database.Description[0].PlainText != "" {
//		description = database.Description[0].PlainText
//	}
//	return models.SiteResponse{
//		Id:          database.ID,
//		Name:        database.Title[0].PlainText,
//		Description: description,
//		Author:      database.CreatedBy.ID,
//		Type:        "database",
//	}, nil
//
//}

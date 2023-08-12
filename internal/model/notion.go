package model

import (
	notion "github.com/jomei/notionapi"
	"time"
)

type StoreObject struct {
	Expiration     time.Time   `json:"expiration"`
	RelevanceScore int         `json:"relevance_score"`
	Object         interface{} `json:"object"`
}
type StoreDatabaseObject struct {
	Expiration     time.Time       `json:"expiration"`
	RelevanceScore int             `json:"relevance_score"`
	Object         notion.Database `json:"object"`
}
type StorePageObject struct {
	Expiration     time.Time   `json:"expiration"`
	RelevanceScore int         `json:"relevance_score"`
	Object         notion.Page `json:"object"`
}
type IamUserNotionUser struct {
	IamUserId    string `json:"iam_id"`
	NotionUserId string `json:"notion_user_id"`
	NotionBotId  string `json:"notion_bot_id"`

	// extra no json, so that the possibility does not arise accidentally output the
	// conclusions, because this would be a big security risk.
	// https://developers.notion.com/docs/authorization#step-5-the-integration-stores-the-access_token-for-future-requests
	NotionAccessToken string `json:"-"`
}
type OAuthCodeRequest struct {
	Code string `json:"code"`
}

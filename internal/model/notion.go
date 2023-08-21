package model

import (
	notion "github.com/jomei/notionapi"
	"golang.org/x/oauth2"
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

// GNFUser
// extra no json for `NotionCredentials`, `GoogleCredentials`, so that the possibility does not arise accidentally output the
// conclusions, because this would be a big security risk.
// https://developers.notion.com/docs/authorization#step-5-the-integration-stores-the-access_token-for-future-requests
type GNFUser struct {
	IamUserId         string `json:"iam_id"`
	NotionCredentials NotionCredentials
	GoogleCredentials GoogleCredentials
}
type NotionCredentials struct {
	UserId      string `json:"user_id"`
	BotId       string `json:"bot_id"`
	AccessToken string `json:"-"`
}
type GoogleCredentials struct {
	Config       oauth2.Config `json:"config"`
	AccessToken  string        `json:"-"`
	RefreshToken string        `json:"-"`
	ExpiresIn    time.Time     `json:"expires_in"`
	TokenType    string        `json:"token_type"`
}

type OAuthCodeRequest struct {
	Code string `json:"code"`
}

type MinimalistDatabase struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedTime string `json:"created_time"`
	Url         string `json:"url"`
}
type DatabasePropertyResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}
type DatabasePropertySelectOptions struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

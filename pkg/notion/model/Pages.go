package model

type CreatePageRequest struct {
	Id          string `json:"id"`
	Title       string `json:"name"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Type        string `json:"type"`
}

type RecordRequest struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

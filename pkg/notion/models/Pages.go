package models

type SiteResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Type        string `json:"type"`
}

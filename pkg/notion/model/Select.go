package model

type Option struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Select struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Options []Option `json:"options"`
}

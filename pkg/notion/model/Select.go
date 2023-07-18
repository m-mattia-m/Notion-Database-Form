package model

type Option struct {
	Id   string `json: "id"`
	Name string `json: "name"`
}

type Select struct {
	Id      string   `json: "id"`
	Name    string   `json: "name"`
	Options []Option `json: "options"`
}

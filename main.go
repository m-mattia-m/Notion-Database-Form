package main

import (
	"Notion-Forms/api"
	"Notion-Forms/pkg/notion"
	"fmt"
)

func main() {
	err := notion.Client()
	if err != nil {
		fmt.Println("Notion-Client-Err: " + err.Error())
	}

	api.Router()
}

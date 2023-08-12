package notion

import (
	"fmt"
	notion "github.com/jomei/notionapi"
)

func (c *Client) ListDatabases() ([]*notion.Database, error) {
	resp, err := c.client.Search.Do(c.ctx, &notion.SearchRequest{
		Filter: notion.SearchFilter{
			Property: "object",
			Value:    "database",
		},
	})

	if err != nil {
		return nil, err
	}

	var databases []*notion.Database
	for _, result := range resp.Results {
		database, status := result.(*notion.Database)
		if !status {
			return nil, fmt.Errorf("can't cast the notion-search-response to database")
		}
		databases = append(databases, database)
	}

	return databases, err
}

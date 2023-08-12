package notion

import (
	"fmt"
	notion "github.com/jomei/notionapi"
)

func (c *Client) ListAllPages() ([]*notion.Page, error) {
	resp, err := c.client.Search.Do(c.ctx, &notion.SearchRequest{
		Filter: notion.SearchFilter{
			Property: "object",
			Value:    "page",
		},
	})
	if err != nil {
		return nil, err
	}

	var pages []*notion.Page
	for _, result := range resp.Results {
		page, status := result.(*notion.Page)
		if !status {
			return nil, fmt.Errorf("can't cast the notion-search-response to page")
		}
		pages = append(pages, page)
	}

	return pages, err
}

func (c *Client) GetPage(id string) (notion.Page, error) {
	resp, err := c.client.Page.Get(c.ctx, notion.PageID(id))
	if err != nil {
		return notion.Page{}, err
	}
	return *resp, nil
}

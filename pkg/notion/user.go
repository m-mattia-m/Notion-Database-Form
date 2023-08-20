package notion

import (
	"Notion-Forms/pkg/notion/model"
	"context"
	"encoding/json"
	notion "github.com/jomei/notionapi"
)

func (c *Client) GetMe() (*model.User, error) {
	user, _ := c.client.User.Me(c.ctx)
	return &model.User{
		Id: user.ID.String(),
	}, nil
}

func (c *Client) Authenticate(redirectUri, code string) (*model.OAuthToken, error) {

	tokenResponse, err := c.client.Authentication.CreateToken(context.TODO(), &notion.TokenCreateRequest{
		Code:        code,
		GrantType:   "authorization_code",
		RedirectUri: redirectUri,
	})
	if err != nil {
		return nil, err
	}

	var token model.OAuthToken
	jsonToken, err := json.Marshal(tokenResponse)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonToken, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil

	//values.Add("code", "4/0Adeu5BVlPkv6Qaq6rKF7Mnzy2rf1h0H_MKZnn2OwjGQ5xj4gwJz8Z0rwhvOVJ-vtObwQ8w")
	//values.Add("client_id", "554221747409-bhikmafi9hod48vcvbm5clig58it5d9e.apps.googleusercontent.com")
	//values.Add("client_secret", "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")
}

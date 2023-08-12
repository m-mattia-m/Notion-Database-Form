package notion

import (
	"Notion-Forms/pkg/notion/model"
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

	tokenResponse, err := createToken()
	//tokenResponse, err := c.client.Authentication.CreateToken(context.TODO(), &notion.TokenCreateRequest{
	//	Code:        code,
	//	GrantType:   "authorization_code",
	//	RedirectUri: redirectUri,
	//})
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
}

func createToken() (*notion.TokenCreateResponse, error) {
	return &notion.TokenCreateResponse{
		AccessToken:          "secret_xa0NmgSMu0GdaW5Tgy7hFqEuJk8vXkjtn241BKNezme",
		BotId:                "27967116-f32f-4948-b06f-e5934ad212bf",
		DuplicatedTemplateId: "",
		Owner: struct {
			OwnerType string `json:"type"`
			User      struct {
				Object string `json:"object"`
				Id     string `json:"id"`
			} `json:"user"`
		}{
			OwnerType: "user",
			User: struct {
				Object string `json:"object"`
				Id     string `json:"id"`
			}{
				Object: "user",
				Id:     "4363b48b-527b-435d-a36e-1eae082114b8",
			},
		},
		WorkspaceIcon: "https://s3-us-west-2.amazonaws.com/public.notion-static.com/0c90febe-6fc3-40da-ad9d-66b2d0bb5079/MM.png",
		WorkspaceId:   "5c30e46e-1cbd-43f8-af70-9454521c4809",
		WorkspaceName: "Mattia's Dev Testing Workspace",
	}, nil
}

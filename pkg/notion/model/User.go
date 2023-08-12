package model

type User struct {
	Id string `json:"id"`
	//Firstname string `json:"firstname"`
	//Lastname  string `json:"lastname"`
}

type OAuthToken struct {
	AccessToken   string `json:"access_token"`
	TokenType     string `json:"token_type"`
	BotId         string `json:"bot_id"`
	WorkspaceName string `json:"workspace_name"`
	WorkspaceIcon string `json:"workspace_icon"`
	WorkspaceId   string `json:"workspace_id"`
	Owner         struct {
		Type string `json:"type"`
		User struct {
			Object string `json:"object"`
			Id     string `json:"id"`
		} `json:"user"`
	} `json:"owner"`
	DuplicatedTemplateId string `json:"duplicated_template_id"`
}

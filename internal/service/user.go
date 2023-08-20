package service

import (
	"Notion-Forms/internal/model"
	notionModel "Notion-Forms/pkg/notion/model"
)

func (svc Clients) ConnectIamUserWithNotionUser(iamUserId, redirectUri, code string) error {
	oidcNotionUser, err := svc.authenticateNotionUser(redirectUri, code)
	if err != nil {
		return err
	}

	err = svc.db.dao.ConnectIamUserWithNotionUser(iamUserId, oidcNotionUser.Owner.User.Id, oidcNotionUser.BotId, oidcNotionUser.AccessToken)
	if err != nil {
		return err
	}

	return nil
}

func (svc Clients) authenticateNotionUser(redirectUri, code string) (*notionModel.OAuthToken, error) {
	notionOauthToken, err := svc.notion.Authenticate(redirectUri, code)
	if err != nil {
		return nil, err
	}
	return notionOauthToken, nil
}

func (svc Clients) GetOwnUser(oidcUser model.OidcUser) (*model.GNFUser, error) {
	return svc.db.dao.GetUserDataByIamUserId(oidcUser.Sub)
}

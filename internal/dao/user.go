package dao

import (
	"Notion-Forms/internal/model"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

func (svc *Dao) ConnectIamUserWithNotionUser(iamUserId, notionUserId, notionBotId, notionOauthAccessToken string) error {
	userConnection := model.IamUserNotionUser{
		IamUserId:         iamUserId,
		NotionUserId:      notionUserId,
		NotionBotId:       notionBotId,
		NotionAccessToken: notionOauthAccessToken,
	}
	_, err := svc.engine.Database(svc.dbName).Collection("users").InsertOne(context.Background(), userConnection)
	return err
}

func (svc *Dao) GetUserDataByIamUserId(iamUserId string) (*model.IamUserNotionUser, error) {
	filter := bson.M{"iamuserid": iamUserId}
	result := svc.engine.Database(svc.dbName).Collection("users").FindOne(context.Background(), filter)
	if result == nil {
		return nil, fmt.Errorf("no user was found with this id")
	}

	var resultObject *model.IamUserNotionUser
	err := result.Decode(&resultObject)
	if err != nil {
		return nil, err
	}

	resultObject.NotionAccessToken = ""
	return resultObject, nil
}

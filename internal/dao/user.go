package dao

import (
	"Notion-Forms/internal/model"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"time"
)

func (svc *Dao) ConnectIamUserWithNotionUser(iamUserId, notionUserId, notionBotId, notionOauthAccessToken string) error {
	filter := bson.M{"iamuserid": iamUserId}
	var gnfUser model.GNFUser
	err := svc.engine.Database(svc.dbName).Collection("users").FindOne(context.Background(), filter).Decode(&gnfUser)
	// Error which is not a NotFound error
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	// All Errors except NotFound were filtered -> therefore NotFound-Error
	if err != nil {
		gnfUser.IamUserId = iamUserId
	}

	gnfUser.NotionCredentials.UserId = notionUserId
	gnfUser.NotionCredentials.BotId = notionBotId
	gnfUser.NotionCredentials.AccessToken = notionOauthAccessToken

	// If the record was found and no error occurs
	if err == nil {
		_, err = svc.engine.Database(svc.dbName).Collection("users").UpdateOne(context.Background(), filter, bson.M{"$set": gnfUser})
		if err != nil {
			return err
		}
	}
	// All Errors except NotFound were filtered -> therefore NotFound-Error
	if err != nil {
		_, err = svc.engine.Database(svc.dbName).Collection("users").InsertOne(context.Background(), gnfUser)
		if err != nil {

		}
	}

	return nil
}

func (svc *Dao) ConnectIamUserWithGoogleUser(
	iamUserId,
	oAuthAccessToken,
	oAuthRefreshToken,
	oAuthTokenType string,
	oAuthExpiresIn time.Time,
	oAuthConfig oauth2.Config,
) error {
	filter := bson.M{"iamuserid": iamUserId}
	var gnfUser model.GNFUser
	err := svc.engine.Database(svc.dbName).Collection("users").FindOne(context.Background(), filter).Decode(&gnfUser)
	// Error which is not a NotFound error
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	// All Errors except NotFound were filtered -> therefore NotFound-Error
	if err != nil {
		gnfUser.IamUserId = iamUserId
	}

	gnfUser.GoogleCredentials.AccessToken = oAuthAccessToken
	gnfUser.GoogleCredentials.RefreshToken = oAuthRefreshToken
	gnfUser.GoogleCredentials.Config = oAuthConfig
	gnfUser.GoogleCredentials.TokenType = oAuthTokenType
	gnfUser.GoogleCredentials.ExpiresIn = oAuthExpiresIn

	// If the record was found and no error occurs
	if err == nil {
		_, err = svc.engine.Database(svc.dbName).Collection("users").UpdateOne(context.Background(), filter, bson.M{"$set": gnfUser})
		if err != nil {
			return err
		}
	}
	// All Errors except NotFound were filtered -> therefore NotFound-Error
	if err != nil {
		_, err = svc.engine.Database(svc.dbName).Collection("users").InsertOne(context.Background(), gnfUser)
		if err != nil {

		}
	}

	return nil
}

func (svc *Dao) GetUserDataByIamUserId(iamUserId string) (*model.GNFUser, error) {
	filter := bson.M{"iamuserid": iamUserId}
	result := svc.engine.Database(svc.dbName).Collection("users").FindOne(context.Background(), filter)
	if result == nil {
		return nil, fmt.Errorf("no user was found with this id")
	}

	var resultObject model.GNFUser
	err := result.Decode(&resultObject)
	if err != nil {
		return nil, err
	}

	return &resultObject, nil
}

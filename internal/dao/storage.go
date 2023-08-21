package dao

import (
	"Notion-Forms/internal/model"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (svc *Dao) SetProvider(databaseId string, provider model.StorageProvider) error {
	filter := bson.M{"databaseid": databaseId}
	var formsStorage model.Storage
	err := svc.engine.Database(svc.dbName).Collection("forms").FindOne(context.Background(), filter).Decode(&formsStorage)
	// Error which is not a NotFound error
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	// All Errors except NotFound were filtered -> therefore NotFound-Error
	if err != nil {
		formsStorage.DatabaseId = databaseId
	}

	formsStorage.StorageProvider = provider

	// If the record was found and no error occurs
	if err == nil {
		_, err = svc.engine.Database(svc.dbName).Collection("forms").UpdateOne(context.Background(), filter, bson.M{"$set": formsStorage})
		if err != nil {
			return err
		}
	}
	// All Errors except NotFound were filtered -> therefore NotFound-Error
	if err != nil {
		_, err = svc.engine.Database(svc.dbName).Collection("forms").InsertOne(context.Background(), formsStorage)
		if err != nil {

		}
	}

	return nil
}

func (svc *Dao) SetLocation(databaseId, folderId string) error {
	filter := bson.M{"databaseid": databaseId}
	var formsStorage model.Storage
	err := svc.engine.Database(svc.dbName).Collection("forms").FindOne(context.Background(), filter).Decode(&formsStorage)
	// Error which is not a NotFound error
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	// All Errors except NotFound were filtered -> therefore NotFound-Error
	if err != nil {
		formsStorage.DatabaseId = databaseId
	}

	formsStorage.ParentFolderId = folderId

	// If the record was found and no error occurs
	if err == nil {
		_, err = svc.engine.Database(svc.dbName).Collection("forms").UpdateOne(context.Background(), filter, bson.M{"$set": formsStorage})
		if err != nil {
			return err
		}
	}
	// All Errors except NotFound were filtered -> therefore NotFound-Error
	if err != nil {
		_, err = svc.engine.Database(svc.dbName).Collection("forms").InsertOne(context.Background(), formsStorage)
		if err != nil {

		}
	}

	return nil
}

func (svc *Dao) GetProvider(databaseId string) (*string, error) {
	filter := bson.M{"databaseid": databaseId}
	var formsStorage model.Storage
	err := svc.engine.Database(svc.dbName).Collection("forms").FindOne(context.Background(), filter).Decode(&formsStorage)
	if err != nil {
		return nil, err
	}

	provider := fmt.Sprintf("%s", formsStorage.StorageProvider)
	return &provider, nil
}

func (svc *Dao) GetStorageFolderId(databaseId string) (*string, error) {
	filter := bson.M{"databaseid": databaseId}
	var formsStorage model.Storage
	err := svc.engine.Database(svc.dbName).Collection("forms").FindOne(context.Background(), filter).Decode(&formsStorage)
	if err != nil {
		return nil, err
	}

	return &formsStorage.ParentFolderId, nil
}

package dao

import (
	"Notion-Forms/internal/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func (svc *Dao) CreateForm(updateForm model.Form) error {
	filter := bson.M{"databaseid": updateForm.DatabaseId}
	var form model.Form
	err := svc.engine.Database(svc.dbName).Collection("forms").FindOne(context.Background(), filter).Decode(&form)
	// Error which is not a NotFound error
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	// All Errors except NotFound were filtered -> therefore NotFound-Error
	if err != nil {
		form.DatabaseId = form.DatabaseId
		_, err = svc.engine.Database(svc.dbName).Collection("forms").InsertOne(context.Background(), updateForm)
		if err != nil {
			return err
		}
		return nil
	}

	// If the record was found and no error occurs
	if err == nil {
		form.UrlId = updateForm.UrlId
		form.Password = updateForm.Password
		form.IamUserId = updateForm.IamUserId
		form.Overrides = updateForm.Overrides

		_, err = svc.engine.Database(svc.dbName).Collection("forms").UpdateOne(context.Background(), filter, bson.M{"$set": form})
		if err != nil {
			return err
		}
	}

	return nil
}

func (svc *Dao) ListForms(oidcUserId string) ([]model.Form, error) {
	ctx := context.Background()
	filter := bson.M{"iamuserid": oidcUserId}
	dbFormsResponse, err := svc.engine.Database(svc.dbName).Collection("forms").Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer dbFormsResponse.Close(ctx)

	var forms []model.Form
	for dbFormsResponse.Next(ctx) {
		var form model.Form
		if err = dbFormsResponse.Decode(&form); err != nil {
			log.Fatal(err)
		}
		forms = append(forms, form)
	}

	return forms, nil
}

func (svc *Dao) GetFormById(databaseid, oidcUserId string) (*model.Form, error) {
	ctx := context.Background()
	filter := bson.M{"databaseid": databaseid, "iamuserid": oidcUserId}
	var form model.Form
	err := svc.engine.Database(svc.dbName).Collection("forms").FindOne(ctx, filter).Decode(&form)
	if err != nil {
		return nil, err
	}

	return &form, nil
}

func (svc *Dao) UpdateFormById(updateForm model.Form) (*model.Form, error) {
	filter := bson.M{"databaseid": updateForm.DatabaseId, "iamuserid": updateForm.IamUserId}
	var form model.Form
	err := svc.engine.Database(svc.dbName).Collection("forms").FindOne(context.Background(), filter).Decode(&form)
	// Error which is not a NotFound error
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	// All Errors except NotFound were filtered -> therefore NotFound-Error
	if err != nil {
		form.DatabaseId = form.DatabaseId
		_, err = svc.engine.Database(svc.dbName).Collection("forms").InsertOne(context.Background(), updateForm)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	// If the record was found and no error occurs
	if err == nil {
		form.UrlId = updateForm.UrlId
		form.Password = updateForm.Password
		form.IamUserId = updateForm.IamUserId
		form.Overrides = updateForm.Overrides

		_, err = svc.engine.Database(svc.dbName).Collection("forms").UpdateOne(context.Background(), filter, bson.M{"$set": form})
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (svc *Dao) DeleteFormById(databaseid, oidcUserId string) error {
	filter := bson.M{"databaseid": databaseid, "iamuserid": oidcUserId}
	_, err := svc.engine.Database(svc.dbName).Collection("forms").DeleteOne(context.Background(), filter)
	return err
}

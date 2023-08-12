package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Dao struct {
	engine *mongo.Client
	dbName string
}

func New(engine *mongo.Client, dbName string) *Dao {
	return &Dao{engine: engine, dbName: dbName}
}

func (svc *Dao) GetConnection() error {
	err := svc.engine.Ping(context.TODO(), &readpref.ReadPref{})
	return err
}

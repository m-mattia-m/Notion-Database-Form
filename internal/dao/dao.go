package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Dao struct {
	engine *mongo.Client
}

func New(engine *mongo.Client) *Dao {
	return &Dao{engine: engine}
}

func (dao *Dao) GetConnection() error {
	err := dao.engine.Ping(context.TODO(), &readpref.ReadPref{})
	return err
}

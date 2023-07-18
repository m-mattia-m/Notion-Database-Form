package service

import (
	"Notion-Forms/internal/dao"
	"Notion-Forms/pkg/notion"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	notion *notion.Client
	dao    Dao
}

// Dao ...
type Dao interface {
	GetConnection() error
}

type DAO struct {
	ctx context.Context
	dao *dao.Dao
}

func New(ctx context.Context, engine *mongo.Client, client *notion.Client) Service {
	return Service{
		notion: client,
		dao: DAO{
			ctx: ctx,
			dao: dao.New(engine),
		},
	}
}

func (svc DAO) GetContext() context.Context {
	return svc.ctx
}

func (svc DAO) GetConnection() error {
	return svc.dao.GetConnection()
}

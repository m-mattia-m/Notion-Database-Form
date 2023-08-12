package service

import (
	"Notion-Forms/internal/dao"
	"Notion-Forms/pkg/cache"
	"Notion-Forms/pkg/iam"
	"Notion-Forms/pkg/logging"
	"Notion-Forms/pkg/notion"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Clients struct {
	notion *notion.Client
	cache  *cache.LRUCache
	db     *Db
	logger *logging.Client
	iam    *iam.Client
}

type Service interface {
	GetConnection() error
	ConnectIamUserWithNotionUser(iamUserId string, notionUserId string, code string) error
}

type Db struct {
	ctx context.Context
	dao *dao.Dao
}

func New(ctx context.Context, db *mongo.Client, dbName string, notion *notion.Client, cache *cache.LRUCache, logger *logging.Client, iam *iam.Client) Service {
	return Clients{
		cache:  cache,
		notion: notion,
		db: &Db{
			ctx: ctx,
			dao: dao.New(db, dbName),
		},
		logger: logger,
		iam:    iam,
	}
}

func (svc Db) GetContext() context.Context {
	return svc.ctx
}

func (svc Clients) GetConnection() error {
	return svc.db.dao.GetConnection()
}

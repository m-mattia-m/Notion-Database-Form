package listener

import (
	"Notion-Forms/internal/dao"
	"Notion-Forms/internal/model"
	"Notion-Forms/pkg/cache"
	"Notion-Forms/pkg/iam"
	"Notion-Forms/pkg/logging"
	"Notion-Forms/pkg/notion"
	"encoding/json"
	"github.com/asaskevich/EventBus"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var (
	Bus EventBus.Bus
)

type Listener struct {
	notion *notion.Client
	cache  *cache.LRUCache
	dao    *dao.Dao
	logger *logging.Client
	iam    *iam.Client
}

func New(db *mongo.Client, dbName string, notion *notion.Client, cache *cache.LRUCache, logger *logging.Client, iam *iam.Client) Listener {
	Bus = EventBus.New()
	return Listener{
		cache:  cache,
		notion: notion,
		dao:    dao.New(db, dbName),
		logger: logger,
		iam:    iam,
	}
}

func (listener Listener) StartListener() error {
	err := Bus.Subscribe("notion:update-database", listener.updateNotionDatabase)
	return err
}

func (listener Listener) updateNotionDatabase(id string) {
	database, err := listener.notion.GetDatabase(id)
	if err != nil {
		listener.logger.Error("notion", "GetDatabase", logging.Message{
			Description: "failed to update notion database information",
			Detail:      err,
			ErrorLevel:  logging.High,
			LogId:       uuid.New().String(),
		})
	}

	databaseString, err := json.Marshal(model.StoreDatabaseObject{
		Expiration:     time.Now(),
		RelevanceScore: 1,
		Object:         database,
	})
	err = listener.cache.Set(string(database.ID), string(databaseString))
	if err != nil {
		listener.logger.Error("cache", "Set", logging.Message{
			Description: "failed to update cache-object by their id",
			Detail:      err,
			ErrorLevel:  logging.High,
			LogId:       uuid.New().String(),
		})
	}
}

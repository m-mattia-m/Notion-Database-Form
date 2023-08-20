package service

import (
	"Notion-Forms/internal/dao"
	"Notion-Forms/internal/model"
	"Notion-Forms/pkg/cache"
	"Notion-Forms/pkg/iam"
	"Notion-Forms/pkg/logging"
	"Notion-Forms/pkg/notion"
	notionModel "Notion-Forms/pkg/notion/model"
	googleDrive "Notion-Forms/pkg/storage/google-drive"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jomei/notionapi"
)

type Clients struct {
	notion  *notion.Client
	cache   *cache.LRUCache
	db      *Db
	logger  *logging.Client
	iam     *iam.Client
	storage *StorageClients
}

type StorageClients struct {
	googleDrive *googleDrive.Client
}

type Service interface {
	GetConnection() error
	SetNotionClient(client *notion.Client) Clients
	ConnectIamUserWithNotionUser(iamUserId, redirectUri, code string) error
	SetAbortResponse(c *gin.Context, targetServiceName string, method string, message string, err error) string
	SetAbortWithoutResponse(targetServiceName string, method string, message string, err error) string
	GetOwnUser(oidcUser model.OidcUser) (*model.GNFUser, error)
	GetMe() (*notionModel.User, error)
	ListDatabases() ([]*notionapi.Database, error)
	GetDatabase(id string) (notionapi.Database, error)
	CreateRecord(databaseId, userId string, requests []notionModel.RecordRequest) (notionapi.Page, error)
	ListSelectOptions(databaseId string, selectName string) ([]notionModel.Select, error)
	ListAllSelectOptions(databaseId string) ([]notionModel.Select, error)
	GetPage(id string) (notionapi.Page, error)
	ListPages() ([]*notionapi.Page, error)
}

type Db struct {
	ctx context.Context
	dao *dao.Dao
}

func New(
	ctx context.Context,
	db *mongo.Client,
	dbName string,
	cache *cache.LRUCache,
	logger *logging.Client,
	iam *iam.Client,
	googleDrive *googleDrive.Client,
) Service {
	return Clients{
		notion: nil,
		cache:  cache,
		db: &Db{
			ctx: ctx,
			dao: dao.New(db, dbName),
		},
		logger: logger,
		iam:    iam,
		storage: &StorageClients{
			googleDrive: googleDrive,
		},
	}
}

func (svc Db) GetContext() context.Context {
	return svc.ctx
}

func (svc Clients) GetConnection() error {
	return svc.db.dao.GetConnection()
}

func (svc Clients) SetNotionClient(client *notion.Client) Clients {
	svc.notion = client
	return svc
}

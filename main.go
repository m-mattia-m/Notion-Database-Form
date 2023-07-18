package main

import (
	"Notion-Forms/internal/api"
	apiV1 "Notion-Forms/internal/api/v1"
	"Notion-Forms/internal/service"
	"Notion-Forms/pkg/notion"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strconv"
)

const (
	maxPageSizeInit     = 500
	defaultPageSizeInit = 10
)

// @title Generated Notion Forms
// @version 1.0
// @description this is the api for generated Notion forms

// @contact.name API Support
// @contact.email develop@mattiamueggler.ch

// @host      localhost:8081
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

// @tokenUrl https://zitadel.upcraft.li/oauth/v2/token
// @authorizationUrl https://zitadel.upcraft.li/oauth/v2/authorize
// @scope.openid Default Grants
// @scope.profile Default Grants
// @scope.email Default Grants
// @scope.roles Default Grants
// @scope.notion-database-form-notion-database-form-swagger-local

func main() {
	err := initConfig()
	if err != nil {
		log.Fatal("error when init the config: " + err.Error())
	}

	dbClient, err := createDbConnection()
	if err != nil {
		log.Fatal("error when starting db: " + err.Error())
	}

	notionClient, err := notion.New()
	if err != nil {
		log.Fatal("error when starting notion-client: " + err.Error())
	}

	svc := service.New(context.Background(), dbClient, notionClient)

	apiConfig, err := getApiConfig()
	if err != nil {
		log.Fatal("error when reading the api-config: " + err.Error())
	}

	apiClient := api.New(svc, *apiConfig)
	err = api.Router(apiClient)
	if err != nil {
		log.Fatal("error when starting router: " + err.Error())
	}
}

func createDbConnection() (*mongo.Client, error) {
	mongoHost, found := os.LookupEnv("MONGO_HOST")
	if !found {
		return nil, fmt.Errorf("env-variable 'MONGO_ROOT_PASSWORD' not found")
	}
	mongoPort, found := os.LookupEnv("MONGO_PORT")
	if !found {
		return nil, fmt.Errorf("env-variable 'MONGO_ROOT_PASSWORD' not found")
	}
	mongoUser, found := os.LookupEnv("MONGO_ROOT_USER")
	if !found {
		return nil, fmt.Errorf("env-variable 'MONGO_ROOT_USER' not found")
	}
	mongoPassword, found := os.LookupEnv("MONGO_ROOT_PASSWORD")
	if !found {
		return nil, fmt.Errorf("env-variable 'MONGO_ROOT_PASSWORD' not found")
	}

	var uri = fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	credentials := options.Credential{
		AuthMechanism: "PLAIN",
		Username:      mongoUser,
		Password:      mongoPassword,
	}
	opts := options.Client().
		ApplyURI(uri).SetAuth(credentials).
		SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func initConfig() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}

func getApiConfig() (*apiV1.ApiConfig, error) {
	runMode := os.Getenv("APP_ENV")
	if runMode == "" {
		return nil, fmt.Errorf("failed to get env-variable: 'APP_ENV'")
	}
	frontendHost := os.Getenv("FRONTEND_HOST_URL")
	if frontendHost == "" {
		return nil, fmt.Errorf("failed to get env-variable: 'FRONTEND_HOST_URL'")
	}
	host := os.Getenv("APP_HOST")
	if host == "" {
		return nil, fmt.Errorf("failed to get env-variable: 'APP_HOST'")
	}
	port, err := strconv.ParseInt(os.Getenv("APP_PORT"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to get env-variable: 'APP_PORT'")
	}
	domain := os.Getenv("APP_DOMAIN")
	if domain == "" {
		return nil, fmt.Errorf("failed to get env-variable: 'APP_DOMAIN'")
	}
	schemas := os.Getenv("APP_SCHEMES")
	if schemas == "" {
		return nil, fmt.Errorf("failed to get env-variable: 'APP_SCHEMES'")
	}

	maxPageSize, err := strconv.ParseInt(os.Getenv("APP_MAX_PAGE_SIZE"), 10, 64)
	if err != nil {
		maxPageSize = maxPageSizeInit
		log.Println(fmt.Sprintf("'APP_MAX_PAGE_SIZE' is not set therefore we use the default value witch is '%d'", maxPageSize))
	}
	defaultPageSize, err := strconv.ParseInt(os.Getenv("APP_DEFAULT_PAGE_SIZE"), 10, 64)
	if err != nil {
		defaultPageSize = defaultPageSizeInit
		log.Println(fmt.Sprintf("'APP_DEFAULT_PAGE_SIZE' is not set therefore we use the default value witch is '%d'", defaultPageSize))

	}

	return &apiV1.ApiConfig{
		RunMode:         runMode,
		FrontendHost:    frontendHost,
		DefaultPageSize: int(defaultPageSize),
		MaxPageSize:     int(maxPageSize),
		Port:            int(port),
		Host:            host,
		Domain:          domain,
		Schemas:         schemas,
	}, nil
}

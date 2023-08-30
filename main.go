package main

import (
	"Notion-Forms/internal/api"
	apiV1 "Notion-Forms/internal/api/v1"
	"Notion-Forms/internal/service"
	"Notion-Forms/pkg/cache"
	"Notion-Forms/pkg/iam"
	"Notion-Forms/pkg/logging"
	googleDrive "Notion-Forms/pkg/storage/google-drive"
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
// @contact.email develop@generated-notion-forms.com

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
// @scope.notion-database-form-notion-database-form-swagger-local  Default Grants

func main() {
	err := initConfig()
	if err != nil {
		log.Fatal("error when init the config: " + err.Error())
	}

	dbClient, dbName, err := createDbConnection()
	if err != nil {
		log.Fatal("error when starting db: " + err.Error())
	}

	cacheClient, err := createLRUCache()
	if err != nil {
		log.Fatal("error when starting cache: " + err.Error())
	}

	loggingClient, err := createLoggingClient()
	if err != nil {
		log.Fatal("error when starting logging-client: " + err.Error())
	}

	iamClient, err := createIamClient()
	if err != nil {
		log.Fatal("error when starting iam-client: " + err.Error())
	}

	googleDriveStorage, err := createGoogleDriveClient()
	if err != nil {
		log.Fatal("error when starting google-drive-storage-client: " + err.Error())
	}

	svc := service.New(context.Background(), dbClient, dbName, cacheClient, loggingClient, iamClient, googleDriveStorage)

	apiConfig, err := getApiConfig()
	if err != nil {
		log.Fatal("error when reading the api-config: " + err.Error())
	}

	err = api.Router(svc, *apiConfig)
	if err != nil {
		log.Fatal("error when starting router: " + err.Error())
	}
}

func createDbConnection() (*mongo.Client, string, error) {
	mongoHost, found := os.LookupEnv("MONGO_HOST")
	if !found {
		return nil, "", fmt.Errorf("env-variable 'MONGO_HOST' not found")
	}
	mongoPort, found := os.LookupEnv("MONGO_PORT")
	if !found {
		return nil, "", fmt.Errorf("env-variable 'MONGO_PORT' not found")
	}
	mongoUser, found := os.LookupEnv("MONGO_ROOT_USER")
	if !found {
		return nil, "", fmt.Errorf("env-variable 'MONGO_ROOT_USER' not found")
	}
	mongoPassword, found := os.LookupEnv("MONGO_ROOT_PASSWORD")
	if !found {
		return nil, "", fmt.Errorf("env-variable 'MONGO_ROOT_PASSWORD' not found")
	}
	mongoDatabaseName, found := os.LookupEnv("MONGO_DATABASE_NAME")
	if !found {
		return nil, "", fmt.Errorf("env-variable 'MONGO_DATABASE_NAME' not found")
	}

	var uri = fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	credentials := options.Credential{
		AuthMechanism: "SCRAM-SHA-256",
		Username:      mongoUser,
		Password:      mongoPassword,
	}
	opts := options.Client().
		ApplyURI(uri).SetAuth(credentials).
		SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, "", err
	}

	return client, mongoDatabaseName, nil
}

func createLRUCache() (*cache.LRUCache, error) {
	redisHost, found := os.LookupEnv("REDIS_HOST")
	if !found {
		return nil, fmt.Errorf("env-variable 'REDIS_HOST' not found")
	}
	redisPort, found := os.LookupEnv("REDIS_PORT")
	if !found {
		return nil, fmt.Errorf("env-variable 'REDIS_PORT' not found")
	}
	redisPassword, found := os.LookupEnv("REDIS_PASSWORD")
	if !found {
		return nil, fmt.Errorf("env-variable 'REDIS_PASSWORD' not found")
	}
	cacheCapacity, found := os.LookupEnv("CACHE_CAPACITY")
	if !found {
		return nil, fmt.Errorf("env-variable 'REDIS_PASSWORD' not found")
	}
	cacheCapacityInt, err := strconv.ParseInt(cacheCapacity, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to convert cache-capacity to int")
	}

	lruCache := cache.NewLRUCache(fmt.Sprintf("%s:%s", redisHost, redisPort), redisPassword, int(cacheCapacityInt))
	return lruCache, nil
}

func createLoggingClient() (*logging.Client, error) {
	appEnv, found := os.LookupEnv("APP_ENV")
	if !found {
		return nil, fmt.Errorf("env-variable 'APP_ENV' not found")
	}

	envDns, found := os.LookupEnv("LOGGING_DNS")
	if !found {
		return nil, fmt.Errorf("env-variable 'LOGGING_DNS' not found")
	}

	enableLogging, found := os.LookupEnv("ENABLE_LOGGING")
	if !found {
		return nil, fmt.Errorf("env-variable 'ENABLE_LOGGING' not found")
	}

	enableExternalLogging, found := os.LookupEnv("EXTERNAL_ENABLE_LOGGING")
	if !found {
		return nil, fmt.Errorf("env-variable 'EXTERNAL_ENABLE_LOGGING' not found")
	}

	return logging.New(appEnv, envDns, enableLogging, enableExternalLogging)
}

func createIamClient() (*iam.Client, error) {
	envIssuer, found := os.LookupEnv("ZITADEL_ISSUER")
	if !found {
		return nil, fmt.Errorf("env-variable 'ZITADEL_ISSUER' not found")
	}
	envApi, found := os.LookupEnv("ZITADEL_API")
	if !found {
		return nil, fmt.Errorf("env-variable 'ZITADEL_API' not found")
	}
	organizationId, found := os.LookupEnv("ZITADEL_ORGANIZATION_ID")
	if !found {
		return nil, fmt.Errorf("env-variable 'ZITADEL_ORGANIZATION_ID' not found")
	}
	organizationName, found := os.LookupEnv("ZITADEL_ORGANIZATION_NAME")
	if !found {
		return nil, fmt.Errorf("env-variable 'ZITADEL_ORGANIZATION_NAME' not found")
	}
	projectId, found := os.LookupEnv("ZITADEL_PROJECT_ID")
	if !found {
		return nil, fmt.Errorf("env-variable 'ZITADEL_PROJECT_ID' not found")
	}
	projectName, found := os.LookupEnv("ZITADEL_PROJECT_NAME")
	if !found {
		return nil, fmt.Errorf("env-variable 'ZITADEL_PROJECT_NAME' not found")
	}

	return iam.New(envIssuer, envApi, organizationId, organizationName, projectId, projectName)

}

func createGoogleDriveClient() (*googleDrive.Client, error) {
	googleClientId, found := os.LookupEnv("GOOGLE_OAUTH_CLIENT_ID")
	if !found {
		return nil, fmt.Errorf("env-variable 'GOOGLE_OAUTH_CLIENT_ID' not found")
	}
	googleClientSecret, found := os.LookupEnv("GOOGLE_OAUTH_CLIENT_SECRET")
	if !found {
		return nil, fmt.Errorf("env-variable 'GOOGLE_OAUTH_CLIENT_SECRET' not found")
	}
	googleRedirectUrl, found := os.LookupEnv("GOOGLE_OAUTH_REDIRECT_URL")
	if !found {
		return nil, fmt.Errorf("env-variable 'GOOGLE_OAUTH_REDIRECT_URL' not found")
	}

	googleDriveClient, err := googleDrive.New(googleDrive.GoogleOauthConfig{
		ClientId:     googleClientId,
		ClientSecret: googleClientSecret,
		RedirectUri:  googleRedirectUrl,
	})
	if err != nil {
		return nil, err
	}

	return googleDriveClient, nil
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
	oidcAuthority := os.Getenv("OIDC_AUTHORITY")
	if runMode == "" {
		return nil, fmt.Errorf("failed to get env-variable: 'OIDC_AUTHORITY'")
	}

	oidcClientId := os.Getenv("OIDC_CLIENT_ID")
	if runMode == "" {
		return nil, fmt.Errorf("failed to get env-variable: 'OIDC_CLIENT_ID'")
	}

	notionClientId := os.Getenv("NOTION_CLIENT_ID")
	if runMode == "" {
		return nil, fmt.Errorf("failed to get env-variable: 'NOTION_CLIENT_ID'")
	}
	notionClientSecret := os.Getenv("NOTION_CLIENT_SECRET")
	if runMode == "" {
		return nil, fmt.Errorf("failed to get env-variable: 'NOTION_CLIENT_SECRET'")
	}
	notionRedirectUri := os.Getenv("NOTION_REDIRECT_URI")
	if runMode == "" {
		return nil, fmt.Errorf("failed to get env-variable: 'NOTION_REDIRECT_URI'")
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
		RunMode:            runMode,
		FrontendHost:       frontendHost,
		DefaultPageSize:    int(defaultPageSize),
		MaxPageSize:        int(maxPageSize),
		Port:               int(port),
		Host:               host,
		Domain:             domain,
		Schemas:            schemas,
		OidcAuthority:      oidcAuthority,
		OidcClientId:       oidcClientId,
		NotionClientId:     notionClientId,
		NotionClientSecret: notionClientSecret,
		NotionRedirectUri:  notionRedirectUri,
	}, nil
}

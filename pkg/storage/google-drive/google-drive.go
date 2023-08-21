package google_drive

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

const (
	GoogleOAuthGrantType = "authorization_code"
)

type GoogleOauthConfig struct {
	ClientId     string
	ClientSecret string
	RedirectUri  string
}

type GoogleOauthTokenConfig struct {
	Config       oauth2.Config `json:"config"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    time.Time     `json:"expires_in"`
	TokenType    string        `json:"token_type"`
	//Scope        string        `json:"scope"`
	//IdToken      string        `json:"id_token"`
}

type Client struct {
	ClientId     string
	ClientSecret string
	RedirectUri  string
	GrandType    string
}

func New(config GoogleOauthConfig) (*Client, error) {
	return &Client{
		ClientId:     config.ClientId,
		ClientSecret: config.ClientSecret,
		RedirectUri:  config.RedirectUri,
		GrandType:    GoogleOAuthGrantType,
	}, nil
}

func (c *Client) Authenticate(code string) (*GoogleOauthTokenConfig, error) {
	config := oauth2.Config{
		ClientID:     c.ClientId,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectUri,
		Scopes:       []string{"https://www.googleapis.com/auth/drive"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	return &GoogleOauthTokenConfig{
		Config:       config,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.Expiry,
		TokenType:    token.TokenType,
	}, nil
}

func (c *Client) UploadFile(oauthTokenConfig GoogleOauthTokenConfig, parentFolderId string, file *multipart.FileHeader) (*string, error) {
	svc, err := drive.NewService(context.Background(), option.WithHTTPClient(
		oauthTokenConfig.Config.Client(
			context.Background(),
			&oauth2.Token{
				AccessToken:  oauthTokenConfig.AccessToken,
				TokenType:    oauthTokenConfig.TokenType,
				RefreshToken: oauthTokenConfig.RefreshToken,
				Expiry:       oauthTokenConfig.ExpiresIn,
			},
		)),
	)
	if err != nil {
		return nil, err
	}

	fileExtension := filepath.Ext(file.Filename)
	fileName := strings.TrimSuffix(file.Filename, fileExtension)
	uploadFile := drive.File{
		Name:    fmt.Sprintf("%s_%s%s", fileName, uuid.New().String(), fileExtension),
		Parents: []string{parentFolderId},
	}
	fileContent, err := file.Open()
	if err != nil {
		return nil, err
	}

	fileResult, err := svc.Files.Create(&uploadFile).Media(fileContent).Do()
	if err != nil {
		return nil, err
	}

	uploadedFile, err := svc.Files.Get(fileResult.Id).Fields("webViewLink").Do()
	if err != nil {
		return nil, err
	}
	return &uploadedFile.WebViewLink, nil
}

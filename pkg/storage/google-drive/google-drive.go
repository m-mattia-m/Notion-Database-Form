package google_drive

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

const (
	GoogleOAuthAuthorizeUrl = "https://oauth2.googleapis.com/token"
	GoogleOAuthGrantType    = "authorization_code"
)

type GoogleOauthConfig struct {
	ClientId     string
	ClientSecret string
	RedirectUri  string
}

type GoogleOauthToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	IdToken     string `json:"id_token"`
}

type Client struct {
	httpQueries url.Values
}

func New(config GoogleOauthConfig) (*Client, error) {

	values := url.Values{}
	values.Add("grant_type", GoogleOAuthGrantType)
	values.Add("client_id", config.ClientId)
	values.Add("client_secret", config.ClientSecret)
	values.Add("redirect_uri", config.RedirectUri)

	return &Client{
		httpQueries: values,
	}, nil
}

func (c *Client) Authenticate(code string) (*GoogleOauthToken, error) {
	c.httpQueries.Add("code", code)
	query := c.httpQueries.Encode()

	req, err := http.NewRequest("POST", GoogleOAuthAuthorizeUrl, bytes.NewBufferString(query))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the response-code was not successfull (200)")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)
	if err != nil {
		return nil, err
	}

	var googleOauthToken GoogleOauthToken
	if err := json.Unmarshal(resBody.Bytes(), &googleOauthToken); err != nil {
		return nil, err
	}

	// TODO: Ich bekomme den ID_token nicht zurück und brauche diesen aber um ein JWT-Token zu generieren -> siehe anleitung

	return &googleOauthToken, nil
}

func (c *Client) UploadFile(path string, file *multipart.FileHeader) error {
	ctx := context.Background()
	service, err := drive.NewService(ctx, option.WithAPIKey(""))
	if err != nil {
		return err
	}

	driveService := drive.NewDrivesService(service)

	fileList, err := driveService.List().Do()
	if err != nil {
		return err
	}
	fmt.Println(fileList)

	return nil
}

package helper

import (
	"Notion-Forms/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"strings"
)

func GetUser(c *gin.Context, cfg model.OidcConfig) (model.OidcUser, error) {
	claims, err := getClaims(c, cfg)
	if err != nil {
		return model.OidcUser{}, err
	}

	roles, err := getRoles(c, cfg)
	if err != nil {
		return model.OidcUser{}, err
	}

	return model.OidcUser{
		Email:             claims.Email,
		EmailVerified:     claims.EmailVerified,
		FamilyName:        claims.FamilyName,
		Gender:            claims.Gender,
		GivenName:         claims.GivenName,
		Locale:            claims.Locale,
		Name:              claims.Name,
		Nickname:          claims.Nickname,
		PreferredUsername: claims.PreferredUsername,
		Sub:               claims.Sub,
		UpdatedAt:         claims.UpdatedAt,
		Roles:             roles,
	}, nil
}

func getClaims(c *gin.Context, cfg model.OidcConfig) (*model.BearerClaims, error) {

	bearer := GetBearer(c)
	if bearer == "" {
		return nil, errors.New("failed to get bearer")
	}

	provider, err := oidc.NewProvider(c.Request.Context(), cfg.OidcAuthority)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't create new provider -> %s", err))
	}

	insecureSkipSignatureCheck := cfg.AppEnv == "DEV"
	var verifier = provider.Verifier(&oidc.Config{ClientID: cfg.OidcClientId, InsecureSkipSignatureCheck: insecureSkipSignatureCheck})

	IDToken, err := verifier.Verify(c.Request.Context(), bearer)
	if err != nil {
		return nil, err
	}

	var claims model.BearerClaims
	if err := IDToken.Claims(&claims); err != nil {
		return nil, errors.New(fmt.Sprintf("can't get custom claims -> %s", err))
	}

	return &claims, nil
}

func getRoles(c *gin.Context, cfg model.OidcConfig) ([]string, error) {
	claims, err := getClaims(c, cfg)
	if err != nil {
		return nil, err
	}

	var rolesMap map[string]interface{}
	roleBytes, err := json.Marshal(claims.UrnZitadelIamOrgProjectRoles)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't marshal roles to []bytes -> %s", err))
	}
	if err := json.Unmarshal(roleBytes, &rolesMap); err != nil {
		return nil, errors.New(fmt.Sprintf("can't unmarshal roles-[]bytes to map -> %s", err))
	}

	var userRoles []string
	for role := range rolesMap {
		userRoles = append(userRoles, role)
	}

	return userRoles, nil
}

func GetBearer(c *gin.Context) string {
	authToken := c.Request.Header.Get("Authorization")
	if authToken == "" {
		return ""
	}
	authTokenSections := strings.Split(authToken, " ")

	if len(authTokenSections) != 2 {
		return ""
	}
	return authTokenSections[1]
}

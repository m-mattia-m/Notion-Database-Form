package api

import (
	v1 "Notion-Forms/internal/api/v1"
	"Notion-Forms/internal/helper"
	"context"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {

	config, found := c.Get("config")
	if !found {
		helper.SetUnauthorizedResponse(c, "Unauthorized")
		return
	}

	configObject := config.(v1.ApiConfig)

	oidcClientId := configObject.OidcClientId
	oidcAuthority := configObject.OidcAuthority
	runMode := configObject.RunMode

	bearer := helper.GetBearer(c)
	if bearer == "" {
		helper.SetUnauthorizedResponse(c, "Unauthorized")
		return
	}

	_, err := validateToken(c.Request.Context(), bearer, oidcClientId, oidcAuthority, runMode)
	if err != nil {
		helper.SetUnauthorizedResponse(c, "Unauthorized")
		return
	}

}

func validateToken(ctx context.Context, bearer, oidcClientId, oidcAuthority, runMode string) (*oidc.Provider, error) {
	provider, err := oidc.NewProvider(ctx, oidcAuthority)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't create new provider -> %s", err))
	}

	insecureSkipSignatureCheck := runMode == "DEV"
	var verifier = provider.Verifier(&oidc.Config{ClientID: oidcClientId, InsecureSkipSignatureCheck: insecureSkipSignatureCheck})

	_, err = verifier.Verify(ctx, bearer)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't verify bearer -> %s", err))
	}

	return provider, nil
}

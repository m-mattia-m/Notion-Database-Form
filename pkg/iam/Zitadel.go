package iam

import (
	"Templates/pkg/helper"
	"Templates/pkg/logging"
	"context"
	"errors"
	"flag"
	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/zitadel-go/v2/pkg/client/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel"
)

var (
	issuer *string
	api    *string
)

func init() {

	envIssuer, err := helper.GetEnv("ZITADEL_ISSUER")
	if err != nil {
		logging.Error("iam/Zitadel", "init", logging.Message{
			Description: "can't get env-var ZITADEL_ISSUER: ",
			Detail:      err,
		})
	}
	envapi, err := helper.GetEnv("ZITADEL_API")
	if err != nil {
		logging.Error("iam/Zitadel", "init", logging.Message{
			Description: "can't get env-var ZITADEL_API: ",
			Detail:      err,
		})
	}

	issuer = flag.String("issuer", envIssuer, "")
	api = flag.String("api", envapi, "")
}

func CreateZitadelUser() error {

	flag.Parse()
	client, err := management.NewClient(
		*issuer,
		*api,
		[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
	)
	if err != nil {
		return errors.New("could not create zitadel-client: '" + err.Error() + "'")
	}
	defer func() {
		err := client.Connection.Close()
		if err != nil {
			logging.Error("iam/zitadel", "CreateZitadelUser", logging.Message{
				Description: "could not close grpc connection:",
				Detail:      err,
			})
		}
	}()
	ctx := context.Background()

	_ = ctx
	return nil
}

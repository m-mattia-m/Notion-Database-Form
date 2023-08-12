package iam

import (
	"context"
	"flag"
	"fmt"
	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/zitadel-go/v2/pkg/client/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel"
	pb "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
)

type Client struct {
	issuer           *string
	api              *string
	organizationId   *string
	organizationName *string
	projectId        *string
	projectName      *string
	zitadel          *management.Client
}

func New(envIssuer, envApi, organizationId, organizationName, projectId, projectName string) (*Client, error) {
	var client Client
	client.issuer = flag.String("issuer", envIssuer, "")
	client.api = flag.String("api", envApi, "")
	client.organizationId = &organizationId
	client.organizationName = &organizationName
	client.projectId = &projectId
	client.projectName = &projectName

	flag.Parse()
	zitadelClient, err := management.NewClient(
		*client.issuer,
		*client.api,
		[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
	)
	if err != nil {
		return nil, fmt.Errorf("could not create zitadel-client: %s", err)
	}

	// currently it looks like it is not needed.
	// if there is a timeout of the session, and Zitadel rejects the connection after a certain time,
	// you have to check it again.
	// #------------#
	//var zitadelConnectionCloseError = make(chan error)
	//defer func() {
	//	err := zitadelClient.Connection.Close()
	//	if err != nil {
	//		zitadelConnectionCloseError <- err
	//		return
	//	}
	//	zitadelConnectionCloseError <- nil
	//}()
	//if <-zitadelConnectionCloseError != nil {
	//	return nil, <-zitadelConnectionCloseError
	//}

	client.zitadel = zitadelClient

	return &client, nil
}

// AddRoleToUser
// zitadelId: the ID can be fetched from the payload of the bearer token from the attribute `sub`
// role: must be specified manually in the frontend which is 1:1 the same as in Zitadel
func (c Client) AddRoleToUser(zitadelId string, role string) error {
	ctx := context.Background()
	_, err := c.zitadel.AddUserGrant(ctx, &pb.AddUserGrantRequest{
		UserId:    zitadelId,
		ProjectId: *c.projectId,
		RoleKeys:  []string{role},
	})
	if err != nil {
		return fmt.Errorf("add-role-to-user'-call failed: %s", err)
	}
	return nil
}

// AddNotionIdToUser
// zitadelId: the ID can be fetched from the payload of the bearer token from the attribute `sub`
// notionId: the ID from the notion user
func (c Client) AddNotionIdToUser(zitadelId string, notionId string) error {
	ctx := context.Background()
	_, err := c.zitadel.SetUserMetadata(ctx, &pb.SetUserMetadataRequest{
		Id:    zitadelId,
		Key:   "notionId",
		Value: []byte(notionId),
	})
	if err != nil {
		return fmt.Errorf("add-notion-id-to-user-meta'-call failed: %s", err)
	}
	return nil
}

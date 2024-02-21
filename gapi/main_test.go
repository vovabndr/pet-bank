package gapi

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	db "pet-bank/db/sqlc"
	"pet-bank/token"
	"pet-bank/utils"
	"pet-bank/worker"
	"testing"
	"time"
)

func NewTestServer(t *testing.T, store db.Store, distributor worker.TaskDistributor) *Server {
	config := utils.Config{
		TokenSymmetricKey:   utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, distributor)
	require.NoError(t, err)

	return server
}

func newContextWithBearerToken(t *testing.T, username string, tokenMaker token.Maker, duration time.Duration) context.Context {
	accessToken, _, err := tokenMaker.CreateToken(username, duration)
	bearerToken := fmt.Sprintf("%s %s", authorizationTypeBearer, accessToken)
	require.NoError(t, err)
	md := metadata.MD{authorizationHeader: []string{bearerToken}}
	return metadata.NewIncomingContext(context.Background(), md)
}

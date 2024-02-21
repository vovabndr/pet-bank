package gapi

import (
	"github.com/stretchr/testify/require"
	db "pet-bank/db/sqlc"
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

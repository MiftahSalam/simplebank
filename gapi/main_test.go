package gapi

import (
	"context"
	"fmt"
	db "simplebank/db/sqlc"
	"simplebank/token"
	"simplebank/util"
	"simplebank/worker"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func newTestServer(t *testing.T, store db.Store, worker worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymetricKey:     util.RandomString(32),
		TokenExpiredDuration: time.Minute,
	}

	server, err := NewServer(config, store, worker)
	require.NoError(t, err)

	return server
}

func newCtxWithBearerToken(t *testing.T, tokenMaker token.TokenManager, username string, duration time.Duration) context.Context {
	token, _, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	bearerToken := fmt.Sprintf("%s %s", authHeaderTypeBearer, token)
	md := metadata.MD{
		authHeaderKey: []string{
			bearerToken,
		},
	}

	return metadata.NewIncomingContext(context.Background(), md)
}

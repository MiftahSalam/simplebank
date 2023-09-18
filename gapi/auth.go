package gapi

import (
	"context"
	"fmt"
	"simplebank/token"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	authHeaderKey        = "authorization"
	authHeaderTypeBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	m, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := m.Get(authHeaderKey)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing auth header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authHeaderTypeBearer {
		return nil, fmt.Errorf("unsupported authorization type")
	}

	access_token := fields[1]
	payload, err := server.tokenManager.VerifyToken(access_token)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %s", err)
	}

	return payload, nil
}

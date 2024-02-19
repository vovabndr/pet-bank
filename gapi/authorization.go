package gapi

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
	"pet-bank/token"
	"strings"
)

const (
	authorizationHeader     = "authorization"
	authorizationTypeBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)

	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != authorizationTypeBearer {
		return nil, errors.Errorf("unsupported authorization type: %s", authorizationTypeBearer)
	}

	accessToken := fields[1]

	return server.tokenMaker.Verify(accessToken)
}

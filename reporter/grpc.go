package reporter

import (
	"context"

	"github.com/forbole/juno/v3/node/remote"
	"google.golang.org/grpc/metadata"
)

// GetRequestContext returns the context to be used for gRPC calls
func GetRequestContext(height int64, authorizationToken string) context.Context {
	return AuthorizedContext(
		remote.GetHeightRequestContext(context.Background(), height),
		authorizationToken,
	)
}

// AuthorizedContext returns an authorized context
func AuthorizedContext(ctx context.Context, authorizationToken string) context.Context {
	if authorizationToken == "" {
		return ctx
	}

	return metadata.AppendToOutgoingContext(ctx, "Authorization", authorizationToken)
}

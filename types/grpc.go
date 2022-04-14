package types

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/forbole/juno/v3/node/remote"
	"google.golang.org/grpc/metadata"
)

// ParseGRPCAddress parses the given gRPC address to get everything needed
func ParseGRPCAddress(grpcAddress string) (address string, insecure bool, headers map[string]string) {
	q, err := url.Parse(grpcAddress)
	if err != nil {
		panic(err)
	}

	headers = make(map[string]string, len(q.Query()))
	for key, values := range q.Query() {
		headers[key] = values[0]
	}

	return fmt.Sprintf("%s:%s", q.Hostname(), q.Port()), !strings.Contains(grpcAddress, "https://"), headers
}

// GetRequestContext returns the context to be used for gRPC calls
func GetRequestContext(height int64, headers map[string]string) context.Context {
	return ContextWithHeaders(
		remote.GetHeightRequestContext(context.Background(), height),
		headers,
	)
}

// ContextWithHeaders returns an authorized context
func ContextWithHeaders(ctx context.Context, headers map[string]string) context.Context {
	var values []string
	for key, value := range headers {
		values = append(values, key)
		values = append(values, value)
	}
	return metadata.AppendToOutgoingContext(ctx, values...)
}

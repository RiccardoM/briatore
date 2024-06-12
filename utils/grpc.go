package utils

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"google.golang.org/grpc/metadata"
)

// ParseAddressHeaders parses the given gRPC address to get everything needed
func ParseAddressHeaders(grpcAddress string) (address string, headers map[string]string) {
	q, err := url.Parse(grpcAddress)
	if err != nil {
		panic(err)
	}

	headers = make(map[string]string, len(q.Query()))
	for key, values := range q.Query() {
		headers[key] = values[0]
	}

	return fmt.Sprintf("%s://%s", q.Scheme, q.Host), headers
}

// GetHeightRequestContext adds the height to the context for querying the state at a given height
func GetHeightRequestContext(context context.Context, height int64) context.Context {
	return metadata.AppendToOutgoingContext(
		context,
		grpctypes.GRPCBlockHeightHeader,
		strconv.FormatInt(height, 10),
	)
}

// GetRequestContext returns the context to be used for gRPC calls
func GetRequestContext(height int64, headers map[string]string) context.Context {
	return ContextWithHeaders(
		GetHeightRequestContext(context.Background(), height),
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

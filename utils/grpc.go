package utils

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	HTTPProtocols = regexp.MustCompile("https?://")
)

// ParseGRPCAddress parses the given gRPC address to get everything needed
func ParseGRPCAddress(grpcAddress string) (address string, headers map[string]string) {
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

// CreateGrpcConnection creates a new gRPC client connection from the given configuration
func CreateGrpcConnection(grpcAddress string) (*grpc.ClientConn, error) {
	var grpcOpts []grpc.DialOption
	if strings.HasPrefix(grpcAddress, "https://") {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	address := HTTPProtocols.ReplaceAllString(grpcAddress, "")
	return grpc.Dial(address, grpcOpts...)
}

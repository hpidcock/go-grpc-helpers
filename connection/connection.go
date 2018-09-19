package connection

import (
	"crypto/tls"
	"errors"
	"net/url"
	"strings"

	"github.com/certifi/gocertifi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	// ErrURLMissingPort is returned when no port is specified.
	ErrURLMissingPort = errors.New("port missing from url")
	// ErrUnknownGRPCScheme is returned when the url scheme is not grpc:// or grpcs://
	ErrUnknownGRPCScheme = errors.New("unkown url scheme provided")
)

// CreateClient creates a client with the grpc:// and grpcs:// scheme.
func CreateClient(uri string, options ...grpc.DialOption) (*grpc.ClientConn, error) {
	address, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	if address.Port() == "" {
		return nil, ErrURLMissingPort
	}

	var conn *grpc.ClientConn
	switch strings.ToLower(address.Scheme) {
	case "grpc":
		options = append([]grpc.DialOption{
			grpc.WithInsecure(),
		}, options...)
	case "grpcs":
		certPool, err := gocertifi.CACerts()
		if err != nil {
			return nil, err
		}
		tlsConfig := &tls.Config{
			ServerName: address.Hostname(),
			RootCAs:    certPool,
		}
		options = append([]grpc.DialOption{
			grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		}, options...)
	default:
		return nil, ErrUnknownGRPCScheme
	}

	conn, err = grpc.Dial(address.Host, options...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

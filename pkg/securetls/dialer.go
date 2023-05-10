package securetls

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"os"
)

var ErrNoCACertsLoadedFromPEM = errors.New("no CA certs were loaded from the PEM file")

func MakeSecureDialer(caFile, clientCertFile, clientCertKey string) func(ctx context.Context, network, addr string) (net.Conn, error) {
	return func(ctx context.Context, network string, addr string) (net.Conn, error) {
		caCert, err := os.ReadFile(caFile)
		if err != nil {
			return nil, err
		}

		certPool := x509.NewCertPool()
		if ok := certPool.AppendCertsFromPEM(caCert); !ok {
			return nil, ErrNoCACertsLoadedFromPEM
		}

		clientCert, err := tls.LoadX509KeyPair(clientCertFile, clientCertKey)
		if err != nil {
			return nil, err
		}

		config := tls.Config{
			MinVersion:   tls.VersionTLS12,
			RootCAs:      certPool,
			Certificates: []tls.Certificate{clientCert},
		}

		dialer := tls.Dialer{
			Config: &config,
		}

		return dialer.DialContext(ctx, network, addr)
	}
}

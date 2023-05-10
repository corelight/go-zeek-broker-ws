package weirdtls

import (
	"context"
	"net"

	"github.com/libp2p/go-openssl"
)

func BrokerDefaultTLSDialer(ctx context.Context, network, addr string) (net.Conn, error) {
	sslCtx, err := openssl.NewCtx()
	if err != nil {
		return nil, err
	}

	err = sslCtx.SetCipherList("AECDH-AES256-SHA@SECLEVEL=0:AECDH-AES256-SHA:P-384")
	if err != nil {
		return nil, err
	}

	return openssl.Dial(network, addr, sslCtx, openssl.InsecureSkipHostVerification)
}

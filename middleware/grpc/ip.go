package server

import (
	"context"
	"net"

	"google.golang.org/grpc/peer"
)

func peerFromContext(ctx context.Context) string {
	if peer, ok := peer.FromContext(ctx); ok {
		switch addr := peer.Addr.(type) {
		case *net.UDPAddr:
			return addr.IP.String()
		case *net.TCPAddr:
			return addr.IP.String()
		default:
			// Hopefully our port removal will work in (almost?) all cases
			return portSuffixRegex.ReplaceAllLiteralString(peer.Addr.String(), "")
		}
	}

	return ""
}

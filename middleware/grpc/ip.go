package server

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"net"
	"strings"
)

func extractGRPCRealIP(ctx context.Context, md metadata.MD) string {
	xForwardedFor := md.Get("x-forwarded-for")
	if len(xForwardedFor) > 0 {
		// When behind a Google Load Balancer, the only two values that we can
		// be sure about are the `n - 2` and `n - 1` (so the last two values
		// in the array). The very last value (`n - 1`) is the Google IP and the
		// `n - 2` value is the actual remote IP that reached the load balancer.
		//
		// When there is more than 2 IPs, all other values prior `n - 2` are
		// those coming from the `X-Forwarded-For` HTTP header received by the load
		// balancer directly, so something a client might have added manually. Since
		// they are coming from an HTTP header and not from Google directly, they
		// can be forged and cannot be trusted.
		//
		// Ideally, to trust the received IP, we should validate it's an actual
		// query coming from Netlify. For now, we are very lenient and trust
		// anything that comes in and looks like an IP.
		//
		// @see https://cloud.google.com/load-balancing/docs/https#x-forwarded-for_header
		if len(xForwardedFor) <= 2 { // 1 or 2
			return strings.TrimSpace(xForwardedFor[0])
		}

		// There is more than 2 addresses, only the element at `n - 2` should be
		// considered, all others cannot be trusted (assuming we got `[a, b, c, d]``,
		// we want to pick element `c` which is at index 2 here so `len(elements) - 2`
		// gives the correct value)
		return strings.TrimSpace(xForwardedFor[len(xForwardedFor)-2]) // more than 2
	}

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

	return "0.0.0.0"
}

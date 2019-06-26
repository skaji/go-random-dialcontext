package random

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"time"
)

// Resolver is
type Resolver interface {
	LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error)
}

// TODO IPv6, IPv4

// DialContext is
func DialContext(resolver Resolver, base func(ctx context.Context, network, addr string) (net.Conn, error)) func(ctx context.Context, network, addr string) (net.Conn, error) {
	if resolver == nil {
		resolver = net.DefaultResolver
	}
	if base == nil {
		base = (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext
	}
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		addrs, err := resolver.LookupIPAddr(ctx, host)
		if err != nil {
			return nil, err
		}
		if len(addrs) == 0 {
			return nil, errors.New("failed to lookup " + host)
		}
		i := rand.Int31n(int32(len(addrs)))
		a := addrs[i]
		return base(ctx, network, net.JoinHostPort(a.IP.String(), port))
	}
}

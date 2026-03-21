package ipv6Info

import (
	"context"
	"errors"
	"fmt"
	"net"
	"regexp"
)

type IpV6Result struct {
	Success       bool
	ServerAddress string
	LocalAddress  string
}

const GoogleIpV6Addr = "2001:4860:4860::8888"

var LocalIpv6AddressError = errors.New("error getting local ipv6 address")

var ipv6AddressRegex = regexp.MustCompile(`\[(?P<address>.+)]`)

func TestIpV6Connection(ctx context.Context, host string) (IpV6Result, error) {

	rtnMe := IpV6Result{
		Success:       false,
		ServerAddress: host,
		LocalAddress:  "",
	}

	d := net.Dialer{}

	localAddr, err := getOutboundIPv6For(ctx, d, host)
	if err != nil {
		return rtnMe, err
	}
	rtnMe.LocalAddress = localAddr

	conn, err := d.DialContext(ctx, "tcp6", net.JoinHostPort(host, "53"))
	if err != nil {
		return rtnMe, err
	}
	defer conn.Close()
	rtnMe.Success = true

	return rtnMe, err
}

// getOutboundIPv6For returns the local IPv6 address the OS would use to reach host.
// Uses UDP so no actual connection or data is sent.
func getOutboundIPv6For(ctx context.Context, d net.Dialer, host string) (string, error) {
	conn, err := d.DialContext(ctx, "udp6", net.JoinHostPort(host, "53"))
	if err != nil {
		return "", fmt.Errorf("could not determine outbound ipv6 address: %w", err)
	}
	defer conn.Close()

	matches := ipv6AddressRegex.FindStringSubmatch(conn.LocalAddr().String())
	if matches == nil {
		return "", fmt.Errorf("failed to parse local address: %w", LocalIpv6AddressError)
	}
	return matches[ipv6AddressRegex.SubexpIndex("address")], nil
}

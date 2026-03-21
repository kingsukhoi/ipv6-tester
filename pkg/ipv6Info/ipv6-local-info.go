package ipv6Info

import (
	"context"
	"fmt"
	"net"
	"slices"
)

type GetIpv6AddressesResult struct {
	Addresses []string
}

func GetLocalIpv6Addresses(ctx context.Context, interfaceName string) (*GetIpv6AddressesResult, error) {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return nil, fmt.Errorf("interface %q not found: %w", interfaceName, err)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	var addresses []string
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		if ip != nil && ip.To4() == nil && ip.To16() != nil {
			addresses = append(addresses, ip.String())
		}
	}

	slices.Sort(addresses)

	return &GetIpv6AddressesResult{Addresses: addresses}, nil
}

// GetInterfaceByIP returns the network interface that has the given IP assigned.
func GetInterfaceByIP(ipString string) (*net.Interface, error) {
	ip := net.ParseIP(ipString)

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		addrs, errL := iface.Addrs()
		if errL != nil {
			continue
		}
		for _, addr := range addrs {
			var ifaceIP net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ifaceIP = v.IP
			case *net.IPAddr:
				ifaceIP = v.IP
			}

			if ifaceIP != nil && ifaceIP.Equal(ip) {
				return &iface, nil
			}
		}
	}
	return nil, &net.AddrError{Err: "no interface found for address", Addr: ip.String()}
}

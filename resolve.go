package resolver

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

type Answer struct {
	Addresses []string
	TTL       uint32
}

// Resolve return a list of Addresses ipv4/ipv6
func (r *Resolver) Resolve(host string) (*Answer, error) {
	// if host is an IP don't, resolve and set TTL to 1 year
	h := strings.Split(host, ":")[0]
	addr := net.ParseIP(h)
	if addr != nil {
		return &Answer{
			Addresses: []string{addr.String()},
			TTL:       31557600,
		}, nil
	}

	ipv4 := &Answer{}
	ipv6 := &Answer{}

	var wg sync.WaitGroup

	wg.Add(2)

	// ipv4
	go func(host string) {
		defer wg.Done()
		ipv4 = r.IPv4(host)
	}(h)

	// ipv6
	go func(host string) {
		defer wg.Done()
		ipv6 = r.IPv6(host)
	}(h)

	wg.Wait()

	if ipv4 == nil && ipv6 == nil {
		return nil, fmt.Errorf("Could not found IP\n")
	}

	ips := &Answer{}

	if ipv4 != nil {
		ips.Addresses = append(ips.Addresses, ipv4.Addresses...)
		ips.TTL = ipv4.TTL
	}

	if ipv6 != nil {
		ips.Addresses = append(ips.Addresses, ipv6.Addresses...)
		if ips.TTL > ipv6.TTL {
			ips.TTL = ipv6.TTL
		}
	}

	return ips, nil
}

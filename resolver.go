package resolver

import "github.com/miekg/dns"

type Resolver struct {
	localServers []string
	server       string
	timeout      int
}

// New return a Resolver intance
func New(server string, t ...int) (*Resolver, error) {
	timeout := 0
	if len(t) > 0 {
		timeout = t[0]
	}
	config, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		return nil, err
	}
	return &Resolver{
		localServers: config.Servers,
		server:       server,
		timeout:      timeout,
	}, nil
}

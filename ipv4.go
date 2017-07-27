package resolver

import (
	"fmt"
	"log"
	"time"

	"github.com/miekg/dns"
)

func (r *Resolver) IPv4(host string) *Answer {
	// start DNS client
	c := dns.Client{
		Timeout: time.Duration(r.timeout) * time.Second,
	}
	m := dns.Msg{}
	m.SetQuestion(fmt.Sprintf("%s.", host), dns.TypeA)
	in, _, err := c.Exchange(&m, fmt.Sprintf("%s:53", r.server))
	var exit bool
	if err != nil {
		exit = true
		log.Printf("server %q not responding, trying with local servers.\n", r.server)
		// if main nameserver not resolving try with local servers
		for i := 0; i < len(r.localServers); i++ {
			in, _, err = c.Exchange(&m, fmt.Sprintf("%s:53", r.localServers[i]))
			if err == nil {
				exit = false
				break
			}
		}
		if exit {
			return nil
		}
	}

	if len(in.Answer) == 0 {
		return nil
	}

	dnsAnswer := Answer{}
	for _, ans := range in.Answer {
		if a, ok := ans.(*dns.A); ok {
			dnsAnswer.Addresses = append(dnsAnswer.Addresses, a.A.String())
			dnsAnswer.TTL = dnsAnswer.TTL + ans.Header().Ttl
		}
	}

	// get the average TTL
	dnsAnswer.TTL = dnsAnswer.TTL / uint32(len(dnsAnswer.Addresses))

	if len(dnsAnswer.Addresses) == 0 {
		return nil
	}
	return &dnsAnswer
}

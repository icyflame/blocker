package blocker

import (
	"context"
	"net"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

const PluginName = "blocker"

func init() {
	plugin.Register(PluginName, setup)
}

func setup(c *caddy.Controller) error {
	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Blocker{
			Next:    next,
			Decider: DefaultBlockDomainsDecider{},
		}
	})

	// All OK, return a nil error.
	return nil
}

type Blocker struct {
	Next    plugin.Handler
	Decider BlockDomainsDecider
}

type BlockDomainsDecider interface {
	IsDomainBlocked(domain string, questionType uint16) bool
}

func (b Blocker) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	if len(r.Question) == 0 {
		// If the query has no question, then we can't do anything with it here. So, I will just
		// pass it to the next plugin.
		return b.Next.ServeDNS(ctx, w, r)
	}

	question := r.Question[0]

	domain := question.Name
	questionType := question.Qtype

	// If the domain is blocked, then reply right now with the answer 0.0.0.0
	if b.Decider.IsDomainBlocked(domain, questionType) {
		// TODO: Count up a metric here for blocked domains

		// TODO: This answer should be calculated based on the type of query
		emptyAnswer := &dns.A{
			Hdr: dns.RR_Header{
				Name:   domain,
				Class:  dns.ClassINET,
				Rrtype: questionType,
			},
			A: net.IPv4zero,
		}

		response := &dns.Msg{
			Answer: []dns.RR{
				emptyAnswer,
			},
		}
		response.SetReply(r)
		w.WriteMsg(response)
		return dns.RcodeSuccess, nil
	}

	return b.Next.ServeDNS(ctx, w, r)
}

// Name ...
func (b Blocker) Name() string {
	return PluginName
}

type DefaultBlockDomainsDecider struct {
	blockedDomains []string
}

// IsBlocked ...
func (d DefaultBlockDomainsDecider) IsDomainBlocked(domain string, questionType uint16) bool {
	if questionType != dns.TypeA {
		return false
	}

	return domain == "fb.com."
}

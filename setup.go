package blocker

import (
	"context"
	"fmt"

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
		return Blocker{Next: next}
	})

	// All OK, return a nil error.
	return nil
}

type Blocker struct {
	Next plugin.Handler
}

func (b Blocker) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	fmt.Print("serving DNS")
	return b.Next.ServeDNS(ctx, w, r)
}

// Name ...
func (b Blocker) Name() string {
	return PluginName
}

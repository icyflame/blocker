package blocker

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
)

const PluginName = "blocker"
const RequiredArgs = 4

func init() {
	plugin.Register(PluginName, setup)
}

func setup(c *caddy.Controller) error {
	var args []string
	c.NextArg() // Skip the name of the plugin, which is returned as an argument
	for c.NextArg() {
		args = append(args, c.Val())
	}

	if len(args) < RequiredArgs {
		// Any errors returned from this setup function should be wrapped with plugin.Error, so we
		// can present a slightly nicer error message to the user.
		return plugin.Error(PluginName, c.ArgErr())
	}

	blocklistFilePath := args[0]
	blocklistUpdateFrequency := args[1]
	blocklistType := args[2]
	if blocklistType != string(BlocklistType_Hosts) && blocklistType != string(BlocklistType_ABP) {
		return plugin.Error(PluginName, c.ArgErr())
	}

	blocklistResponseType := args[3]
	if blocklistResponseType != string(ResponseType_Empty) && blocklistResponseType != string(ResponseType_NXDOMAIN) {
		return plugin.Error(PluginName, c.ArgErr())
	}

	logger := clog.NewWithPlugin(PluginName)

	decider, shutdownHooks, err := PrepareBlocklist(blocklistFilePath, blocklistUpdateFrequency, blocklistType, logger)
	if err != nil {
		return plugin.Error(PluginName, err)
	}

	for _, hook := range shutdownHooks {
		c.OnShutdown(hook)
	}

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Blocker{
			Next:         next,
			Decider:      decider,
			ResponseType: blocklistResponseType,
		}
	})

	// All OK, return a nil error.
	return nil
}

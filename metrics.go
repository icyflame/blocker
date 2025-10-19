package blocker

import (
	"github.com/coredns/coredns/plugin"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// RequestBlockCount is the number of DNS requests being blocked.
	RequestBlockCount = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: PluginName,
		Name:      "blocked_requests_total",
		Help:      "Counter of DNS requests being blocked.",
	})
	// RequestAllowCount is the number of DNS requests being Allowed.
	RequestAllowCount = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: PluginName,
		Name:      "allowed_requests_total",
		Help:      "Counter of DNS requests being allowed.",
	})
)

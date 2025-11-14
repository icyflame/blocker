package blocker

import (
	"context"
	"net"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metadata"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Blocker struct {
	Next         plugin.Handler
	Decider      BlockDomainsDecider
	ResponseType string
}

type ResponseType string

const (
	ResponseType_Empty    ResponseType = "empty"
	ResponseType_NXDOMAIN ResponseType = "nxdomain"
)

const MetadataRequestBlocked = "blocker/request-blocked"

// Metadata implements the Metadata plugin's required interface in the blocker function.
func (b Blocker) Metadata(ctx context.Context, state request.Request) context.Context {
	metadata.SetValueFunc(ctx, MetadataRequestBlocked, func() string {
		return "NO"
	})
	return ctx
}

func (b Blocker) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	if len(r.Question) == 0 {
		// If the query has no question, then we can't do anything with it here. So, I will just
		// pass it to the next plugin.
		return b.Next.ServeDNS(ctx, w, r)
	}

	question := r.Question[0]

	// TODO: After some time, check if clients are making non-A/AAAA requests to the DNS server. I
	// don't think of any reason a browser should request anything except the IP address for a
	// domain. But I don't want to block everything right away.
	//
	// One thing that I should keep in mind is that some domains have CNAMEs in order to cloak their
	// real purpose. So, if this is CNAME, and the domain matches the known list of CNAME trackers,
	// then the request should be blocked as well.
	//
	// Do the same thing for question class (INET)
	questionType := question.Qtype
	if questionType != dns.TypeA && questionType != dns.TypeAAAA {
		return b.Next.ServeDNS(ctx, w, r)
	}

	questionClass := question.Qclass
	if questionClass != dns.ClassINET {
		return b.Next.ServeDNS(ctx, w, r)
	}

	// If the domain is blocked, then reply right now with the answer 0.0.0.0 for IPv4 and :: for
	// IPv6
	domain := question.Name
	if b.Decider.IsDomainBlocked(domain) {
		RequestBlockCount.Inc()
		switch ResponseType(b.ResponseType) {
		case ResponseType_Empty:
			response := &dns.Msg{
				Answer: []dns.RR{
					GetEmptyAnswerForQuestionType(questionType, domain),
				},
			}
			response.SetReply(r)
			w.WriteMsg(response)
			metadata.SetValueFunc(ctx, MetadataRequestBlocked, func() string {
				return "YES"
			})
			return dns.RcodeSuccess, nil

		case ResponseType_NXDOMAIN:
			response := &dns.Msg{}
			response.SetRcode(r, dns.RcodeNameError)
			w.WriteMsg(response)
			metadata.SetValueFunc(ctx, MetadataRequestBlocked, func() string {
				return "YES"
			})
			return dns.RcodeNameError, nil
		}
	}
	RequestAllowCount.Inc()

	return b.Next.ServeDNS(ctx, w, r)
}

// Name ...
func (b Blocker) Name() string {
	return PluginName
}

// GetEmptyAnswerForQuestionType ...
func GetEmptyAnswerForQuestionType(questionType uint16, domain string) dns.RR {
	responseHeader := dns.RR_Header{
		Name:   domain,
		Class:  dns.ClassINET,
		Rrtype: questionType,
	}
	if questionType == dns.TypeAAAA {
		return &dns.AAAA{
			Hdr:  responseHeader,
			AAAA: net.IPv6zero,
		}
	}

	return &dns.A{
		Hdr: responseHeader,
		A:   net.IPv4zero,
	}
}

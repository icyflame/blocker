package blocker

type BlockDomainsDecider interface {
	IsDomainBlocked(domain string) bool
}

type BlocklistType string

const BlocklistType_Hosts BlocklistType = "hosts"

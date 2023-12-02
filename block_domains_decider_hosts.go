package blocker

import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type BlockDomainsDeciderHosts struct {
	blocklist     map[string]bool
	blocklistFile string
	lastUpdated   time.Time
	log           Logger
}

// Name ...
func NewBlockDomainsDeciderHosts(filePath string, logger Logger) BlockDomainsDecider {
	return &BlockDomainsDeciderHosts{
		blocklistFile: filePath,
		log:           logger,
		blocklist:     map[string]bool{},
	}
}

// IsDomainBlocked ...
func (d *BlockDomainsDeciderHosts) IsDomainBlocked(domain string) bool {
	return d.blocklist[domain]
}

// StartBlocklistUpdater ...
func (d *BlockDomainsDeciderHosts) StartBlocklistUpdater(ticker *time.Ticker) {
	go func() {
		for true {
			tick := <-ticker.C
			d.log.Debugf("Ticker arrived at time: %v", tick)

			if d.IsBlocklistUpdateRequired() {
				d.log.Debug("update required")
				d.UpdateBlocklist()
			} else {
				d.log.Debug("update not required")
			}
		}
	}()
}

// UpdateBlocklist ...
func (d *BlockDomainsDeciderHosts) UpdateBlocklist() error {
	// Update process
	blocklistContent, err := os.Open(d.blocklistFile)
	if err != nil {
		d.log.Errorf("could not read blocklist file: %s", d.blocklistFile)
		return err
	}
	defer blocklistContent.Close()

	numBlockedDomainsBefore := len(d.blocklist)
	lastUpdatedBefore := d.lastUpdated

	scanner := bufio.NewScanner(blocklistContent)
	for scanner.Scan() {
		hostLine := scanner.Text()
		comps := strings.Split(hostLine, " ")
		if len(comps) < 2 {
			// Bad line in the input file
			d.log.Warningf("unformatted line present in the input file: %s", hostLine)
			continue
		}

		domain := comps[1]
		d.blocklist[dns.Fqdn(domain)] = true
	}

	d.lastUpdated = time.Now()

	d.log.Infof("updated blocklist; blocked domains: before: %d, after: %d; last updated: before: %v, after: %v",
		numBlockedDomainsBefore, len(d.blocklist), lastUpdatedBefore, d.lastUpdated)

	return nil
}

// IsBlocklistUpdateRequired ...
func (d *BlockDomainsDeciderHosts) IsBlocklistUpdateRequired() bool {
	blocklistFileStat, _ := os.Stat(d.blocklistFile)
	return blocklistFileStat.ModTime().After(d.lastUpdated)
}

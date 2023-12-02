package blocker

import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type BlockDomainsDeciderHosts struct {
	Blocklist     map[string]bool
	BlocklistFile string
	LastUpdated   time.Time
	Log           Logger
}

// IsDomainBlocked ...
func (d *BlockDomainsDeciderHosts) IsDomainBlocked(domain string) bool {
	return d.Blocklist[domain]
}

// StartBlocklistUpdater ...
func (d *BlockDomainsDeciderHosts) StartBlocklistUpdater(ticker *time.Ticker) {
	go func() {
		for true {
			tick := <-ticker.C
			d.Log.Debugf("Ticker arrived at time: %v", tick)

			if d.IsBlocklistUpdateRequired() {
				d.Log.Debug("update required")
				d.UpdateBlocklist()
			} else {
				d.Log.Debug("update not required")
			}
		}
	}()
}

// UpdateBlocklist ...
func (d *BlockDomainsDeciderHosts) UpdateBlocklist() error {
	// Update process
	blocklistContent, err := os.Open(d.BlocklistFile)
	if err != nil {
		d.Log.Errorf("could not read blocklist file: %s", d.BlocklistFile)
		return err
	}
	defer blocklistContent.Close()

	numBlockedDomainsBefore := len(d.Blocklist)
	lastUpdatedBefore := d.LastUpdated

	scanner := bufio.NewScanner(blocklistContent)
	for scanner.Scan() {
		hostLine := scanner.Text()
		comps := strings.Split(hostLine, " ")
		if len(comps) < 2 {
			// Bad line in the input file
			d.Log.Warningf("unformatted line present in the input file: %s", hostLine)
			continue
		}

		domain := comps[1]
		d.Blocklist[dns.Fqdn(domain)] = true
	}

	d.LastUpdated = time.Now()

	d.Log.Infof("updated blocklist; blocked domains: before: %d, after: %d; last updated: before: %v, after: %v",
		numBlockedDomainsBefore, len(d.Blocklist), lastUpdatedBefore, d.LastUpdated)

	return nil
}

// IsBlocklistUpdateRequired ...
func (d *BlockDomainsDeciderHosts) IsBlocklistUpdateRequired() bool {
	blocklistFileStat, _ := os.Stat(d.BlocklistFile)
	return blocklistFileStat.ModTime().After(d.LastUpdated)
}

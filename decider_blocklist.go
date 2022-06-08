package blocker

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type BlocklistBasedDecider struct {
	Blocklist     map[string]bool
	BlocklistFile string
	LastUpdated   time.Time
}

// PrepareBlocklist ...
func PrepareBlocklist(filePath string) (BlockDomainsDecider, []func() error, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, nil, err
	}

	decider := &BlocklistBasedDecider{
		Blocklist:     map[string]bool{},
		BlocklistFile: filePath,
	}

	ticker := time.NewTicker(time.Second)
	decider.StartBlocklistUpdater(ticker)

	stopTicker := func() error {
		fmt.Println("Ticker was stopped.")
		ticker.Stop()
		return nil
	}

	shutdownHooks := []func() error{
		stopTicker,
	}

	return decider, shutdownHooks, nil
}

// IsDomainBlocked ...
func (d *BlocklistBasedDecider) IsDomainBlocked(domain string) bool {
	return d.Blocklist[domain]
}

// StartBlocklistUpdater ...
func (d *BlocklistBasedDecider) StartBlocklistUpdater(ticker *time.Ticker) {
	go func() {
		for true {
			tick := <-ticker.C
			fmt.Println("Ticker arrived at time: ", tick)

			if d.IsBlocklistUpdateRequired() {
				fmt.Println("update required")
				d.UpdateBlocklist()
			} else {
				fmt.Println("update not required")
			}
		}
	}()
}

// UpdateBlocklist ...
func (d *BlocklistBasedDecider) UpdateBlocklist() error {
	// Update process
	blocklistContent, err := os.Open(d.BlocklistFile)
	if err != nil {
		fmt.Println("[ERROR] could not read blocklist file", d.BlocklistFile)
		return err
	}
	defer blocklistContent.Close()

	scanner := bufio.NewScanner(blocklistContent)
	for scanner.Scan() {
		hostLine := scanner.Text()
		comps := strings.Split(hostLine, " ")
		if len(comps) < 2 {
			// Bad line in the input file
			fmt.Println("[WARN] unformatted line present in the input file: ", hostLine)
			continue
		}

		domain := comps[1]
		d.Blocklist[dns.Fqdn(domain)] = true
	}

	d.LastUpdated = time.Now()
	return nil
}

// IsBlocklistUpdateRequired ...
func (d *BlocklistBasedDecider) IsBlocklistUpdateRequired() bool {
	blocklistFileStat, _ := os.Stat(d.BlocklistFile)
	return blocklistFileStat.ModTime().After(d.LastUpdated)
}

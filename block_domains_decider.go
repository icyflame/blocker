package blocker

import (
	"fmt"
	"os"
	"time"
)

// BlockDomainsDecider is the interface which must be implemented by any type which intends to
// become a blocker. The purpose of each of the functions is described below.
type BlockDomainsDecider interface {
	IsDomainBlocked(domain string) bool
	StartBlocklistUpdater(ticker *time.Ticker)
	UpdateBlocklist() error
}

type BlocklistType string

const BlocklistType_Hosts BlocklistType = "hosts"
const BlocklistType_ABP BlocklistType = "abp"

// PrepareBlocklist ...
func PrepareBlocklist(filePath string, blocklistUpdateFrequency string, blocklistType string, logger Logger) (BlockDomainsDecider, []func() error, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, nil, err
	}

	frequency, err := time.ParseDuration(blocklistUpdateFrequency)
	if err != nil {
		return nil, nil, err
	}

	var decider BlockDomainsDecider
	switch BlocklistType(blocklistType) {
	case BlocklistType_Hosts:
		decider = NewBlockDomainsDeciderHosts(filePath, logger)
	case BlocklistType_ABP:
		decider = NewBlockDomainsDeciderABP(filePath, logger)
	}

	// Always update the blocklist when the server starts up
	decider.UpdateBlocklist()

	// Setup periodic updation of the blocklist
	ticker := time.NewTicker(frequency)
	decider.StartBlocklistUpdater(ticker)

	stopTicker := func() error {
		fmt.Println("[INFO] Ticker was stopped.")
		ticker.Stop()
		return nil
	}

	shutdownHooks := []func() error{
		stopTicker,
	}

	return decider, shutdownHooks, nil
}

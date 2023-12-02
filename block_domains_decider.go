package blocker

import (
	"fmt"
	"os"
	"time"
)

type BlockDomainsDecider interface {
	IsDomainBlocked(domain string) bool
}

type BlocklistType string

const BlocklistType_Hosts BlocklistType = "hosts"

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

	decider := &BlockDomainsDeciderHosts{
		Blocklist:     map[string]bool{},
		BlocklistFile: filePath,
		Log:           logger,
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

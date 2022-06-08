package blocker

import (
	"os"
	"time"
)

type BlocklistBasedDecider struct {
	Blocklist   map[string]bool
	LastUpdated time.Time
}

// PrepareBlocklist ...
func PrepareBlocklist(filePath string) (BlockDomainsDecider, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	decider := &BlocklistBasedDecider{
		Blocklist:   map[string]bool{},
		LastUpdated: fileInfo.ModTime(),
	}

	decider.Blocklist["fb.com."] = true
	decider.Blocklist["google.com."] = true

	return decider, nil
}

// IsDomainBlocked ...
func (d *BlocklistBasedDecider) IsDomainBlocked(domain string, questionType uint16) bool {
	return d.Blocklist[domain]
}

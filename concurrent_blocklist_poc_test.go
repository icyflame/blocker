package blocker

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

type nopLogger struct{}

func (nopLogger) Error(v ...any)                   {}
func (nopLogger) Errorf(format string, v ...any)   {}
func (nopLogger) Warning(v ...any)                 {}
func (nopLogger) Warningf(format string, v ...any) {}
func (nopLogger) Info(v ...any)                    {}
func (nopLogger) Infof(format string, v ...any)    {}
func (nopLogger) Debug(v ...any)                   {}
func (nopLogger) Debugf(format string, v ...any)   {}
func (nopLogger) Fatal(v ...any)                   {}
func (nopLogger) Fatalf(format string, v ...any)   {}

// PoC for concurrent map read/write bug in BlockDomainsDeciderABP
// At the same time a test for thread-safe access to the blocklist
func TestConcurrentBlocklistAccess_PoC(t *testing.T) {
	// Create a temporary blocklist file
	f, err := os.CreateTemp("", "blocklist-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())

	// Write initial blocklist entries
	for i := 0; i < 10000; i++ {
		fmt.Fprintf(f, "||bad%d.com.^\n", i)
	}
	f.Sync()

	decider := NewBlockDomainsDeciderABP(f.Name(), nopLogger{}).(*BlockDomainsDeciderABP)

	// Writer goroutine: simulates UpdateBlocklist
	go func() {
		for i := 0; i < 10; i++ {
			decider.UpdateBlocklist()
		}
	}()

	// Reader goroutines: simulates IsDomainBlocked
	var wg sync.WaitGroup
	for r := range 10 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < 10000; i++ {
				decider.IsDomainBlocked(fmt.Sprintf("bad%d.com.", i))
			}
		}(r)
	}
	wg.Wait()
}

package blocker

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestIsDomainBlocked(t *testing.T) {
	type testCase struct {
		blocklist    map[string]bool
		queryDomains []string
		wantResults  []bool
	}

	testCases := map[string]testCase{
		"base case": {
			blocklist: map[string]bool{
				"example.com": true,
			},
			queryDomains: []string{
				"example.com",
				"example.com.org",
				"subdomain.example.com",
				"example.org",
			},
			wantResults: []bool{
				true,
				false,
				false,
				false,
			},
		},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := []bool{}
			decider := BlockDomainsDeciderHosts{
				blocklist: tc.blocklist,
			}
			for _, domain := range tc.queryDomains {
				_ = domain
				got = append(got, decider.IsDomainBlocked(domain))
			}

			if diff := pretty.Compare(got, tc.wantResults); diff != "" {
				t.Errorf("unexpected (-got +want): \n%s", diff)
			}
		})
	}
}

package blocker

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestIsDomainBlocked_ABP(t *testing.T) {
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
				"subdomain.example.com",
				"sub1.sub2.example.com",
				"example.com.org",
				"example.org",
			},
			wantResults: []bool{
				true,
				true,
				true,
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
			decider := BlockDomainsDeciderABP{
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

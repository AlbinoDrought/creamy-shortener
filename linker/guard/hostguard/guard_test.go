package hostguard

import (
	"net/url"
	"testing"

	"github.com/AlbinoDrought/creamy-shortener/linker"
)

func TestMakeAllowed(t *testing.T) {
	type allowed struct {
		url     string
		allowed bool
	}
	tests := []struct {
		name   string
		hosts  []string
		checks []allowed
		want   linker.LinkGuard
	}{
		{
			name: "single host",
			hosts: []string{
				"example.com",
			},
			checks: []allowed{
				{"https://example.com/foo?bar", true},
				{"http://example.com/", true},
				{"https://example.com:6969/foo?bar", false},
			},
		},
		{
			name: "localhost with two non-standard ports",
			hosts: []string{
				"localhost:3000",
				"localhost:8080",
			},
			checks: []allowed{
				{"http://localhost:3000/some/url", true},
				{"http://localhost:8080/some/url", true},
				{"http://localhost/some/url", false},
				{"http://localhost:42069/some/url", false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			guard := Make(tt.hosts)
			for _, check := range tt.checks {
				link, err := url.Parse(check.url)
				if err != nil {
					t.Errorf("failed parsing %v: %v", link, err)
				}

				allowed := guard.Allowed(link)
				if allowed != check.allowed {
					t.Errorf("link %v should have been %v but was %v", check.url, check.allowed, allowed)
				}
			}
		})
	}
}

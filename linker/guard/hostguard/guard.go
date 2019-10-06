package hostguard

import (
	"net/url"

	"github.com/AlbinoDrought/creamy-shortener/linker"
)

type guard struct {
	hosts map[string]bool
}

func (g *guard) Allowed(url *url.URL) bool {
	_, ok := g.hosts[url.Host]
	return ok
}

// Make a link guard that allows links based on a host whitelist
func Make(allowedHosts []string) linker.LinkGuard {
	hostMap := make(map[string]bool)
	for _, host := range allowedHosts {
		hostMap[host] = true
	}
	return &guard{hostMap}
}

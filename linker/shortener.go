package linker

import (
	"errors"
	"net/url"
)

// A LinkGuard determines if a URL is allowed to be shortened
type LinkGuard interface {
	Allowed(url *url.URL) bool
}

// A LinkRepository persists ID <=> URL pairs
type LinkRepository interface {
	Get(id string) (*url.URL, error)
	Set(id string, url *url.URL) error
}

// LinkMapper generates IDs from URLs
type LinkMapper interface {
	Map(link *url.URL) (string, error)
}

// A LinkShortener turns fully-qualified URLS into shorter IDs, and back again
type LinkShortener struct {
	guard  LinkGuard
	repo   LinkRepository
	mapper LinkMapper
}

// Allowed returns true if the link is allowed to be shortened
func (repo *LinkShortener) Allowed(link *url.URL) bool {
	return repo.guard.Allowed(link)
}

// Shorten a long URL to a shorter ID
func (repo *LinkShortener) Shorten(link *url.URL) (string, error) {
	// make sure we are allowed to shorten the link:
	if !repo.guard.Allowed(link) {
		return "", errors.New("link not allowed")
	}

	// turn the link into a shorter ID:
	id, err := repo.mapper.Map(link)
	if err != nil {
		return "", err
	}

	// save it for later:
	err = repo.repo.Set(id, link)
	if err != nil {
		return "", err
	}

	return id, nil
}

// Expand a shorter ID to a long URL
func (repo *LinkShortener) Expand(id string) (*url.URL, error) {
	// fetch the URL we saved for this ID:
	url, err := repo.repo.Get(id)
	if err != nil {
		return nil, err
	}

	// make sure it's still allowed:
	if !repo.guard.Allowed(url) {
		return nil, errors.New("link not allowed")
	}

	return url, nil
}

// Make a default shortener with passed implementations
func Make(guard LinkGuard, mapper LinkMapper, repo LinkRepository) *LinkShortener {
	return &LinkShortener{
		guard,
		repo,
		mapper,
	}
}

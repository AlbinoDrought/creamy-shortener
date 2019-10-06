package memoryrepo

import (
	"errors"
	"net/url"

	"github.com/AlbinoDrought/creamy-shortener/linker"
)

type repo struct {
	urls map[string]*url.URL
}

func (r *repo) Get(id string) (*url.URL, error) {
	url, ok := r.urls[id]
	if !ok {
		return nil, errors.New("url not found")
	}
	return url, nil
}

func (r *repo) Set(id string, url *url.URL) error {
	r.urls[id] = url
	return nil
}

// Make an in-memory temporary link repository
func Make() linker.LinkRepository {
	return &repo{}
}

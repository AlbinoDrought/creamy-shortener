package aferorepo

import (
	"net/url"
	"os"

	"github.com/AlbinoDrought/creamy-shortener/linker"

	"github.com/spf13/afero"
)

type repo struct {
	fs afero.Fs
}

func (r *repo) Get(id string) (*url.URL, error) {
	rawLink, err := afero.ReadFile(r.fs, id)
	if err != nil {
		return nil, err
	}

	return url.Parse(string(rawLink))
}

func (r *repo) Set(id string, url *url.URL) error {
	return afero.WriteFile(r.fs, id, []byte(url.String()), os.ModePerm)
}

// Make a link repository that persists to an afero filesystem
func Make(fs afero.Fs) (linker.LinkRepository, error) {
	if err := fs.MkdirAll(".", os.ModePerm); err != nil {
		return nil, err
	}
	return &repo{fs}, nil
}

// Must make a link repository with the given filesystem, or panic on failure
func Must(fs afero.Fs) linker.LinkRepository {
	linkRepository, err := Make(fs)
	if err != nil {
		panic(err)
	}
	return linkRepository
}

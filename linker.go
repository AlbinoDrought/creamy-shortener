package main

import (
	"errors"
	"net/url"
	"os"
	"path"

	"github.com/spf13/afero"

	"github.com/multiformats/go-multihash"
)

type linker struct {
	fs           afero.Fs
	appURL       string
	hashMode     string
	allowedHosts []string
}

func makeLinker(fs afero.Fs, appURL, hashMode string) linker {
	return linker{
		fs,
		appURL,
		hashMode,
		[]string{},
	}
}

func (repo *linker) dataPart(piece string) (string, error) {
	hash, err := multihash.EncodeName([]byte(piece), repo.hashMode)
	if err != nil {
		return "", err
	}

	mh, err := multihash.Sum(hash, multihash.Names[repo.hashMode], -1)
	if err != nil {
		return "", err
	}

	return mh.B58String(), nil
}

func (repo *linker) Allowed(url *url.URL) error {
	host := url.Host

	for _, allowedHost := range repo.allowedHosts {
		if host == allowedHost {
			return nil
		}
	}

	return errors.New("host not allowed")
}

func (repo *linker) Allow(host string) error {
	repo.allowedHosts = append(repo.allowedHosts, host)
	return nil
}

func (repo *linker) Shorten(link *url.URL) (string, error) {
	if err := repo.Allowed(link); err != nil {
		return "", err
	}

	linkString := link.String()
	linkDataPart, err := repo.dataPart(linkString)
	if err != nil {
		return "", err
	}

	hostPath := path.Join("link", linkDataPart)

	repo.fs.MkdirAll(path.Dir(hostPath), os.ModePerm)

	err = afero.WriteFile(repo.fs, hostPath, []byte(linkString), os.ModePerm)
	if err != nil {
		return "", err
	}

	fullURL, err := url.Parse(repo.appURL)
	if err != nil {
		return "", err
	}
	fullURL.Path = "/l/" + linkDataPart

	return fullURL.String(), nil
}

func (repo *linker) Expand(shortenedPart string) (*url.URL, error) {
	hostPath := path.Join("link", shortenedPart)

	contents, err := afero.ReadFile(repo.fs, hostPath)
	if err != nil {
		return nil, err
	}

	return url.Parse(string(contents))
}

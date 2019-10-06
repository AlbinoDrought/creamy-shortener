package main

import (
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	"github.com/multiformats/go-multihash"
)

type linker struct {
	directory    string
	appURL       string
	hashMode     string
	allowedHosts []string
}

func makeLinker(directory, appURL, hashMode string) linker {
	return linker{
		directory,
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

	return hex.EncodeToString(hash), nil
}

func (repo *linker) localPath(pieces ...string) string {
	piecePath := path.Join(pieces...)
	return path.Join(repo.directory, piecePath)
}

func (repo *linker) dataPath(category, value string) (string, error) {
	dataPart, err := repo.dataPart(value)
	if err != nil {
		return "", err
	}

	return repo.localPath(category, dataPart), nil
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
	hostPath, err := repo.dataPath("link", linkString)

	if err != nil {
		return "", err
	}

	os.MkdirAll(path.Dir(hostPath), os.ModePerm)

	err = ioutil.WriteFile(hostPath, []byte(linkString), os.ModePerm)
	if err != nil {
		return "", err
	}

	shortenedPart, err := repo.dataPart(linkString)
	if err != nil {
		return "", err
	}

	fullURL, err := url.Parse(repo.appURL)
	if err != nil {
		return "", err
	}
	fullURL.Path = "/l/" + shortenedPart

	return fullURL.String(), nil
}

func (repo *linker) Expand(shortenedPart string) (*url.URL, error) {
	hostPath := repo.localPath("link", shortenedPart)

	contents, err := ioutil.ReadFile(hostPath)
	if err != nil {
		return nil, err
	}

	return url.Parse(string(contents))
}

package main

import (
	"encoding/hex"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	"github.com/multiformats/go-multihash"
)

type linker struct {
	directory string
	appURL    string
}

func (repo *linker) dataPart(piece string) (string, error) {
	hash, err := multihash.EncodeName([]byte(piece), "sha2-256")

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
	hostPath, err := repo.dataPath("host", url.Host)

	if err != nil {
		return err
	}

	_, err = os.Stat(hostPath)

	return err
}

func (repo *linker) Allow(host string) error {
	hostPath, err := repo.dataPath("host", host)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Dir(hostPath), os.ModePerm)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(hostPath, []byte(host), os.ModePerm)
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

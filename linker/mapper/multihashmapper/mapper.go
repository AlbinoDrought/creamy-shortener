package multihashmapper

import (
	"errors"
	"net/url"

	"github.com/AlbinoDrought/creamy-shortener/linker"
	"github.com/multiformats/go-multihash"
)

type mapper struct {
	multihashCode uint64
}

func (m *mapper) Map(link *url.URL) (string, error) {
	hash, err := multihash.Sum([]byte(link.String()), m.multihashCode, -1)
	if err != nil {
		return "", err
	}

	return hash.B58String(), nil
}

// Make a multihash LinkMapper using the given multihash mode
func Make(multihashName string) (linker.LinkMapper, error) {
	code, ok := multihash.Names[multihashName]
	if !ok {
		return nil, errors.New("bad multihash name")
	}

	return &mapper{code}, nil
}

// Must make a multihash LinkMapper with the given mode, or panic on failure
func Must(multihashName string) linker.LinkMapper {
	mapper, err := Make(multihashName)
	if err != nil {
		panic(err)
	}
	return mapper
}

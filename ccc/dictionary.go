package ccc

import (
	"errors"
	"github.com/crawlcoin/ccc/providers"
)

const SHARED = "shared"

var INVALID_HOSTS = map[string]bool{"shared": true}

type Version struct {
	Version int
}

func Combined(p providers.DictionaryProvider, host string, version int, sharedVersion int) ([]byte, error) {
	hostDict, err := HostDictionary(p, host, version)
	if err != nil {
		return nil, err
	}
	sharedDict, err := SharedDictionary(p, sharedVersion)
	if err != nil {
		return nil, err
	}
	return Append(hostDict, sharedDict), nil
}

func Append(dict1 []byte, dict2 []byte) []byte {
	return append(dict1, dict2...)
}

func SharedDictionary(p providers.DictionaryProvider, version int) ([]byte, error) {
	if version == 0 {
		return []byte{}, nil
	}
	return p.SharedDictionary(version)
}

func HostDictionary(p providers.DictionaryProvider, host string, version int) ([]byte, error) {
	if _, prs := INVALID_HOSTS[host]; prs {
		return nil, errors.New("Invalid host, some names are reserved.")
	}
	if version == 0 {
		return []byte{}, nil
	}
	return p.CustomDictionary(host, version)
}

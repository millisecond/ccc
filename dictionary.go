package ccc

import (
	"errors"
)

const SHARED = "shared"

var INVALID_HOSTS = map[string]bool{"shared": true}

type Version struct {
	Version int
}

func Combined(p DictionaryProvider, id string, customVersion int, sharedVersion int) ([]byte, error) {
	customDict, err := CustomDictionary(p, id, customVersion)
	if err != nil {
		return nil, err
	}
	sharedDict, err := SharedDictionary(p, sharedVersion)
	if err != nil {
		return nil, err
	}
	return Append(customDict, sharedDict), nil
}

func Append(dict1 []byte, dict2 []byte) []byte {
	return append(dict1, dict2...)
}

func SharedDictionary(p DictionaryProvider, version int) ([]byte, error) {
	if version == 0 {
		return []byte{}, nil
	}
	return p.SharedDictionary(version)
}

func CustomDictionary(p DictionaryProvider, id string, version int) ([]byte, error) {
	if _, prs := INVALID_HOSTS[id]; prs {
		return nil, errors.New("Invalid id, some names are reserved.")
	}
	if version == 0 {
		return []byte{}, nil
	}
	return p.CustomDictionary(id, version)
}

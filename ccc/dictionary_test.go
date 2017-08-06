package ccc

import (
	"github.com/facebookgo/ensure"
	"testing"
	"github.com/crawlcoin/ccc/providers"
)

func TestAppend(t *testing.T) {
	ensure.DeepEqual(t, Append([]byte{1, 2}, []byte{3, 4}), []byte{1, 2, 3, 4})
}

func TestVersions(t *testing.T) {
	provider := providers.MemoryDictionaryProvider{}
	provider.SharedDictionaries = map[int][]byte{1: []byte("SHARED")}
	provider.CustomDictionaries = map[string]map[int][]byte{"": {1: []byte("CUSTOM")}}

	combined, err := Combined(provider, "", 1, 1)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, string(combined), "CUSTOMSHARED")
}

func TestZeroDictVersion(t *testing.T) {
	provider := providers.FileDictionaryProvider{}

	shared, err := SharedDictionary(provider, 0)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, shared, []byte{})

	custom, err := CustomDictionary(provider, "crawlcoin.com", 0)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, custom, []byte{})
}

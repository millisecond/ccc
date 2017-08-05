package ccc

import (
	"testing"
	"github.com/facebookgo/ensure"
	"github.com/crawlcoin/ccc/providers"
)

func TestRoundtrip(t *testing.T) {
	raw := []byte{1, 2, 3}

	provider := providers.FileDictionaryProvider{}

	compressed, err := Compress(provider, raw, "", 0, 0)
	ensure.Nil(t, err)

	ensure.NotDeepEqual(t, raw, compressed)

	decompressed, err := Decompress(provider, compressed, "", 0, 0)
	ensure.Nil(t, err)

	ensure.DeepEqual(t, raw, decompressed)
}

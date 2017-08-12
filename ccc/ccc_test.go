package ccc

import (
	"github.com/crawlcoin/ccc/providers"
	"github.com/facebookgo/ensure"
	"testing"
)

func TestRoundtripZero(t *testing.T) {
	raw := []byte{1, 2, 3}
	provider := providers.FileDictionaryProvider{}

	compressed, err := Compress(provider, raw, "", 0, 0)
	ensure.Nil(t, err)
	ensure.NotDeepEqual(t, raw, compressed)

	decompressed, err := Decompress(provider, compressed, "", 0, 0)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, raw, decompressed)
}

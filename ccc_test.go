package ccc

import (
	"compress/zlib"
	"github.com/crawlcoin/ccc/providers"
	"github.com/facebookgo/ensure"
	"testing"
)

func TestRoundtripZero(t *testing.T) {
	raw := []byte{1, 2, 3}
	provider := providers.FileDictionaryProvider{}

	compressed, err := BrotliCompress(provider, raw, "", 0, 0)
	ensure.Nil(t, err)
	ensure.NotDeepEqual(t, raw, compressed)

	decompressed, err := BrotliDecompress(provider, compressed, "", 0, 0)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, raw, decompressed)

	compressed, err = ZlibCompress(provider, zlib.BestCompression, raw, "", 0, 0)
	ensure.Nil(t, err)
	ensure.NotDeepEqual(t, raw, compressed)

	decompressed, err = ZlibDecompress(provider, compressed, "", 0, 0)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, raw, decompressed)
}

func TestRoundtrip(t *testing.T) {
	provider := providers.NewMemoryDictionaryProvider()

	id := "test"
	customVersion := 1
	sharedVersion := 1

	err := provider.AddCustom(id, customVersion, []byte{1, 2})
	ensure.Nil(t, err)
	err = provider.AddShared(sharedVersion, []byte{3, 4})
	ensure.Nil(t, err)

	testBytes := [][]byte{
		{1, 2},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	for _, raw := range testBytes {
		compressed, err := BrotliCompress(provider, raw, id, customVersion, sharedVersion)
		ensure.Nil(t, err)
		ensure.NotDeepEqual(t, raw, compressed)

		decompressed, err := BrotliDecompress(provider, compressed, id, customVersion, sharedVersion)
		ensure.Nil(t, err)
		ensure.DeepEqual(t, raw, decompressed)

		for _, level := range []int{zlib.BestCompression, zlib.BestSpeed, zlib.DefaultCompression} {
			compressed, err = ZlibCompress(provider, level, raw, id, customVersion, sharedVersion)
			ensure.Nil(t, err)
			ensure.NotDeepEqual(t, raw, compressed)

			decompressed, err = ZlibDecompress(provider, compressed, id, customVersion, sharedVersion)
			ensure.Nil(t, err)
			ensure.DeepEqual(t, raw, decompressed)
		}
	}
}

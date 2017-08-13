package examples

// COPY ANY CHANGES BELOW INTO README.md

import (
	"github.com/crawlcoin/ccc"
	"github.com/crawlcoin/ccc/providers"
)

func cccRoundtrip(input []byte) []byte {
	mem := providers.NewMemoryDictionaryProvider()
	id := "test"
	customVersion := 1
	sharedVersion := 1
	mem.AddCustom(id, customVersion, []byte{1, 2})
	mem.AddShared(sharedVersion, []byte{3, 4})

	compressed, _ := ccc.BrotliCompress(mem, input, id, customVersion, sharedVersion)
	decompressed, _ := ccc.BrotliDecompress(mem, compressed, id, customVersion, sharedVersion)
	return decompressed
}


# Crawlcoin Compression

A wrapper around Zlib and Brotli that lazy-loads versioned compression dictionaries over HTTP, file-system, or memory.  Supports both a shared dictionary and a per-ID custom dictionary.  The ID can be anything app-specific, commonly a domain name.

Dictionaries can be manually created or by using [dictator](https://github.com/vkrasnov/dictator).

Part of the [Crawlcoin](https://crawlcoin.com) system but generally usable.

Overview
---

Each dictionary must be static for that id/version.  If a new dictionary is created, the version should be bumped.  This allows for backward compatibility with any existing compressed file while allowing refinement over time.  The version of the dictionary used to compress some bytes is not included in the compressed file and must be stored separately.

If 0 is specified for a particular dictionary, no dictionary is used (including shared).

Usage
---

To use the bindings, you just need to import the ccc package, provider a dictionary provider, and compress/decompress.

Compression + decompression example with no error handling:

```go
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
```

Check out `providers/url_test.go` for complete examples of HTTP dictionaries and caching. 

Development
---

### Testing

`make test` to run all tests locally

`make bench` to run all benchmarks

### Submitting Patches

Before sending commits or patches run 

`make fmt && make fulltest` 

or if you want to `git commit -a` you can run the convenience target 

`make commit`

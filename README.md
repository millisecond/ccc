
## Crawlcoin Compression

Please see the [documentation](https://github.com/Crawlcoin/documentation) for details on Crawlcoin generally. 

### Overview

A wrapper around Brotli compression with custom dictionaries.  Controlling both client and server we can create domain-specific dictionaries that offer much higher compression ratios than a standard approach. 

Each dictionary is static for that version which allows for aggressive caching and accessible by a constant URL scheme:

Examples: 

(https://dictionaries.crawlcoin.com/shared/1.dict)[https://dictionaries.crawlcoin.com/shared/1.dict] is the v1 file of our shared dictionary. 

(https://dictionaries.crawlcoin.com/host/cnn.com/2.dict)[https://dictionaries.crawlcoin.com/host/cnn.com/2.dict] is the v2 file of a dictionary for `cnn.com` this base file will never change - only new versions added.
 
(https://dictionaries.crawlcoin.com/host/cnn.com/dictionary.json)[https://dictionaries.crawlcoin.com/host/cnn.com/dictionary.json] would contain `{"version":2}` in this case. 

If 0 is specified for a particular dictionary, it is ignored (including shared).

## NOTES

Dictionary creation is based on [dictator](https://github.com/vkrasnov/dictator). 

## Development

### Testing

`make test` to run all tests locally

`make bench` to run all benchmarks

### Submitting Patches

Before sending commits or patches run 

`make fmt && make fulltest` 

or if you want to `git commit -a` you can run the convenience target 

`make commit`

package providers

import (
	"strconv"
	"github.com/hashicorp/golang-lru"
	"sync"
	"io/ioutil"
	"net/http"
)

const CACHE_SIZE = 128
const LOCAL_PREFIX = "~/.crawlcoin/dictionaries/"
const REMOTE_PREFIX = "https://dictionaries.crawlcoin.com/"

var cache *lru.Cache
var lruMutex = &sync.RWMutex{}

type URLDictionaryProvider struct {
	BaseSharedURL string
	BaseCustomURL string
	UseMemoryCache    bool
	UseFSCache    bool
}

func (p URLDictionaryProvider) SharedDictionary(version int) ([]byte, error) {
	return nil, nil
}

func (p URLDictionaryProvider) CustomDictionary(id string, version int) ([]byte, error) {
	return nil, nil
}

func (p URLDictionaryProvider) fileOrRemote(relativePath string, version int) ([]byte, error) {
	var dict []byte
	var err error
	localFilename := LOCAL_PREFIX + relativePath + "/" + strconv.Itoa(version) + ".dict"
	if p.UseMemoryCache {
		dict, err = func() ([]byte, error) {
			err := initCache()
			if err != nil {
				return nil, err
			}
			c, pres := cache.Get(localFilename)
			if !pres {
				return nil, nil
			}
			return c.([]byte), nil
		}()
	}
	if len(dict) > 0 {
		return dict, nil
	}
	if p.UseFSCache {
		if dict, err = ioutil.ReadFile(localFilename); err != nil {
			return nil, err
		}
	}
	if p.UseMemoryCache && len(dict) > 0 {
		cache.Add(localFilename, dict)
		return dict, nil
	}
	resp, err := http.Get(REMOTE_PREFIX + relativePath)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if dict, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	}
	if len(dict) > 0 {
		if p.UseMemoryCache {
			cache.Add(localFilename, dict)
		}
		if p.UseFSCache {
			// Write it out to cache for next time.
			ioutil.WriteFile(localFilename, dict, 0644)
		}
	}
	return dict, nil
}

func initCache() error {
	if cache != nil {
		return nil
	}
	lruMutex.Lock()
	defer lruMutex.Unlock()
	if cache != nil {
		return nil
	}
	var err error
	cache, err = lru.New(CACHE_SIZE)
	if err != nil {
		return err
	}
	return nil
}

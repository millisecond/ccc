package ccc

import (
	"github.com/hashicorp/golang-lru"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

const DEFAULT_CACHE_SIZE = 512

type URLDictionaryProvider struct {
	BaseSharedURL string
	BaseCustomURL string

	UseMemoryCache  bool
	MemoryCacheSize int

	UseFSCache     bool
	FSSharedPrefix string
	FSCustomPrefix string

	mut    sync.RWMutex
	lruMut sync.RWMutex
	cache  *lru.Cache
}

// Create a basic URL provider, all requests will go over HTTP
func NewURLDictionaryProvider(BaseSharedURL string, BaseCustomURL string) *URLDictionaryProvider {
	return &URLDictionaryProvider{
		BaseSharedURL: BaseSharedURL,
		BaseCustomURL: BaseCustomURL,
	}
}

// More complicated constructor with caching options
func NewCachedURLDictionaryProvider(BaseSharedURL string, BaseCustomURL string, UseMemoryCache bool, MemoryCacheSize int, UseFSCache bool, FSSharedPrefix string, FSCustomPrefix string) (*URLDictionaryProvider, error) {
	p := NewURLDictionaryProvider(BaseSharedURL, BaseCustomURL)
	p.UseMemoryCache = UseMemoryCache
	p.MemoryCacheSize = MemoryCacheSize
	p.UseFSCache = UseFSCache
	p.FSSharedPrefix = FSSharedPrefix
	p.FSCustomPrefix = FSCustomPrefix
	if p.UseMemoryCache {
		p.lruMut.Lock()
		defer p.lruMut.Unlock()
		if p.MemoryCacheSize == 0 {
			p.MemoryCacheSize = DEFAULT_CACHE_SIZE
		}
		var err error
		p.cache, err = lru.New(p.MemoryCacheSize)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (p *URLDictionaryProvider) SharedDictionary(version int) ([]byte, error) {
	return p.fileOrRemote(p.BaseSharedURL, "", true, version)
}

func (p *URLDictionaryProvider) CustomDictionary(id string, version int) ([]byte, error) {
	return p.fileOrRemote(p.BaseCustomURL, id+"/", false, version)
}

func (p *URLDictionaryProvider) fileOrRemote(base string, path string, shared bool, version int) ([]byte, error) {
	var dict []byte
	var err error
	localFilename := p.FSCustomPrefix
	if shared {
		localFilename = p.FSSharedPrefix
	}
	localFilename += path + strconv.Itoa(version) + ".dict"
	if p.UseMemoryCache {
		dict, err = func() ([]byte, error) {
			c, pres := p.cache.Get(localFilename)
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
		dict, _ = ioutil.ReadFile(localFilename)
		// ignore errors, fallback to URL
		if len(dict) > 0 {
			return dict, nil
		}
	}
	if p.UseMemoryCache && len(dict) > 0 {
		p.cache.Add(localFilename, dict)
		return dict, nil
	}
	url := base + path + strconv.Itoa(version)
	resp, err := http.Get(url)
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
			p.cache.Add(localFilename, dict)
		}
		if p.UseFSCache {
			// Write it out to cache for next time.
			ioutil.WriteFile(localFilename, dict, 0644)
		}
	}
	return dict, nil
}

package ccc

import (
	"errors"
	"github.com/hashicorp/golang-lru"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

const CACHE_SIZE = 128
const LOCAL_PREFIX = "~/.crawlcoin/dictionaries/"
const REMOTE_PREFIX = "https://dictionaries.crawlcoin.com/"
const SHARED = "shared"

var INVALID_HOSTS = map[string]bool{"shared": true}

var cache *lru.Cache
var lruMutex = &sync.RWMutex{}

var localOverride func(path string, version int) []byte
var remoteOverride func(path string, version int) []byte

type Version struct {
	Version int
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

func Combined(host string, version int, sharedVersion int) ([]byte, error) {
	hostDict, err := HostDictionary(host, version)
	if err != nil {
		return nil, err
	}
	sharedDict, err := SharedDictionary(sharedVersion)
	if err != nil {
		return nil, err
	}
	return Append(hostDict, sharedDict), nil
}

func Append(dict1 []byte, dict2 []byte) []byte {
	return append(dict1, dict2...)
}

func SharedDictionary(version int) ([]byte, error) {
	return fileOrRemote(SHARED, version)
}

func HostDictionary(host string, version int) ([]byte, error) {
	if _, prs := INVALID_HOSTS[host]; prs {
		return nil, errors.New("Invalid host, some names are reserved.")
	}
	return fileOrRemote(host, version)
}

func fileOrRemote(relativePath string, version int) ([]byte, error) {
	var dict []byte
	var err error
	localFilename := LOCAL_PREFIX + relativePath + "/" + strconv.Itoa(version) + ".dict"
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
	if len(dict) > 0 {
		return dict, nil
	}
	if localOverride == nil {
		if dict, err = ioutil.ReadFile(localFilename); err != nil {
			return nil, err
		}
	} else {
		dict = localOverride(relativePath, version)
	}
	if len(dict) > 0 {
		cache.Add(localFilename, dict)
		return dict, nil
	}
	if remoteOverride == nil {
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
	} else {
		dict = remoteOverride(relativePath, version)
	}
	if len(dict) > 0 {
		cache.Add(localFilename, dict)
		// Write it out to cache for next time.
		ioutil.WriteFile(localFilename, dict, 0644)
	}
	return dict, nil
}

func Create(sample []byte) string {
	return ""
}

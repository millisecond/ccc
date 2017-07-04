package ccz

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"errors"
	"encoding/json"
)

const LOCAL_PREFIX = "~/.crawlcoin/dictionaries/"
const REMOTE_PREFIX = "https://dictionaries.crawlcoin.com/"
const SHARED = "shared"
var INVALID_HOSTS = map[string]bool{"shared":true}

type Version struct {
	Version int
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

func SharedDictionaryLatest() (int, []byte, error) {
	return 0, nil, nil
}

func HostDictionary(host string, version int) ([]byte, error) {
	if _, prs := INVALID_HOSTS[host]; prs {
		return nil, errors.New("Invalid host, some names are reserved.")
	}
	return fileOrRemote(host, version)
}

func HostDictionaryLatest(host string) (int, []byte, error) {
	if _, prs := INVALID_HOSTS[host]; prs {
		return 0, nil, errors.New("Invalid host, some names are reserved.")
	}
	return 0, nil, nil
}

func latestVersion(relativePath string) (int, error) {
	remote, err := http.Get(REMOTE_PREFIX + relativePath)
	if remote != nil && remote.Body != nil {
		defer remote.Body.Close()
	}
	if err != nil {
		return 0, err
	}
	remoteBytes, err := ioutil.ReadAll(remote.Body)
	if err != nil {
		return 0, err
	}
	v := &Version{}
	err = json.Unmarshal(remoteBytes, v)
	if err != nil {
		return 0, err
	}
	return v.Version, nil
}

func fileOrRemote(relativePath string, version int) ([]byte, error) {
	localFilename := LOCAL_PREFIX + relativePath + "/" + strconv.Itoa(version) + ".dict"
	local, err := ioutil.ReadFile(localFilename)
	if err != nil {
		return nil, err
	}
	if len(local) > 0 {
		return local, nil
	}
	remote, err := http.Get(REMOTE_PREFIX + relativePath)
	if remote != nil && remote.Body != nil {
		defer remote.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	remoteBytes, err := ioutil.ReadAll(remote.Body)
	if err != nil {
		return nil, err
	}
	// Write it out to cache for next time.
	ioutil.WriteFile(localFilename, remoteBytes, 0644)
	return remoteBytes, nil
}

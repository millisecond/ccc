package ccc

import (
	"github.com/facebookgo/ensure"
	"testing"
)

func TestAppend(t *testing.T) {
	ensure.DeepEqual(t, Append([]byte{1, 2}, []byte{3, 4}), []byte{1, 2, 3, 4})
}

func TestVersions(t *testing.T) {
	localCalls := 0
	remoteCalls := 0
	localOverride = func(path string, version int) []byte {
		localCalls += 1
		if version != 1 {
			return nil
		}
		return []byte("LOCAL")
	}

	remoteOverride = func(path string, version int) []byte {
		remoteCalls += 1
		return []byte("REMOTE")
	}

	combined, err := Combined("", 1, 2)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, string(combined), "LOCALREMOTE")

	ensure.DeepEqual(t, localCalls, 2)
	ensure.DeepEqual(t, remoteCalls, 1)

	// Test LRU cache
	combined, err = Combined("", 1, 2)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, string(combined), "LOCALREMOTE")

	ensure.DeepEqual(t, localCalls, 2)
	ensure.DeepEqual(t, remoteCalls, 1)
}

func TestZeroDictVersion(t *testing.T) {
	shared, err := SharedDictionary(0)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, shared, []byte{})

	custom, err := HostDictionary("crawlcoin.com", 0)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, custom, []byte{})
}

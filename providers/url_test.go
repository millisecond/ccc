package providers

import (
	"github.com/facebookgo/ensure"
	"net/http"
	"sync/atomic"
	"testing"
)

var sharedCount = int32(0)

func init() {
	http.HandleFunc("/shared/1", func(rw http.ResponseWriter, req *http.Request) {
		atomic.AddInt32(&sharedCount, 1)
		rw.Write([]byte{3, 4})
	})
	http.HandleFunc("/custom/test/1", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte{1, 2})
	})
	go http.ListenAndServe(":7080", nil)

}

func TestURLNotFound(t *testing.T) {
	urlProvider := NewURLDictionaryProvider("", "")

	_, err := urlProvider.CustomDictionary("test", 1)
	ensure.NotNil(t, err)

	_, err = urlProvider.SharedDictionary(1)
	ensure.NotNil(t, err)
}

func TestURLGeneral(t *testing.T) {
	urlProvider := NewURLDictionaryProvider(
		"http://localhost:7080/shared/",
		"http://localhost:7080/custom/",
	)

	custom, err := urlProvider.CustomDictionary("test", 1)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, custom, []byte{1, 2})

	ensure.DeepEqual(t, sharedCount, int32(0))
	shared, err := urlProvider.SharedDictionary(1)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, shared, []byte{3, 4})
	ensure.DeepEqual(t, sharedCount, int32(1))

	shared, err = urlProvider.SharedDictionary(1)
	ensure.DeepEqual(t, sharedCount, int32(2))
}

func TestURLMemCache(t *testing.T) {
	sharedCount = int32(0)

	urlProvider, err := NewCachedURLDictionaryProvider(
		"http://localhost:7080/shared/",
		"http://localhost:7080/custom/",
		true, 128,
		false, "", "",
	)
	ensure.Nil(t, err)

	ensure.DeepEqual(t, sharedCount, int32(0))
	shared, err := urlProvider.SharedDictionary(1)
	ensure.DeepEqual(t, shared, []byte{3, 4})
	ensure.Nil(t, err)
	ensure.DeepEqual(t, sharedCount, int32(1))

	shared, err = urlProvider.SharedDictionary(1)
	ensure.DeepEqual(t, shared, []byte{3, 4})
	ensure.Nil(t, err)
	ensure.DeepEqual(t, sharedCount, int32(1))
}

func TestURLFSCache(t *testing.T) {
	sharedCount = int32(0)
	sharedDir, customDir := emptyTestDirs(t)

	urlProvider, err := NewCachedURLDictionaryProvider(
		"http://localhost:7080/shared/",
		"http://localhost:7080/custom/",
		false, 0,
		true, sharedDir, customDir,
	)
	ensure.Nil(t, err)

	ensure.DeepEqual(t, sharedCount, int32(0))
	shared, err := urlProvider.SharedDictionary(1)
	ensure.DeepEqual(t, shared, []byte{3, 4})
	ensure.Nil(t, err)
	ensure.DeepEqual(t, sharedCount, int32(1))

	shared, err = urlProvider.SharedDictionary(1)
	ensure.DeepEqual(t, shared, []byte{3, 4})
	ensure.Nil(t, err)
	ensure.DeepEqual(t, sharedCount, int32(1))
}

package providers

import (
	"github.com/facebookgo/ensure"
	"io/ioutil"
	"os"
	"testing"
)

func TestFileNotFound(t *testing.T) {
	sharedDir := os.TempDir()
	customDir := os.TempDir()
	fileProvider := NewFileDictionaryProvider(sharedDir, customDir)

	_, err := fileProvider.CustomDictionary("test", 1)
	ensure.NotNil(t, err)

	_, err = fileProvider.SharedDictionary(1)
	ensure.NotNil(t, err)
}

func TestFileGeneral(t *testing.T) {
	sharedDir := os.TempDir()
	customDir := os.TempDir()
	fileProvider := NewFileDictionaryProvider(sharedDir, customDir)

	os.MkdirAll(customDir+"test/", 0755)
	err := ioutil.WriteFile(customDir+"test/1.dict", []byte{1, 2}, 0755)
	ensure.Nil(t, err)
	err = ioutil.WriteFile(sharedDir+"1.dict", []byte{3, 4}, 0755)
	ensure.Nil(t, err)

	custom, err := fileProvider.CustomDictionary("test", 1)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, custom, []byte{1, 2})

	shared, err := fileProvider.SharedDictionary(1)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, shared, []byte{3, 4})

	shared, err = fileProvider.SharedDictionary(1)
}

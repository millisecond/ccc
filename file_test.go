package ccc

import (
	"github.com/facebookgo/ensure"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFileNotFound(t *testing.T) {
	sharedDir, customDir := emptyTestDirs(t)

	fileProvider := NewFileDictionaryProvider(sharedDir, customDir)

	_, err := fileProvider.CustomDictionary("test", 22)
	ensure.NotNil(t, err)

	_, err = fileProvider.SharedDictionary(22)
	ensure.NotNil(t, err)
}

func TestFileGeneral(t *testing.T) {
	sharedDir, customDir := emptyTestDirs(t)

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

func emptyTestDirs(t *testing.T) (string, string) {
	sharedDir := os.TempDir() + "shared/"
	customDir := os.TempDir() + "custom/"
	err := removeContents(sharedDir)
	ensure.Nil(t, err)
	err = removeContents(customDir)
	ensure.Nil(t, err)
	err = os.MkdirAll(sharedDir, 0755)
	ensure.Nil(t, err)
	err = os.MkdirAll(customDir, 0755)
	ensure.Nil(t, err)
	return sharedDir, customDir
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

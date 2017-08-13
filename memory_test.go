package ccc

import (
	"github.com/facebookgo/ensure"
	"testing"
)

func TestMemoryNotFound(t *testing.T) {
	mem := &MemoryDictionaryProvider{}

	_, err := mem.CustomDictionary("test", 1)
	ensure.NotNil(t, err)

	_, err = mem.SharedDictionary(1)
	ensure.NotNil(t, err)
}

func TestMemoryGeneral(t *testing.T) {
	mem := NewMemoryDictionaryProvider()

	err := mem.AddCustom("test", 1, []byte{1, 2})
	ensure.Nil(t, err)
	err = mem.AddShared(1, []byte{3, 4})
	ensure.Nil(t, err)

	custom, err := mem.CustomDictionary("test", 1)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, custom, []byte{1, 2})

	shared, err := mem.SharedDictionary(1)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, shared, []byte{3, 4})
}

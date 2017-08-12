package providers

import (
	"fmt"
	"io/ioutil"
	"strconv"
)

type FileDictionaryProvider struct {
	BaseSharedDirectory string
	BaseCustomDirectory string
}

func NewFileDictionaryProvider(BaseSharedDirectory string, BaseCustomDirectory string) *FileDictionaryProvider {
	return &FileDictionaryProvider{
		BaseCustomDirectory: BaseCustomDirectory,
		BaseSharedDirectory: BaseSharedDirectory,
	}
}

func (p FileDictionaryProvider) SharedDictionary(version int) ([]byte, error) {
	return p.file(p.BaseSharedDirectory, "", true, version)
}

func (p FileDictionaryProvider) CustomDictionary(id string, version int) ([]byte, error) {
	return p.file(p.BaseCustomDirectory, id+"/", false, version)
}

func (p *FileDictionaryProvider) file(base string, path string, shared bool, version int) ([]byte, error) {
	localFilename := p.BaseCustomDirectory
	if shared {
		localFilename = p.BaseSharedDirectory
	}
	localFilename += path + strconv.Itoa(version) + ".dict"
	fmt.Println(localFilename)
	return ioutil.ReadFile(localFilename)
}

package providers

import (
	"errors"
	"sync"
)

type MemoryDictionaryProvider struct {
	SharedDictionaries map[int][]byte
	CustomDictionaries map[string]map[int][]byte
	mut                sync.RWMutex
}

func NewMemoryDictionaryProvider() *MemoryDictionaryProvider {
	p := &MemoryDictionaryProvider{}
	p.CustomDictionaries = make(map[string]map[int][]byte)
	p.SharedDictionaries = make(map[int][]byte)
	return p
}

func (p *MemoryDictionaryProvider) AddCustom(id string, version int, dict []byte) error {
	p.mut.Lock()
	defer p.mut.Unlock()
	if p.CustomDictionaries == nil {
		return errors.New("Memory provider not init'd")
	}
	if custom, pres := p.CustomDictionaries[id]; pres {
		custom[version] = dict
	} else {
		p.CustomDictionaries[id] = map[int][]byte{version: dict}
	}
	return nil
}

func (p *MemoryDictionaryProvider) AddShared(version int, dict []byte) error {
	p.mut.Lock()
	defer p.mut.Unlock()
	if p.SharedDictionaries == nil {
		return errors.New("Memory provider not init'd")
	}
	p.SharedDictionaries[version] = dict
	return nil
}

func (p *MemoryDictionaryProvider) SharedDictionary(version int) ([]byte, error) {
	p.mut.RLock()
	defer p.mut.RUnlock()
	if dict, pres := p.SharedDictionaries[version]; pres {
		return dict, nil
	}
	return nil, errors.New("No version found for shared dictionary.")
}

func (p *MemoryDictionaryProvider) CustomDictionary(id string, version int) ([]byte, error) {
	p.mut.RLock()
	defer p.mut.RUnlock()
	if customs, pres := p.CustomDictionaries[id]; pres {
		if dict, pres := customs[version]; pres {
			return dict, nil
		}
	}
	return nil, errors.New("No version found for custom dictionary.")
}

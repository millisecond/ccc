package providers

import "errors"

type MemoryDictionaryProvider struct {
	SharedDictionaries map[int][]byte
	CustomDictionaries map[string]map[int][]byte
}

func (p MemoryDictionaryProvider) SharedDictionary(version int) ([]byte, error) {
	if dict, pres := p.SharedDictionaries[version]; pres {
		return dict, nil
	}
	return nil, errors.New("No version found for shared dictionary.")
}

func (p MemoryDictionaryProvider) CustomDictionary(id string, version int) ([]byte, error) {
	if customs, pres := p.CustomDictionaries[id]; pres {
		if dict, pres := customs[version]; pres {
			return dict, nil
		}
	}
	return nil, errors.New("No version found for custom dictionary.")
}

package ccc

type DictionaryProvider interface {
	// Return the bytes for a shared dictionary at a specific version
	SharedDictionary(version int) ([]byte, error)

	// Return the vytes for a custom dictionary at a specific version
	CustomDictionary(version int) ([]byte, error)
}

type URLDictionaryProvider struct {
	BaseSharedURL string
	BaseCustomURL string
}

func (p *URLDictionaryProvider) SharedDictionary(version int) ([]byte, error) {
	return nil, nil
}

func (p *URLDictionaryProvider) CustomDictionary(version int) ([]byte, error) {
	return nil, nil
}

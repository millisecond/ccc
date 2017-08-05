package providers

type FileDictionaryProvider struct {
	BaseSharedURL string
	BaseCustomURL string
}

func (p FileDictionaryProvider) SharedDictionary(version int) ([]byte, error) {
	return nil, nil
}

func (p FileDictionaryProvider) CustomDictionary(id string, version int) ([]byte, error) {
	return nil, nil
}

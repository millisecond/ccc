package providers

type DictionaryProvider interface {
	// Return the bytes for a shared dictionary at a specific version
	SharedDictionary(version int) ([]byte, error)

	// Return the vytes for a custom dictionary at a specific version
	CustomDictionary(id string, version int) ([]byte, error)
}

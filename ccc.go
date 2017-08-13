package ccc

import (
	"bytes"
	"compress/zlib"
	"gopkg.in/kothar/brotli-go.v0/dec"
	"gopkg.in/kothar/brotli-go.v0/enc"
	"io"
	"io/ioutil"
)

// Given a dictionary provider, compress some bytes with Brotli using versioned dictionaries.  Use 0 to ignore that type of dictionary.
func BrotliCompress(provider DictionaryProvider, b []byte, id string, customVersion int, sharedVersion int) ([]byte, error) {
	var err error
	var encoded []byte
	dict, err := Combined(provider, id, customVersion, sharedVersion)
	if err != nil {
		return nil, err
	}
	if dict != nil && len(dict) > 0 {
		encoded, err = enc.CompressBufferDict(nil, b, dict, nil)
	} else {
		encoded, err = enc.CompressBuffer(nil, b, nil)
	}
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

// Given a dictionary provider, decompress some bytes with Brotli using versioned dictionaries.  Use 0 to ignore that type of dictionary.
func BrotliDecompress(provider DictionaryProvider, b []byte, id string, customVersion int, sharedVersion int) ([]byte, error) {
	var err error
	var decoded []byte
	dict, err := Combined(provider, id, customVersion, sharedVersion)
	if err != nil {
		return nil, err
	}
	if dict != nil && len(dict) > 0 {
		decoded, err = dec.DecompressBufferDict(b, dict, nil)
	} else {
		decoded, err = dec.DecompressBuffer(b, nil)
	}
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

// Given a dictionary provider, compress some bytes with zlib using versioned dictionaries.  Use 0 to ignore that type of dictionary.
func ZlibCompress(provider DictionaryProvider, level int, b []byte, id string, customVersion int, sharedVersion int) ([]byte, error) {
	var err error
	dict, err := Combined(provider, id, customVersion, sharedVersion)
	if err != nil {
		return nil, err
	}
	var w *zlib.Writer
	var zipOut bytes.Buffer
	if dict != nil && len(dict) > 0 {
		w, err = zlib.NewWriterLevelDict(&zipOut, level, dict)
	} else {
		w, err = zlib.NewWriterLevel(&zipOut, level)
	}
	if err != nil {
		return nil, err
	}
	_, err = w.Write(b)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return zipOut.Bytes(), nil
}

// Given a dictionary provider, decompress some bytes with zlib using versioned dictionaries.  Use 0 to ignore that type of dictionary.
func ZlibDecompress(provider DictionaryProvider, b []byte, id string, customVersion int, sharedVersion int) ([]byte, error) {
	var err error
	dict, err := Combined(provider, id, customVersion, sharedVersion)
	if err != nil {
		return nil, err
	}
	var r io.ReadCloser
	if dict != nil && len(dict) > 0 {
		r, err = zlib.NewReaderDict(bytes.NewReader(b), dict)
	} else {
		r, err = zlib.NewReader(bytes.NewReader(b))
	}
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}

package ccc

import (
	"gopkg.in/kothar/brotli-go.v0/enc"
	"gopkg.in/kothar/brotli-go.v0/dec"
	"github.com/crawlcoin/ccc/providers"
)

func Compress(provider providers.DictionaryProvider, b []byte, id string, customVersion int, sharedVersion int) ([]byte, error) {
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

func Decompress(provider providers.DictionaryProvider, b []byte, id string, customVersion int, sharedVersion int) ([]byte, error) {
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

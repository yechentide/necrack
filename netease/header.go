package netease

import (
	"bytes"
	"fmt"
)

var (
	headerNetEaseCurrent = []byte{0x80, 0x1D, 0x30, 0x01}
	headerNetEaseLegacy  = []byte{0x90, 0x1D, 0x30, 0x01}
	headerVanillaBedrock = []byte{0x4D, 0x41, 0x4E, 0x49} // "MANI"
)

type HeaderType int

const (
	HeaderTypeNetEaseCurrent HeaderType = iota
	HeaderTypeNetEaseLegacy
	HeaderTypeVanillaBedrock
	HeaderTypeUnknown
)

func identifyHeader(data []byte) HeaderType {
	if len(data) < 4 {
		return HeaderTypeUnknown
	}

	header := data[:4]

	if bytes.Equal(header, headerNetEaseCurrent) {
		return HeaderTypeNetEaseCurrent
	}
	if bytes.Equal(header, headerNetEaseLegacy) {
		return HeaderTypeNetEaseLegacy
	}
	if bytes.Equal(header, headerVanillaBedrock) {
		return HeaderTypeVanillaBedrock
	}

	return HeaderTypeUnknown
}

func ValidateDecryptableFile(data []byte) error {
	headerType := identifyHeader(data)

	switch headerType {
	case HeaderTypeNetEaseCurrent:
		return nil
	case HeaderTypeNetEaseLegacy:
		return fmt.Errorf("legacy NetEase encryption (AES-CFB8) is not supported")
	case HeaderTypeVanillaBedrock:
		return fmt.Errorf("vanilla Bedrock MANIFEST format, no decryption needed")
	default:
		return fmt.Errorf("unknown or invalid header format")
	}
}

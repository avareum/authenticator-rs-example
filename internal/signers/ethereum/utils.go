package ethereum

import (
	hexlib "encoding/hex"
	"errors"
	"strings"
)

func GetMethodSignature(hex string) ([]byte, error) {
	hex = strings.TrimPrefix(hex, "0x")
	if len(hex) < 8 {
		return nil, errors.New("invalid hex")
	}
	return hexlib.DecodeString(hex[:8])
}

func GetInputData(hex string) ([]byte, error) {
	hex = strings.TrimPrefix(hex, "0x")
	if len(hex) <= 8 {
		return nil, errors.New("invalid hex")
	}
	return hexlib.DecodeString(hex[8:])
}

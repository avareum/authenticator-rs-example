package ethereum

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
)

func GetMethodSignature(bytecode string) ([]byte, error) {
	if len(bytecode) < 10 {
		return nil, errors.New("invalid bytecode")
	}
	return common.Hex2Bytes(bytecode[2:10]), nil
}

func GetInputData(bytecode string) ([]byte, error) {
	if len(bytecode) <= 10 {
		return nil, errors.New("invalid bytecode")
	}
	return common.Hex2Bytes(bytecode[10:]), nil
}

package ethereum

import (
	"errors"
)

func MethodSignature(bytecode string) (string, error) {
	if len(bytecode) < 10 {
		return "", errors.New("invalid bytecode")
	}
	return bytecode[2:10], nil
}

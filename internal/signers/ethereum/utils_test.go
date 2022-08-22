package ethereum

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/test-go/testify/require"
)

const FTT = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"burnFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

func Test_ExtractBytecode(t *testing.T) {
	t.Run("should extract input data", func(t *testing.T) {
		// FTT token transfered
		// https://etherscan.io/tx/0x9ce54ebd9b60c64e675885b31c4230f514b0574d2989709d5c602fdc6f4dc758
		bytecode := "0xa9059cbb000000000000000000000000f16e9b0d03470827a95cdfd0cb8a8a3b46969b91000000000000000000000000000000000000000000000000301264884eb2c000"
		erc20, err := abi.JSON(strings.NewReader(FTT))
		require.Nil(t, err)

		// Extract method signature
		sig, err := GetMethodSignature(bytecode)
		require.Nil(t, err)
		method, err := erc20.MethodById(sig)
		require.Nil(t, err)

		// Extract input data
		inputData, err := GetInputData(bytecode)
		require.Nil(t, err)
		inputs, err := method.Inputs.Unpack(inputData)
		require.Nil(t, err)

		require.Equal(t, common.HexToAddress("0xf16E9B0D03470827A95CDfd0Cb8a8A3b46969B91"), inputs[0])
		require.Equal(t, big.NewInt(3463941600000000000), inputs[1])
	})
}

func Test_ExtractMethodSignature(t *testing.T) {
	t.Run("should extract method signature", func(t *testing.T) {
		inputData := "0xa9059cbb000000000000000000000000f16e9b0d03470827a95cdfd0cb8a8a3b46969b91000000000000000000000000000000000000000000000000301264884eb2c000"
		sig, err := GetMethodSignature(inputData)
		require.Nil(t, err)
		require.Equal(t, common.Hex2Bytes("a9059cbb"), sig)
	})

	t.Run("should reject invalid bytecode", func(t *testing.T) {
		inputData := "0x1"
		_, err := GetMethodSignature(inputData)
		require.EqualError(t, err, "invalid bytecode")
	})
}

func Test_ExtractInputData(t *testing.T) {
	t.Run("should extract method signature", func(t *testing.T) {
		inputData := "0xa9059cbb000000000000000000000000f16e9b0d03470827a95cdfd0cb8a8a3b46969b91000000000000000000000000000000000000000000000000301264884eb2c000"
		sig, err := GetInputData(inputData)
		require.Nil(t, err)
		require.Equal(t, common.Hex2Bytes("000000000000000000000000f16e9b0d03470827a95cdfd0cb8a8a3b46969b91000000000000000000000000000000000000000000000000301264884eb2c000"), sig)
	})

	t.Run("should reject invalid bytecode", func(t *testing.T) {
		inputData := "0x1"
		_, err := GetInputData(inputData)
		require.EqualError(t, err, "invalid bytecode")
	})
}

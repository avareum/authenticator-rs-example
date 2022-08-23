package ethereum

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/test-go/testify/require"
)

func Test_ExtractMethodSignature(t *testing.T) {
	rawInputDataHex := "0xa9059cbb000000000000000000000000f16e9b0d03470827a95cdfd0cb8a8a3b46969b91000000000000000000000000000000000000000000000000301264884eb2c000"
	expectedSignatures := common.Hex2Bytes("a9059cbb")

	t.Run("should extract signature from hex", func(t *testing.T) {
		sig, err := GetMethodSignature(rawInputDataHex)
		require.Nil(t, err)
		require.Equal(t, expectedSignatures, sig)
	})

	t.Run("should extract signature from hex data (without 0x prefix)", func(t *testing.T) {
		// this case found while extracting ethereum transaction.Data()
		sig, err := GetMethodSignature(rawInputDataHex[2:])
		require.Nil(t, err)
		require.Equal(t, expectedSignatures, sig)
	})

	t.Run("should reject invalid hex", func(t *testing.T) {
		inputData := "0x1"
		_, err := GetMethodSignature(inputData)
		require.EqualError(t, err, "invalid hex")
	})

}

func Test_ExtractInputData(t *testing.T) {
	rawInputDataHex := "0xa9059cbb000000000000000000000000f16e9b0d03470827a95cdfd0cb8a8a3b46969b91000000000000000000000000000000000000000000000000301264884eb2c000"
	expectedInputData := common.Hex2Bytes("000000000000000000000000f16e9b0d03470827a95cdfd0cb8a8a3b46969b91000000000000000000000000000000000000000000000000301264884eb2c000")

	t.Run("should extract input data", func(t *testing.T) {
		sig, err := GetInputData(rawInputDataHex)
		require.Nil(t, err)
		require.Equal(t, expectedInputData, sig)
	})

	t.Run("should extract input data (without 0x prefix)", func(t *testing.T) {
		// this case found while extracting ethereum transaction.Data()
		sig, err := GetInputData(rawInputDataHex[2:])
		require.Nil(t, err)
		require.Equal(t, expectedInputData, sig)
	})

	t.Run("should reject invalid hex", func(t *testing.T) {
		rawInputDataHex := "0x1"
		_, err := GetInputData(rawInputDataHex)
		require.EqualError(t, err, "invalid hex")
	})
}

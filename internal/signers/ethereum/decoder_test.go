package ethereum

import (
	"math/big"
	"strings"
	"testing"

	"github.com/avareum/avareum-hubble-signer/internal/signers/ethereum/types"
	"github.com/avareum/avareum-hubble-signer/tests/fixtures"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/test-go/testify/require"
)

const FTT = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"burnFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

func Test_EthereumTransactionDecoder(t *testing.T) {

	// FTT token transfered
	// https://etherscan.io/getRawTx?tx=0x9ce54ebd9b60c64e675885b31c4230f514b0574d2989709d5c602fdc6f4dc758
	txRawHex := "0x02f8b10180843b9aca008507ebabd723830186a09450d1c9771902476076ecfc8b2a83ad6b9355a4c980b844a9059cbb000000000000000000000000f16e9b0d03470827a95cdfd0cb8a8a3b46969b91000000000000000000000000000000000000000000000000301264884eb2c000c001a01dd61519d498e5bbf86bf722e616bad59a48bbbe8cf9ed0183698512a4ea7bb7a076f73fa82830d251af9b11de00069e74040cad3e4dfd062c62921191848a3849"

	suite := fixtures.NewTestSuite()
	decoder := NewEthereumTransactionDecoder()
	sender := types.MustNewEthereumKey()
	receiver := types.MustNewEthereumKey()
	originalTx := suite.Ethereum.NewTransferTransaction(*sender, receiver.PublicKey, 1)

	t.Run("ethereum tx decoder", func(t *testing.T) {

		t.Run("should decode from bin", func(t *testing.T) {
			bin, err := originalTx.MarshalBinary()
			require.Nil(t, err)

			tx, err := decoder.TryDecode(bin)
			require.Nil(t, err)
			require.NotNil(t, tx)
			require.Equal(t, originalTx.Hash(), tx.Hash())
			require.Equal(t, originalTx.To(), tx.To())
		})

		t.Run("should decode from JSON", func(t *testing.T) {
			bin, err := originalTx.MarshalJSON()
			require.Nil(t, err)

			tx, err := decoder.TryDecode(bin)
			require.Nil(t, err)
			require.NotNil(t, tx)
			require.Equal(t, originalTx.Hash(), tx.Hash())
			require.Equal(t, originalTx.To(), tx.To())
		})

		t.Run("should decode from raw hex", func(t *testing.T) {
			tx, err := decoder.TryDecodeFromHex(txRawHex)
			require.Nil(t, err)
			require.NotNil(t, tx)
			require.Equal(t, "0x9ce54ebd9b60c64e675885b31c4230f514b0574d2989709d5c602fdc6f4dc758", tx.Hash().String())
			require.Equal(t, "0x50D1c9771902476076eCFc8B2A83Ad6b9355a4c9", tx.To().String())
		})
	})

	t.Run("ethereum tx input data decoder", func(t *testing.T) {
		tx, err := decoder.TryDecodeFromHex(txRawHex)
		require.Nil(t, err)

		t.Run("should decode signature", func(t *testing.T) {
			sig, err := decoder.GetMethodSignature(tx)
			require.Nil(t, err)
			require.Equal(t, common.Hex2Bytes("a9059cbb"), sig)
		})

		t.Run("should decode input data", func(t *testing.T) {
			abi, err := abi.JSON(strings.NewReader(FTT))
			require.Nil(t, err)
			inputs, err := decoder.GetInputData(abi, tx)
			require.Nil(t, err)
			require.Equal(t, common.HexToAddress("0xf16E9B0D03470827A95CDfd0Cb8a8A3b46969B91"), inputs[0])
			require.Equal(t, big.NewInt(3463941600000000000), inputs[1])
		})
	})
}

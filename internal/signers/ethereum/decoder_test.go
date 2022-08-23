package ethereum

import (
	"testing"

	"github.com/avareum/avareum-hubble-signer/internal/signers/ethereum/types"
	"github.com/avareum/avareum-hubble-signer/tests/fixtures"
	"github.com/test-go/testify/require"
)

func Test_EthereumTransactionDecoder(t *testing.T) {
	t.Run("ethereum decoder", func(t *testing.T) {
		suite := fixtures.NewTestSuite()
		decoder := NewEthereumTransactionDecoder()
		sender := types.MustNewEthereumKey()
		receiver := types.MustNewEthereumKey()
		originalTx := suite.Ethereum.NewTransferTransaction(*sender, receiver.PublicKey, 1)

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
			// https://etherscan.io/getRawTx?tx=0x0e5d12200272547fe6f14a7c559ca7ebc2f58123dc84268d761ab20c1cbbb175
			rawHex := "0x02f8b401833d8d6e84773594008517bfac7c00830329189450d1c9771902476076ecfc8b2a83ad6b9355a4c980b844a9059cbb00000000000000000000000043c0199357fa579a969a5296bc5fc26999b027870000000000000000000000000000000000000000000000087bd94bbcac098000c080a0022365a7725ff27cd52898f82332ff03b47c314b9229eec2b5453d03b301201da044eabf05d280d46dd50e913e104ef10e97661fbc9a1870631fb59442ca57dd90"

			tx, err := decoder.TryDecodeFromHex(rawHex)
			require.Nil(t, err)
			require.NotNil(t, tx)
			require.Equal(t, "0x0e5d12200272547fe6f14a7c559ca7ebc2f58123dc84268d761ab20c1cbbb175", tx.Hash().String())
			require.Equal(t, "0x50D1c9771902476076eCFc8B2A83Ad6b9355a4c9", tx.To().String())
		})
	})
}

package solana

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testSigner() *SolanaSigner {
	return NewSolanaSigner(SolanaSignerOptions{
		RPC:       "https://api.devnet.solana.com",
		Websocket: "wss://api.devnet.solana.com",
	})
}

// use inspector to decode the tx
// https://explorer.solana.com/tx/inspector?cluster=mainnet-beta
const testMsg = "AbAZKENOaUB1GJXX+ogNc1VooFm4X4ZBtGMSWMdPK7UKWhwRPvWf809M3OTooDCOFby+RCjtUypMmWK2O9zt0woBAAMEOmOhjMUpoACRpftiCvGaUNCoEMVgnnn5EaaD1JJam3tb+PYtxHudzqZUf/lfXBETAqDBznGh79ae8htPwAp7+uIeH0bS2nwl+69kGmXwJ1UoyC8mNwVtKWtnJBoIy3W+DY8YU/p0tPYxNs5WLcYbs7o0YSO7AOib7lNyQz/LNtLHgaWC48UuismEblIqJeB4rveKoQfAjvkAVZ6K9c4eXwEDAgECKXD+xDS/h+ig0ENnFALWVwo1aH54NCT5Q0jeBlX5t+iz8aL7tYsOWu0A"

func Test_Signer_Decoder(t *testing.T) {
	assert := assert.New(t)
	s := testSigner()

	t.Run("should decode tx", func(t *testing.T) {
		rawTx, _ := base64.StdEncoding.DecodeString(testMsg)
		tx, err := s.decodeTx(rawTx)
		assert.Nil(err)

		t.Run("should contain accounts", func(t *testing.T) {
			assert.Equal(4, len(tx.Message.AccountKeys))
		})

		t.Run("should contain program id", func(t *testing.T) {
			program, err := tx.ResolveProgramIDIndex(tx.Message.Instructions[0].ProgramIDIndex)
			assert.Nil(err)
			assert.Equal("uvrXzrjzPv7cAsjgnZy3hmMUnff2KiKaJpj8Ab3ffny", program.String())
		})
	})

}

func Test_Signer(t *testing.T) {
	assert := assert.New(t)
	s := testSigner()
	rawTx, _ := base64.StdEncoding.DecodeString(testMsg)

	t.Run("should fetch raw key from secret", func(t *testing.T) {
		priv, err := s.parseSignerKey()
		assert.Nil(err)
		assert.Equal("4vvm7viTdtn42UqKQuvetA4omNak5mnwXpQNqhZCdy6i", priv.PublicKey().String())

		t.Run("should build new tx with new recent blockhash", func(t *testing.T) {
			tx, err := s.decodeTx(rawTx)
			assert.Nil(err)

			sigs, err := s.sign(tx, priv)
			assert.Nil(err)
			assert.Greater(len(sigs), 0)
		})
	})

}

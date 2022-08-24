package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/avareum/avareum-hubble-signer/tests/fixtures"
	solanalib "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
)

var (
	receiverTarget = solanalib.MustPublicKeyFromBase58("BuokBaKhv1cXkU2hUSGUEKcBMyQjT4XFhQq5NkVwnwJc")
	fund           = solanalib.MustPublicKeyFromBase58("J9JqCbhirq89BQ8xBrrVAsi738Xe9RRWUR2Hys8B7fWx")
)

func NewTestTxPayload(suite *fixtures.TestSuite) []byte {
	originalTx := suite.Solana.NewTx(fund, system.NewTransferInstruction(
		1*solanalib.LAMPORTS_PER_SOL,
		fund,
		receiverTarget,
	).Build())
	payload, err := originalTx.Message.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return payload
}

func main() {
	coreService, err := solanalib.NewRandomPrivateKey()
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("SERVICE_test_service", ed25519.PrivateKey(coreService), 0644)
	if err != nil {
		panic(err)
	}

	suite := fixtures.NewTestSuite()
	transferPayload := NewTestTxPayload(suite)
	transferPayloadSignature := suite.ACL.MustSignPayloadWithKey(ed25519.PrivateKey(coreService), transferPayload)

	fmt.Println("Caller:", "test_service")
	fmt.Println("Wallet:", fund.String())
	fmt.Println("Payload (base64):", base64.StdEncoding.EncodeToString(transferPayload))
	fmt.Println("Payload Signature (base64):", base64.StdEncoding.EncodeToString(transferPayloadSignature))
}

// Caller: test_service
// Wallet: J9JqCbhirq89BQ8xBrrVAsi738Xe9RRWUR2Hys8B7fWx
// Payload (base64): AQABA/63fF+ujkuWhMEUohkc4vvXeTuoTFRbClGCxZvAz8eRohvs7+oGdi9uFWPFLFzm+bG1DIbYxC1jOPE5WmSwSr0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACDErr9EbrNsMpTwtFyAbeCzJZfuEg9x8J3I+ERtXC80AQICAAEMAgAAAADKmjsAAAAA
// Payload Signature (base64): FZEoulJ0zOa033VZyzl+3seLvIf+Bc6+pW1ICRNKlgpiLBeMf+eLvGbtshKPngFNK/Pl2ZJYQpbAmxK/UfveBQ==

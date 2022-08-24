package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/avareum/avareum-hubble-signer/tests/fixtures"
	"github.com/ethereum/go-ethereum/common"
	solanalib "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
)

var (
	solanaFund     = solanalib.MustPublicKeyFromBase58("J9JqCbhirq89BQ8xBrrVAsi738Xe9RRWUR2Hys8B7fWx")
	solanaReceiver = solanalib.MustPublicKeyFromBase58("BuokBaKhv1cXkU2hUSGUEKcBMyQjT4XFhQq5NkVwnwJc")
	bscFund        = common.HexToAddress("0xFbb2FBFacA0151148b933bE6a024cc2914032221")
	bscReceiver    = common.HexToAddress("0x9e6619A6a6cc869F384EF95f00322EE19CE12556")
)

func NewSolanaTestPayload(suite *fixtures.TestSuite) []byte {
	originalTx := suite.Solana.NewTx(solanaFund, system.NewTransferInstruction(
		1*solanalib.LAMPORTS_PER_SOL,
		solanaFund,
		solanaReceiver,
	).Build())
	payload, err := originalTx.Message.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return payload
}

func NewEvmTestPayload(suite *fixtures.TestSuite) []byte {
	originalTx := suite.Ethereum.NewTransferTransaction(bscFund, bscReceiver, 0.01)
	payload, err := originalTx.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return payload
}

func main() {
	servicePriv, err := os.ReadFile("keys/SERVICE_test_service")
	if err != nil {
		panic(err)
	}

	suite := fixtures.NewTestSuite()

	fmt.Println("Caller:", "test_service")
	fmt.Println("Service:", solanalib.PrivateKey(servicePriv).PublicKey().String())

	fmt.Println("Solana")
	solanaTransferPayload := NewSolanaTestPayload(suite)
	solanaTransferPayloadSignature := suite.ACL.MustSignPayloadWithKey(servicePriv, solanaTransferPayload)
	fmt.Println("Wallet:", solanaFund.String())
	fmt.Println("Payload (base64):", base64.StdEncoding.EncodeToString(solanaTransferPayload))
	fmt.Println("Payload signature (base64):", base64.StdEncoding.EncodeToString(solanaTransferPayloadSignature))

	fmt.Println("BSC")
	bscTransferPayload := NewEvmTestPayload(suite)
	bscTransferPayloadSignature := suite.ACL.MustSignPayloadWithKey(servicePriv, bscTransferPayload)
	fmt.Println("Wallet:", bscFund.Hex())
	fmt.Println("Payload (base64):", base64.StdEncoding.EncodeToString(bscTransferPayload))
	fmt.Println("Payload signature (base64):", base64.StdEncoding.EncodeToString(bscTransferPayloadSignature))

}

/*
Caller: test_service
Service: J4v7pQY7JcSNqfWtEx2Wo75XRHVVXFRywRb47Hs4CMfp
Solana
Wallet: J9JqCbhirq89BQ8xBrrVAsi738Xe9RRWUR2Hys8B7fWx
Payload (base64): AQABA/63fF+ujkuWhMEUohkc4vvXeTuoTFRbClGCxZvAz8eRohvs7+oGdi9uFWPFLFzm+bG1DIbYxC1jOPE5WmSwSr0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAOQysSuCqZ+b2RcWn50kD89/rNy4aakquFrj1KGn9HeQAQICAAEMAgAAAADKmjsAAAAA
Payload signature (base64): ngSCgEGIrZwYaQTNHmFLQmsrpcxSzpr3CJ8E9cBwVC+dEoxw2zy/fWqnmwOgHWZ9ez3MNbArPcH3uL1izST7Ag==
BSC
Wallet: 0xFbb2FBFacA0151148b933bE6a024cc2914032221
Payload (base64): 64CFBKgXyACCUgiUnmYZpqbMhp84TvlfADIu4ZzhJVaHI4byb8EAAICAgIA=
Payload signature (base64): f4fevGBccyGcGeXwdMhgbo6jE3Sy0tdYg/Bm/AV6Zr4U+baL1ASpFM1M7/6fS0d5iawDURv8CFb0Nq0Zg3XZDA==
*/

// Result:
// https://explorer.solana.com/tx/inspector?cluster=devnet&message=AQABA%252F63fF%252BujkuWhMEUohkc4vvXeTuoTFRbClGCxZvAz8eRohvs7%252BoGdi9uFWPFLFzm%252BbG1DIbYxC1jOPE5WmSwSr0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACDErr9EbrNsMpTwtFyAbeCzJZfuEg9x8J3I%252BERtXC80AQICAAEMAgAAAADKmjsAAAAA
// https://testnet.bscscan.com/tx/0xe7d4d0d2428271d647d8380f4ad5e94837c8c8c57d982c8a2880addec5a821cf

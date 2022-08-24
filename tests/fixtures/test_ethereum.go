package fixtures

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/avareum/avareum-hubble-signer/internal/signers/ethereum/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type EthereumTestSuite struct {
	client   *backends.SimulatedBackend
	coinbase *ecdsa.PrivateKey
}

func NewEthereumTestSuite() *EthereumTestSuite {
	e := &EthereumTestSuite{}
	e.Init()
	return e
}

func (e *EthereumTestSuite) Init() error {
	e.coinbase = types.MustNewEthereumKey()
	balance := new(big.Int)
	balance.SetString("10000000000000000000000", 10) // 10k eth in wei

	address := crypto.PubkeyToAddress(e.coinbase.PublicKey)
	genesisAlloc := map[common.Address]core.GenesisAccount{
		address: {
			Balance: balance,
		},
	}
	blockGasLimit := uint64(4712388)
	client := backends.NewSimulatedBackend(genesisAlloc, blockGasLimit)
	e.client = client
	return nil
}

func (e *EthereumTestSuite) FaucetTo(to ecdsa.PublicKey) {
	tx := e.NewTransferTransaction(crypto.PubkeyToAddress(e.coinbase.PublicKey), crypto.PubkeyToAddress(to), 10)
	signedTx, err := ethtypes.SignTx(tx, ethtypes.HomesteadSigner{}, e.coinbase)
	if err != nil {
		log.Fatal(err)
	}
	e.client.SendTransaction(context.Background(), signedTx)
	e.client.Commit()
}

func (e *EthereumTestSuite) SendTransaction(tx *ethtypes.Transaction) {
	e.client.SendTransaction(context.Background(), tx)
	e.client.Commit()
}

func (e *EthereumTestSuite) TransactionReceipt(hash common.Hash) (*ethtypes.Receipt, error) {
	return e.client.TransactionReceipt(context.Background(), hash)
}

func (e *EthereumTestSuite) NewTransferTransaction(from common.Address, to common.Address, amount float64) *ethtypes.Transaction {
	nonce, err := e.client.PendingNonceAt(context.Background(), from)
	if err != nil {
		log.Fatal(err)
	}

	weiValue := big.NewFloat(1000000000000000000)
	ethAmount := big.NewFloat(amount)
	actualAmount, _ := ethAmount.Mul(ethAmount, weiValue).Int(nil)
	tx := ethtypes.NewTx(&ethtypes.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Gas:      uint64(21000),           // The gas limit for a standard ETH transfer is 21000 units.
		GasPrice: big.NewInt(20000000000), // fix for test at 20 Gwei
		Value:    actualAmount,
	})
	return tx
}

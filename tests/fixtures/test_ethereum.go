package fixtures

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
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
	e.coinbase = e.MustNewWallet()
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
	fmt.Println(e.client.BalanceAt(context.Background(), address, nil))
	return nil
}

func (e *EthereumTestSuite) MustNewWallet() *ecdsa.PrivateKey {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	return privateKey
}

func (e *EthereumTestSuite) AirdropTo(to ecdsa.PublicKey) {
	tx := e.NewTransferTransaction(*e.coinbase, to, 10)
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, e.coinbase)
	if err != nil {
		log.Fatal(err)
	}
	e.client.SendTransaction(context.Background(), signedTx)
	e.client.Commit()
	fmt.Println(e.client.BalanceAt(context.Background(), crypto.PubkeyToAddress(to), nil))
}

func (e *EthereumTestSuite) SendTransaction(tx *types.Transaction) {
	e.client.SendTransaction(context.Background(), tx)
	e.client.Commit()
}

func (e *EthereumTestSuite) TransactionReceipt(hash common.Hash) (*types.Receipt, error) {
	return e.client.TransactionReceipt(context.Background(), hash)
}

func (e *EthereumTestSuite) NewTransferTransaction(from ecdsa.PrivateKey, to ecdsa.PublicKey, amount int64) *types.Transaction {
	toAddress := crypto.PubkeyToAddress(to)
	nonce, err := e.client.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(from.PublicKey))
	if err != nil {
		log.Fatal(err)
	}
	ethAmount := big.NewInt(amount)
	weiValue := big.NewInt(1000000000000000000)
	gasLimit := uint64(21000) // The gas limit for a standard ETH transfer is 21000 units.
	gasPrice, err := e.client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	tx := types.NewTransaction(nonce, toAddress, ethAmount.Mul(ethAmount, weiValue), gasLimit, gasPrice, nil)
	return tx
}

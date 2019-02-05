package node

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"log"
	"math/big"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EthInstance struct that holds the connection to the node
type EthInstance struct {
	client     *ethclient.Client
	privateKey *ecdsa.PrivateKey
}

// New instance of the client
func New(url, privateKeyHex string) (*EthInstance, error) {
	instance := new(EthInstance)

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	instance.privateKey = privateKey

	client, err := ethclient.Dial(url)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	instance.client = client

	return instance, err
}

// IsHexAddress check if string is address
func IsHexAddress(s string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(s)
}

// Faucet to transfer eth from a fixed account to an address
func (e *EthInstance) Faucet(toAddrStr string, wei *big.Int) (string, error) {
	publicKey := e.privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		err := errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		log.Println(err)
		return "", err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := e.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println(err)
		return "", err
	}

	value := wei
	gasLimit := uint64(21000)
	gasPrice, err := e.client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err)
		return "", err
	}
	toAddress := common.HexToAddress(toAddrStr)

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	chainID, err := e.client.NetworkID(context.Background())
	if err != nil {
		log.Println(err)
		return "", err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), e.privateKey)
	if err != nil {
		log.Println(err)
		return "", err
	}

	err = e.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return signedTx.Hash().String(), nil
}

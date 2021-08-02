package main

import (
	"context"
	"discrete-systems/monorepo/pkg/eth"
	"log"
	"math/big"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	providerUrl := "https://eth-goerli.alchemyapi.io/v2/_gg7wSSi0KMBsdKnGVfHDueq6xMB9EkC"
	client, _ := ethclient.Dial(providerUrl)

	var gasTipCap *big.Int

	log.Println(os.Args[2])
	if os.Args[2] != "" {
		v, _ := strconv.Atoi(os.Args[2])
		gasTipCap = big.NewInt(int64(v))
	} else {
		gasTipCap, _ = client.SuggestGasTipCap(context.Background())
	}
	log.Println("gasTipCap: ", gasTipCap)

	privKeyHex := os.Args[1]

	privateKey, publicKeyECDSA, _ := eth.PrivKeyHexToEcdsaKeypair(privKeyHex)
	fromAddress := eth.PublicKeyEcdsaToAddress(*publicKeyECDSA)

	nonce, _ := eth.GetNonceForTx(client, fromAddress.Hex())

	gasLimit := uint64(21000)
	chainID := int64(5)

	toAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")
	value := *big.NewInt(0)

	dynamicFeeTx := eth.NewUnsignedDynamicFeeTx(
		chainID,
		nonce,
		*gasTipCap,
		*gasTipCap,
		gasLimit,
		toAddress.Hex(),
		value,
		nil,
		types.AccessList{},
	)

	tx := types.NewTx(&dynamicFeeTx)
	signer := types.NewLondonSigner(big.NewInt(chainID))
	signature, _ := crypto.Sign(signer.Hash(tx).Bytes(), privateKey)
	signedTx, _ := tx.WithSignature(signer, signature)

	err := client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal("err:", err)
	}

	txHash := signedTx.Hash().Hex()

	log.Println("https://goerli.etherscan.io/tx/" + txHash)
}

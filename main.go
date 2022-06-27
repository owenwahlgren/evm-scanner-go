package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	WssNodeEndpoint  = "ws://localhost:8545"
	HttpNodeEndpoint = "http://localhost:3000"
)

func main() {
	ctx := context.Background()
	txnsHash := make(chan common.Hash)

	baseClient, err := rpc.Dial(WssNodeEndpoint)
	if err != nil {
		log.Fatalln(err)
	}

	ethClient, err := ethclient.Dial(WssNodeEndpoint)
	if err != nil {
		log.Fatalln(err)
	}

	chainID, err := ethClient.NetworkID(ctx)
	if err != nil {
		log.Fatal(err)
	}

	subscriber := gethclient.New(baseClient)
	_, err = subscriber.SubscribePendingTransactions(ctx, txnsHash)

	if err != nil {
		log.Fatalln(err)
	}

	signer := types.NewLondonSigner(chainID)
	defer func() {
		fmt.Println("error")
		baseClient.Close()
		ethClient.Close()
	}()

	for txnHash := range txnsHash {
		fmt.Println("tx by hash:", txnHash)
		txn, _, err := ethClient.TransactionByHash(ctx, txnHash)
		if err != nil {
			fmt.Println("error here!")
		}

		message, err := txn.AsMessage(signer, nil)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(txnHash.String())
		fmt.Println(message.To())
	}
}

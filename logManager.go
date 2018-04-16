package main

import (
	"github.com/ethereum/go-ethereum"
	"os"
	"fmt"
	"log"
	"math/big"
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"encoding/json"
)

var fltr ethereum.FilterQuery

func createFile(address string) {
	fmt.Print("Creating file for account: ", address, "\n")

	cl, err := ethclient.Dial(*server)
	if err != nil {
		log.Panic("Connection Error: ", err)
	}

	f, err := os.Create(fmt.Sprint("logs/", address, ".json"))
	if err != nil {
		log.Panic("File Creation Error: ", err)
	}
	fltr.Addresses = []common.Address{common.HexToAddress(address)}
	fltr.FromBlock = big.NewInt(int64(*fromBlock))
	ctx := context.Background()
	lgs, err := cl.FilterLogs(ctx, fltr)
	if err != nil {
		log.Panic("Filter Error: ", err)
	}
	b, err := json.Marshal(lgs)
	if err != nil {
		log.Panic("Marshal Error: ", err)
	}
	_, err = f.Write(b)
	if err != nil {
		log.Panic("File Write Error: ", err)
	}
	f.Close()
	fmt.Print("File Created\n")
}
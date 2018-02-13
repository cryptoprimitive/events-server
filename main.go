package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event/filter"
	"log"
)

func main() {
	server := "http://localhost:8545"
	//server := "https://mainnet.infura.io"

	cl, err := ethclient.Dial(server)
	if err != nil {
		log.Panic("Connection Error: ", err)
	}

	var filter ethereum.FilterQuery
	//0x4df81d58993ff6f6e3a721b2ac0a08a5cd78ce9e

	filter.Addresses = []common.Address{common.HexToAddress("0x4df81d58993ff6f6e3a721b2ac0a08a5cd78ce9e")}
	ctx := context.Background()
	lgs, err := cl.FilterLogs(ctx, filter)
	if err != nil {
		log.Panic("Filter Error: ", err)
	}
	fmt.Print(lgs)

	//An Early block transaction to test a geth node that isn't fully sync'd.
	//ctx := context.Background()
	//tx, pending, _ := cl.TransactionByHash(ctx, common.HexToHash("0x1b728581a737edb547b39381e13969e191c7263030ba966291fcb707b9440c87"))
	//if !pending {
	//	fmt.Println(tx)
	//}
}

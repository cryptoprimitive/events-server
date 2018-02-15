package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	server := "http://localhost:8545"
	//server := "https://mainnet.infura.io"

	cl, err := ethclient.Dial(server)
	if err != nil {
		log.Panic("Connection Error: ", err)
	}

	var fltr ethereum.FilterQuery
	//0x4df81d58993ff6f6e3a721b2ac0a08a5cd78ce9e
	//0x6090A6e47849629b7245Dfa1Ca21D94cd15878Ef

	fltr.Addresses = []common.Address{common.HexToAddress(r.URL.Path[1:])}
	ctx := context.Background()
	lgs, err := cl.FilterLogs(ctx, fltr)
	if err != nil {
		log.Panic("Filter Error: ", err)
	}


	//An Early block transaction to test a geth node that isn't fully sync'd.
	//ctx := context.Background()
	//tx, pending, _ := cl.TransactionByHash(ctx, common.HexToHash("0x14fcfe755cf24fe8d36464d58e0e333d2a0a59fdad07193ed0fcf6c34557772c"))
	//if !pending {
	//	fmt.Println(tx)
	//}

	// Encode Response to the writer

	for _, l := range lgs {
		//b, err := l.MarshalJSON()
		//if err != nil {
		//	log.Panic("JSON Marshal Error: ", err)
		//}
		fmt.Fprint(w, l)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

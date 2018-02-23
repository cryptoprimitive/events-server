package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"net/http"
	"flag"
)

var serverMode = flag.String("mode","production","Set to 'testing' to enable address and tx lookup.")
var server = "http://localhost:8545"
//var server = "https://mainnet.infura.io"

func addressHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Address Balance requested: ", r.URL.Path[6:], "\n")

	cl, err := ethclient.Dial(server)
	if err != nil {
		log.Panic("Connection Error: ", err)
	}

	ctx := context.Background()
	balance, err := cl.BalanceAt(ctx, common.StringToAddress(r.URL.Path[6:]), nil)
	fmt.Fprint(w, balance)
}

func txHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Transaction requested: ", r.URL.Path[4:], "\n")

	cl, err := ethclient.Dial(server)
	if err != nil {
		log.Panic("Connection Error: ", err)
	}

	ctx := context.Background()
	tx, pending, _ := cl.TransactionByHash(ctx, common.HexToHash(r.URL.Path[4:]))
	if !pending {
		fmt.Fprint(w, tx)
	} else {
		fmt.Fprint(w, "Warning: Transaction Pending\n")
	}
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	if *serverMode == "testing" {
		fmt.Fprint(w, "Events Requested: ", r.URL.Path[8:], "\n")
	}

	cl, err := ethclient.Dial(server)
	if err != nil {
		log.Panic("Connection Error: ", err)
	}

	var fltr ethereum.FilterQuery
	//0x4df81d58993ff6f6e3a721b2ac0a08a5cd78ce9e
	//0x6090A6e47849629b7245Dfa1Ca21D94cd15878Ef

	fltr.Addresses = []common.Address{common.HexToAddress(r.URL.Path[8:])}
	ctx := context.Background()
	lgs, err := cl.FilterLogs(ctx, fltr)
	if err != nil {
		log.Panic("Filter Error: ", err)
	}


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
	flag.Parse()
	if *serverMode == "testing" {
		http.HandleFunc("/addr/", addressHandler)
		http.HandleFunc("/tx/",txHandler)
	}
	http.HandleFunc("/events/", eventsHandler)
	http.ListenAndServe(":8080", nil)
}

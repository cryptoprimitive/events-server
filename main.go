package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"net/http"
	"flag"
	"math/big"
	"strconv"
	"github.com/ethereum/go-ethereum/core/types"
	"encoding/json"
)

var serverMode = flag.String("serverMode","production","Set to 'testing' to enable debug access.")
var server = flag.String("host", "http://localhost:8545", "Set server host.")
//var server = "https://mainnet.infura.io"
var fromBlock = flag.Int("fromBlock", 0, "Set to block to start server queries from.")

func syncHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Node sync status requested\n")

	cl, err := ethclient.Dial(*server)
	if err != nil {
		log.Panic("Connection Error: ", err)
	}

	ctx := context.Background()
	prog, err := cl.SyncProgress(ctx)
	if err != nil {
		log.Panic("Error Fetching Sync Status: ", err)
	}

	if prog == nil {
		fmt.Fprint(w, "Syncing complete!\n")
	} else {
		fmt.Fprint(w, "Current Block: ", prog.CurrentBlock, "\nHighestBlock: ", prog.HighestBlock)
	}
}

func addressHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Address Balance Requested: ", r.URL.Path[6:], "\n")

	cl, err := ethclient.Dial(*server)
	if err != nil {
		log.Panic("Connection Error: ", err)
	}

	ctx := context.Background()
	balance, err := cl.BalanceAt(ctx, common.HexToAddress(r.URL.Path[6:]), nil)
	if err != nil {
		log.Panic("Error Fetching Balance: ", err)
	}

	fmt.Fprint(w, balance)
}

func txHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Transaction Requested: ", r.URL.Path[4:], "\n")

	cl, err := ethclient.Dial(*server)
	if err != nil {
		log.Panic("Connection Error: ", err)
	}

	ctx := context.Background()
	tx, pending, err := cl.TransactionByHash(ctx, common.HexToHash(r.URL.Path[4:]))
	if err != nil {
		log.Panic("Error Fetching Transaction: ", err)
	}

	if !pending {
		fmt.Fprint(w, tx)
	} else {
		fmt.Fprint(w, "Warning: Transaction Pending\n")
		fmt.Fprint(w, tx)
	}
}

func blockHandler(w http.ResponseWriter, r *http.Request) {
	cl, err := ethclient.Dial(*server)
	if err != nil {
		log.Panic("Connection Error: ", err)
	}

	ctx := context.Background()
	i, err := strconv.Atoi(r.URL.Path[7:])
	if err != nil {
		log.Panic("Error Coverting Block Number: ", err)
	}
	block, err := cl.BlockByNumber(ctx, big.NewInt(int64(i)))
	if err != nil {
		log.Panic("Block Fetch Error: ", err)
	}
	fmt.Fprint(w, block)
}

func blockeventsHandler(w http.ResponseWriter, r *http.Request) {
	if *serverMode == "testing" {
		fmt.Fprint(w, "Block Events Requested: ", r.URL.Path[13:], "\n")
	}
	cl, err := ethclient.Dial(*server)
	if err != nil {
		log.Panic("Connection Error: ", err)
	}

	ctx := context.Background()

	i, err := strconv.Atoi(r.URL.Path[13:])
	if err != nil {
		log.Panic("Error Coverting Block Number: ", err)
	}

	block, err := cl.BlockByNumber(ctx, big.NewInt(int64(i)))
	if err != nil {
		log.Panic("Block Fetch Error: ", err)
	}

	txs := block.Transactions()

	for _, t := range txs {
		receipt, err := cl.TransactionReceipt(ctx, t.Hash())
		if err != nil {
			log.Panic("Receipt Error: ", err)
		}
		for _, lg := range receipt.Logs {
			b, err := lg.MarshalJSON()
			if err != nil {
				log.Panic("JSON Marshalling Error: ", err)
			}
			_, err = w.Write(b)
			if err != nil {
				log.Panic("Error Writing Output: ", err)
			}
		}
	}
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	addr := r.URL.Path[8:]
	if *serverMode == "testing" {
		fmt.Fprint(w, "Events Requested: ", addr, "\n")
		fmt.Fprint(w, "Starting from Block: ", *fromBlock, "\n")
	}
	evtReturner := eventReturner{address: addr}
	eventReturnerChan <- &evtReturner

	v, err := json.Marshal(evtReturner.logs)
	if err != nil {
		log.Panic("Marshal Error: ", err)
	}
	_, err = w.Write(v)
	if err != nil {
		log.Panic("Error Writing Logs")
	}
	//0x4df81d58993ff6f6e3a721b2ac0a08a5cd78ce9e
	//0x6090A6e47849629b7245Dfa1Ca21D94cd15878Ef
}

var eventReturnerChan = make(chan *eventReturner)
var newEventsChan = make(chan *types.Log)

func main() {
	flag.Parse()
	if *serverMode == "testing" {
		http.HandleFunc("/addr/", addressHandler)
		http.HandleFunc("/tx/", txHandler)
		http.HandleFunc("/block/", blockHandler)
		http.HandleFunc("/sync", syncHandler)
	}

	go FileManager(eventReturnerChan, newEventsChan)

	http.HandleFunc("/blockevents/", blockeventsHandler)
	http.HandleFunc("/events/", eventsHandler)
	http.ListenAndServe(":8080", nil)
}

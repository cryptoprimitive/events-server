package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"os"
)

var fltr ethereum.FilterQuery

const (
	EVENTS_PENDING = 0
	EVENTS_SYNCED  = 1
)

type addressFile struct {
	fileName string
	status   int
	err      error
}

func logFilename(address string) string {
	return fmt.Sprint("logs/", address, ".json")
}

func createFile(address string) (addrFile *addressFile) {
	fmt.Print("Creating file for account: ", address, "\n")

	cl, err := ethclient.Dial(*server)
	if err != nil {
		addrFile.err = err
		return addrFile
		//log.Panic("Connection Error: ", err)
	}
	filename := logFilename(address)
	f, err := os.Create(filename)
	if err != nil {
		addrFile.err = err
		return addrFile
		//log.Panic("File Creation Error: ", err)
	}
	fltr.Addresses = []common.Address{common.HexToAddress(address)}
	fltr.FromBlock = big.NewInt(int64(*fromBlock))
	ctx := context.Background()
	lgs, err := cl.FilterLogs(ctx, fltr)
	if err != nil {
		addrFile.err = err
		return addrFile
		//log.Panic("Filter Error: ", err)
	}
	b, err := json.Marshal(lgs)
	if err != nil {
		addrFile.err = err
		return addrFile
		//log.Panic("Marshal Error: ", err)
	}
	_, err = f.Write(b)
	if err != nil {
		addrFile.err = err
		return addrFile
		//log.Panic("File Write Error: ", err)
	}
	addrFile.fileName = filename
	f.Close()
	addrFile.status = EVENTS_SYNCED
	fmt.Print("File Created\n")
	return addrFile
}

type eventReturner struct {
	address string
	logs    []types.Log
	err     error
}

func FileManager(eventsReturnerChan chan *eventReturner, newEventsChan chan *types.Log) {
	dir, err := os.Open("logs")
	if err != nil {
		log.Panic("Error Opening Directory: ", err)
	}
	files, err := dir.Readdirnames(-1)
	var AllAccounts map[string]addressFile
	for _, fileString := range files {
		AllAccounts[fileString] = addressFile{fileName: fileString, status: isSynced(fileString)}
	}
	for {
		select {
		//case newEvent := <-newEventsChan:
		//handle new event
		case returner := <-eventsReturnerChan:
			filename := logFilename(returner.address)
			account, exists := AllAccounts[filename]
			if exists == false {
				createFile(returner.address)
				returner.err = fmt.Errorf("Account Being Created")
			}
			if account.status == EVENTS_PENDING {
				returner.err = fmt.Errorf("Events Being Synced")
				break
			}

			f, err := os.Open(account.fileName)
			if err != nil {
				log.Panic("File error: ", err)
			}
			stat, err := f.Stat()
			if err != nil {
				log.Panic("Stat Error: ", err)
			}

			b := make([]byte, stat.Size())
			_, err = f.Read(b)
			if err != nil {
				log.Panic("File Read Error: ", err)
			}

			lgs := make([]types.Log, 0)
			err = json.Unmarshal(b, &lgs)
			if err != nil {
				log.Panic("Unmarshal Error: ", err)
			}
			returner.logs = lgs
			//Return new events
		}
	}
}

func isSynced(fileString string) int {
	return EVENTS_SYNCED
}

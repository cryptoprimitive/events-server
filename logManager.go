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
	return fmt.Sprint(address, ".json")
}

func createFile(address string) (addrFile *addressFile) {
	fmt.Print("Creating file for account: ", address, "\n")

	addrFile = new(addressFile)

	cl, err := ethclient.Dial(*server)
	if err != nil {
		fmt.Print(err)
		addrFile.err = err
		return addrFile
		//log.Panic("Connection Error: ", err)
	}
	filename := logFilename(address)
	f, err := os.Create(fmt.Sprint("logs/", filename))
	if err != nil {
		fmt.Print(err)
		addrFile.err = err
		return addrFile
		//log.Panic("File Creation Error: ", err)
	}
	fltr.Addresses = []common.Address{common.HexToAddress(address)}
	fltr.FromBlock = big.NewInt(int64(*fromBlock))
	ctx := context.Background()
	lgs, err := cl.FilterLogs(ctx, fltr)
	if err != nil {
		fmt.Print(err)
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
	fmt.Print("File Created\n", addrFile, "\n")
	return addrFile
}

type eventReturner struct {
	address string
	logs    []types.Log
	err     error
}

type createAccount struct {
	filename string
	address string
	addrFile addressFile
}

func createAccountFunc(createAccountStruct createAccount, accountCreatedChan chan *createAccount) {
	createAccountStruct.addrFile = *createFile(createAccountStruct.address)
	accountCreatedChan <- &createAccountStruct
}

func FileManager(eventsReturnerChan chan *eventReturner, newEventsChan chan *types.Log) {
	var accountCreatedChan = make(chan *createAccount)

	dir, err := os.Open("logs")
	if err != nil {
		log.Panic("Error Opening Directory: ", err)
	}
	files, err := dir.Readdirnames(-1)
	var AllAccounts = make(map[string]addressFile)
	for _, fileString := range files {
		AllAccounts[fileString] = addressFile{fileName: fileString, status: isSynced(fileString)}
	}
	//fmt.Print(AllAccounts)
	for {
		select {
		case newEvent := <-newEventsChan:
		//handle new event
			address := newEvent.Address.Hex()
			fmt.Print(address, "\n", AllAccounts, "\n")
			filename := logFilename(address)
			account, exists := AllAccounts[filename]
			fmt.Print(account, "\n", exists, "\n", filename, "\n", AllAccounts)
			if exists == true {
				f, err := os.OpenFile(fmt.Sprint("logs/", account.fileName), os.O_RDWR, 0644)
				if err != nil {
					log.Panic("File error: ", err)
				}
				_, err = f.Seek(-1,2)
				if err != nil {
					log.Panic("Error seeking for file write: ", err)
				}
				b, err := json.Marshal(newEvent)
				if err != nil {
					log.Panic("Marshal Error: ", err)
				}
				_, err = f.Write([]byte{','})
				if err != nil {
					fmt.Print("Error Writing Literal: ", err)
				}
				fmt.Print("Writing: ", b[1:])
				_, err = f.Write(b)
				if err != nil {
					log.Panic("Error Writing: ", err)
				}
				_, err = f.Write([]byte{']'})
				if err != nil {
					log.Panic("Error Writing: ", err)
				}
				f.Close()

				//add event to file
			}
		case returner := <-eventsReturnerChan:
			filename := logFilename(returner.address)
			account, exists := AllAccounts[filename]
			//fmt.Print(account, exists, filename, "\n", AllAccounts)
			if exists == false {
				addrfile := new(addressFile)
				addrfile.fileName = filename
				addrfile.status = EVENTS_PENDING
				AllAccounts[filename] = *addrfile
				createAccountStruct := createAccount{filename: filename, address: returner.address, addrFile: *addrfile}
				go createAccountFunc(createAccountStruct, accountCreatedChan)
				returner.err = fmt.Errorf("Account Being Create %s", returner.address)
				eventReturnerChan <- returner
				break
			}
			if account.status == EVENTS_PENDING {
				returner.err = fmt.Errorf("Events Being Synced for account %s", returner.address)
				eventReturnerChan <- returner
				break
			}

			f, err := os.Open(fmt.Sprint("logs/", account.fileName))
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
			eventReturnerChan <- returner
		case accountCreatedStruct := <-accountCreatedChan:
			AllAccounts[accountCreatedStruct.filename] = accountCreatedStruct.addrFile
		}
	}
}

func SubListener(newEventsChan chan *types.Log) {
	ctx := context.Background()
	headerChan := make(chan *types.Header)
	cl, err := ethclient.Dial(*server)
	if err != nil {
		log.Panic("Connection Error: ", err)
	}
	subscription, err := cl.SubscribeNewHead(ctx, headerChan)
	if err != nil {
		fmt.Print("Header Subscription Fail: ", err,"\nEvents not updating")
		return
	}
	defer subscription.Unsubscribe()

	for {
		header := <-headerChan
		fmt.Print(header, "\n")
		num := header.Number
		block, err := cl.BlockByNumber(ctx, num)
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
				newEventsChan <- lg
			}
		}
	}
}

func testListener(newEventsChan chan *types.Log) {
	f, err := os.Open("testdata/test.json")
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
	for _, lg := range lgs {
		//fmt.Print("Sending along channel:", lg,"\n")
		newEventsChan <- &lg
	}
}

func isSynced(fileString string) int {
	return EVENTS_SYNCED
}

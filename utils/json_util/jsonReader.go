package main

import (
	"os"
	"log"
	//"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"encoding/json"
	//"io"
)

func main() {
	//f, err := os.Open("events/test.json")
	f, err := os.Open("events/testdata.json")
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
	//var lg types.Log
	//err = json.Unmarshal(b, &lg)
	//if err != nil {
	//	log.Panic("Unmarshal Error: ", err)
	//}
	//lgs := make([]types.Log, 1)
	//lgs[0] = lg

	//fmt.Print(lg, "\n")

	lgs := make([]types.Log, 0)
	err = json.Unmarshal(b, &lgs)
	if err != nil {
		log.Panic("Unmarshal Error: ", err)
	}
	v, err := json.Marshal(lgs)
	if err != nil {
		log.Panic("Marshal Error: ", err)
	}
	_, _ = os.Stdout.Write(v)
	//fmt.Print(lgs, "\n")
}
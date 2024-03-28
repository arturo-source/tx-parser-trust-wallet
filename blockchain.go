package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Blockchain struct {
	LastBlockNum int
	Suscribers   map[string][]Transaction
}

func newBlockchain() Parser {
	return &Blockchain{
		LastBlockNum: 0,
		Suscribers:   make(map[string][]Transaction),
	}
}

func hexToDec(hex string) int {
	numStr := strings.TrimPrefix(hex, "0x")
	num, err := strconv.ParseInt(numStr, 16, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing number: %s\n", err)
		return 0
	}

	return int(num)
}

func (b *Blockchain) GetCurrentBlock() int {
	data := RPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []any{"latest", false},
		ID:      1,
	}

	respData, err := doRequest(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting block: %s\n", err)
		return 0
	}

	num := hexToDec(respData.Result.Number)
	b.LastBlockNum = num
	return num
}

func (b *Blockchain) Subscribe(address string) bool {
	if _, ok := b.Suscribers[address]; ok {
		return false
	}

	b.Suscribers[address] = []Transaction{}
	return true
}

func (b *Blockchain) GetTransactions(address string) []Transaction {
	if _, ok := b.Suscribers[address]; !ok {
		fmt.Fprintf(os.Stderr, "Address %s not subscribed\n", address)
	}

	return b.Suscribers[address]
}

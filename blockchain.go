package main

import (
	"fmt"
	"os"
	"time"
)

type Blockchain struct {
	LastBlockNum int
	Subscribers  map[string][]Transaction
}

func newBlockchain() Parser {
	b := &Blockchain{
		LastBlockNum: 0,
		Subscribers:  make(map[string][]Transaction),
	}

	go b.backgroundListening()
	return b
}

func (b *Blockchain) backgroundListening() {
	ticker := time.NewTicker(time.Second)

	for range ticker.C {
		b.readBlocks()
	}
}

func (b *Blockchain) readBlocks() {
	newLastBlockNum, err := getLatestBlock()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting latest block: %s\n", err)
		return
	}

	for blockNum := b.LastBlockNum + 1; blockNum <= newLastBlockNum; blockNum++ {
		block, err := getBlockByNumber(blockNum)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting block: %s\n", err)
			return
		}

		b.addTransactionsFromBlock(block)
	}

	b.LastBlockNum = newLastBlockNum
}

func (b *Blockchain) addTransactionsFromBlock(block Block) {
	for _, tx := range block.Result.Transactions {
		for sub := range b.Subscribers {
			switch sub {
			case tx.From, tx.To:
				b.Subscribers[sub] = append(b.Subscribers[sub], tx)
			}
		}
	}
}

func (b *Blockchain) GetCurrentBlock() int {
	return b.LastBlockNum
}

func (b *Blockchain) Subscribe(address string) bool {
	if _, ok := b.Subscribers[address]; ok {
		return false
	}

	b.Subscribers[address] = []Transaction{}
	return true
}

func (b *Blockchain) GetTransactions(address string) []Transaction {
	if _, ok := b.Subscribers[address]; !ok {
		fmt.Fprintf(os.Stderr, "Address %s not subscribed\n", address)
	}

	return b.Subscribers[address]
}

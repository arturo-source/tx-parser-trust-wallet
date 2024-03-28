package main

import (
	"fmt"
	"os"
	"time"
)

type Blockchain struct {
	LastBlockNum int
	Database     *Memory
}

func newBlockchain() (Parser, error) {
	b := &Blockchain{}

	lastBlockNum, err := getLatestBlock()
	if err != nil {
		return b, fmt.Errorf("error getting latest block: %s", err)
	}

	b.LastBlockNum = lastBlockNum
	b.Database = NewMemory()
	go b.backgroundListening()

	return b, nil
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
		for _, subAddr := range b.Database.GetAllSubscribers() {
			switch subAddr {
			case tx.From, tx.To:
				b.Database.AddTransaction(subAddr, tx)
			}
		}
	}
}

func (b *Blockchain) GetCurrentBlock() int {
	return b.LastBlockNum
}

func (b *Blockchain) Subscribe(address string) bool {
	if b.Database.SubscriberExist(address) {
		return false
	}

	b.Database.AddSubscriber(address)
	return true
}

func (b *Blockchain) GetTransactions(address string) []Transaction {
	if !b.Database.SubscriberExist(address) {
		fmt.Fprintf(os.Stderr, "Address %s not subscribed\n", address)
	}

	return b.Database.GetSubscriberTransactions(address)
}

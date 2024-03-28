package main

import (
	"fmt"
	"os"
	"time"
)

type Blockchain struct {
	LastBlockNum int
	Subscribers  Memory
}

func newBlockchain() (Parser, error) {
	b := &Blockchain{}

	lastBlockNum, err := getLatestBlock()
	if err != nil {
		return b, fmt.Errorf("error getting latest block: %s", err)
	}

	b.LastBlockNum = lastBlockNum
	b.Subscribers = NewMemory()
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
		for _, subAddr := range b.Subscribers.GetAllSubscribers() {
			switch subAddr {
			case tx.From, tx.To:
				b.Subscribers.AddTransaction(subAddr, tx)
			}
		}
	}
}

func (b *Blockchain) GetCurrentBlock() int {
	return b.LastBlockNum
}

func (b *Blockchain) Subscribe(address string) bool {
	if b.Subscribers.SubscriberExist(address) {
		return false
	}

	b.Subscribers.AddSubscriber(address)
	return true
}

func (b *Blockchain) GetTransactions(address string) []Transaction {
	if !b.Subscribers.SubscriberExist(address) {
		fmt.Fprintf(os.Stderr, "Address %s not subscribed\n", address)
	}

	return b.Subscribers.GetSubscriberTransactions(address)
}

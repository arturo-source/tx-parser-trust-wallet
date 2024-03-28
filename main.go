package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func run() error {
	port := flag.String("port", "8080", "port to listen on")
	flag.Parse()

	blockchain, err := newBlockchain()
	if err != nil {
		return fmt.Errorf("error creating blockchain: %s", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/getCurrentBlock", func(w http.ResponseWriter, r *http.Request) {
		blockNum := blockchain.GetCurrentBlock()
		w.Write([]byte(strconv.Itoa(blockNum)))
	})

	mux.HandleFunc("/suscribe", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		alreadySubscribed := blockchain.Subscribe(address)
		if alreadySubscribed {
			w.Write([]byte("Already subscribed"))
		} else {
			w.Write([]byte("Subscribed"))
		}
	})

	mux.HandleFunc("/getTransactions", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		transactions := blockchain.GetTransactions(address)
		w.Write([]byte(fmt.Sprintf("%v", transactions)))
	})

	if err := http.ListenAndServe(":"+*port, mux); err != nil {
		return fmt.Errorf("error starting server: %s", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

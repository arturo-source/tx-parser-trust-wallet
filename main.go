package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func main() {
	port := flag.String("port", "8080", "port to listen on")
	flag.Parse()

	blockchain, err := newBlockchain()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating blockchain: %s\n", err)
		os.Exit(1)
	}

	http.HandleFunc("/getCurrentBlock", func(w http.ResponseWriter, r *http.Request) {
		blockNum := blockchain.GetCurrentBlock()
		w.Write([]byte(strconv.Itoa(blockNum)))
	})

	http.HandleFunc("/suscribe", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		alreadySubscribed := blockchain.Subscribe(address)
		if alreadySubscribed {
			w.Write([]byte("Already subscribed"))
		} else {
			w.Write([]byte("Subscribed"))
		}
	})

	http.HandleFunc("/getTransactions", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		transactions := blockchain.GetTransactions(address)
		w.Write([]byte(fmt.Sprintf("%v", transactions)))
	})

	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %s\n", err)
		os.Exit(1)
	}
}

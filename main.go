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

	blockchain := Blockchain{}

	http.HandleFunc("/getCurrentBlock", func(w http.ResponseWriter, r *http.Request) {
		blockNum := blockchain.GetCurrentBlock()
		w.Write([]byte(strconv.Itoa(blockNum)))
	})

	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %s\n", err)
		os.Exit(1)
	}
}

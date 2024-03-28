package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const RPC_URL = "https://cloudflare-eth.com"

type Blockchain struct{}

type RPCRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      int    `json:"id"`
}

type RPCResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		BaseFeePerGas         string   `json:"baseFeePerGas"`
		BlobGasUsed           string   `json:"blobGasUsed"`
		Difficulty            string   `json:"difficulty"`
		ExcessBlobGas         string   `json:"excessBlobGas"`
		ExtraData             string   `json:"extraData"`
		GasLimit              string   `json:"gasLimit"`
		GasUsed               string   `json:"gasUsed"`
		Hash                  string   `json:"hash"`
		LogsBloom             string   `json:"logsBloom"`
		Miner                 string   `json:"miner"`
		MixHash               string   `json:"mixHash"`
		Nonce                 string   `json:"nonce"`
		Number                string   `json:"number"`
		ParentBeaconBlockRoot string   `json:"parentBeaconBlockRoot"`
		ParentHash            string   `json:"parentHash"`
		ReceiptsRoot          string   `json:"receiptsRoot"`
		Sha3Uncles            string   `json:"sha3Uncles"`
		Size                  string   `json:"size"`
		StateRoot             string   `json:"stateRoot"`
		Timestamp             string   `json:"timestamp"`
		TotalDifficulty       string   `json:"totalDifficulty"`
		Transactions          []string `json:"transactions"`
		TransactionsRoot      string   `json:"transactionsRoot"`
		Uncles                []any    `json:"uncles"`
		Withdrawals           []struct {
			Index          string `json:"index"`
			ValidatorIndex string `json:"validatorIndex"`
			Address        string `json:"address"`
			Amount         string `json:"amount"`
		} `json:"withdrawals"`
		WithdrawalsRoot string `json:"withdrawalsRoot"`
	} `json:"result"`
	ID int `json:"id"`
}

func (b Blockchain) GetCurrentBlock() int {
	data := RPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []any{"latest", false},
		ID:      1,
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding data: %s\n", err)
		return 0
	}

	req, err := http.NewRequest("POST", RPC_URL, buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating request: %s\n", err)
		return 0
	}

	req.Header.Set("Content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error doing request: %s\n", err)
		return 0
	}
	defer resp.Body.Close()

	var respData RPCResponse
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding response: %s\n", err)
		return 0
	}

	numStr := strings.TrimPrefix(respData.Result.Number, "0x")
	num, err := strconv.ParseInt(numStr, 16, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing number: %s\n", err)
		return 0
	}

	return int(num)
}

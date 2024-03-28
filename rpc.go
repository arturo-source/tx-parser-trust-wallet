package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const RPC_URL = "https://cloudflare-eth.com"

type RPCRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      int    `json:"id"`
}

type Block struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		BaseFeePerGas         string        `json:"baseFeePerGas"`
		BlobGasUsed           string        `json:"blobGasUsed"`
		Difficulty            string        `json:"difficulty"`
		ExcessBlobGas         string        `json:"excessBlobGas"`
		ExtraData             string        `json:"extraData"`
		GasLimit              string        `json:"gasLimit"`
		GasUsed               string        `json:"gasUsed"`
		Hash                  string        `json:"hash"`
		LogsBloom             string        `json:"logsBloom"`
		Miner                 string        `json:"miner"`
		MixHash               string        `json:"mixHash"`
		Nonce                 string        `json:"nonce"`
		Number                string        `json:"number"`
		ParentBeaconBlockRoot string        `json:"parentBeaconBlockRoot"`
		ParentHash            string        `json:"parentHash"`
		ReceiptsRoot          string        `json:"receiptsRoot"`
		Sha3Uncles            string        `json:"sha3Uncles"`
		Size                  string        `json:"size"`
		StateRoot             string        `json:"stateRoot"`
		Timestamp             string        `json:"timestamp"`
		TotalDifficulty       string        `json:"totalDifficulty"`
		Transactions          []Transaction `json:"transactions"`
		TransactionsRoot      string        `json:"transactionsRoot"`
		Uncles                []any         `json:"uncles"`
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

func doRequest(data RPCRequest) (Block, error) {
	var block Block

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(data)
	if err != nil {
		return block, fmt.Errorf("error encoding data: %s", err)
	}

	req, err := http.NewRequest("POST", RPC_URL, buf)
	if err != nil {
		return block, fmt.Errorf("error creating request: %s", err)
	}

	req.Header.Set("Content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return block, fmt.Errorf("error doing request: %s", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&block)
	if err != nil {
		return block, fmt.Errorf("error decoding response: %s", err)
	}
	return block, nil
}

func hexToDec(hex string) (int, error) {
	numStr := strings.TrimPrefix(hex, "0x")
	num, err := strconv.ParseInt(numStr, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing number: %s", err)
	}

	return int(num), nil
}

func getLatestBlock() (int, error) {
	data := RPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []any{"latest", true},
		ID:      1,
	}

	block, err := doRequest(data)
	if err != nil {
		return 0, fmt.Errorf("error getting block: %s", err)
	}

	num, err := hexToDec(block.Result.Number)
	return num, err
}

func decToHex(dec int) string {
	return fmt.Sprintf("0x%x", dec)
}

func getBlockByNumber(number int) (Block, error) {
	numberHex := decToHex(number)
	data := RPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []any{numberHex, true},
		ID:      1,
	}

	block, err := doRequest(data)
	if err != nil {
		return block, fmt.Errorf("error getting block: %s", err)
	}

	return block, nil
}

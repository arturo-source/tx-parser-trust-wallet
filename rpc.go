package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const RPC_URL = "https://cloudflare-eth.com"

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

func doRequest(data RPCRequest) (RPCResponse, error) {
	var respData RPCResponse

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(data)
	if err != nil {
		return respData, fmt.Errorf("error encoding data: %s", err)
	}

	req, err := http.NewRequest("POST", RPC_URL, buf)
	if err != nil {
		return respData, fmt.Errorf("error creating request: %s", err)
	}

	req.Header.Set("Content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return respData, fmt.Errorf("error doing request: %s", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return respData, fmt.Errorf("error decoding response: %s", err)
	}
	return respData, nil
}

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type EmbeddingsReq struct {
	Prompt         string `json:"input"`
	Model          string `json:"model"`
	Dimensions     int    `json:"dimensions"`
	EncodingFormat string `json:"encoding_format"` // "float"
}

type EmbeddingsResp struct {
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
	Data []struct {
		Index     int       `json:"index"`
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

func (reqjson EmbeddingsReq) Send(baseurl, apikey string) (*EmbeddingsResp, error) {
	client := http.Client{
		Timeout: 600 * time.Second,
	}

	reqbyte, err := json.Marshal(reqjson)
	if err != nil {
		return nil, errors.New("couldn't marshal json")
	}
	req, err := http.NewRequest("POST", baseurl+"/embeddings", bytes.NewReader(reqbyte))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apikey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respdata, _ := io.ReadAll(resp.Body)
		return nil, errors.New("returned " + resp.Status + "\n" + string(respdata))
	}
	respdata, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("couldn't read response data")
	}
	var respjson EmbeddingsResp
	err = json.Unmarshal(respdata, &respjson)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal response")
	}
	return &respjson, nil
}

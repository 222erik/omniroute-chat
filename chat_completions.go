package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

// these two aren't used
type chatDevSystemUserFuncMessage struct {
	Role    string `json:"role"` // system, developer, user, function
	Content string `json:"content"`
	Name    string `json:"name"` // your name
}
type chatAssistantMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"` // always 'assistant'
	Refusal string `json:"refusal"`
}

type chatMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type ReqChat struct {
	Model       string        `json:"model"`
	Stream      bool          `json:"stream"`      // show the model's response while it's generating
	Temperature float32       `json:"temperature"` // the higher the more random
	ToolUsage   string        `json:"took_choice"` // none, auto, required, or type:function
	Verbosity   string        `json:"verbosity"`   // low, meduim, high
	Messages    []chatMessage `json:"messages"`
}

type RespChat struct {
	Id    string `json:"id"`
	Usage struct {
		CompletionsTokens int `json:"completion_tokens"`
		PromptTokens      int `json:"prompt_tokens"`
		TotalTokens       int `json:"total_tokens"`
	}
	Choices []struct { // every choice is a generated message my the model
		FinishReason string `json:"finish_reason"` // 'stop' for successful
		Index        int    `json:"index"`
		chatMessage  `json:"message"`
	}
}

func (reqjson ReqChat) Send(baseURL, APIkey string) (*RespChat, error) {
	client := http.Client{
		Timeout: 60 * time.Second,
	}

	var reqbyte []byte
	json.Unmarshal(reqbyte, reqjson)
	req, err := http.NewRequest("POST", baseURL, bytes.NewReader(reqbyte))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+APIkey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("returned " + resp.Status)
	}
	respdata, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("couldn't read response data")
	}
	var respjson RespChat
	err = json.Unmarshal(respdata, &respjson)
	if err != nil {
		return nil, errors.New("couldn't unmarshal response")
	}
	return &respjson, nil
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type RespModels struct {
	Object string `json:"object"`
	Data   []struct {
		Id           string `json:"id"`
		Created      int    `json:"created"`
		Object       string `json:"object"`
		OwnedBy      string `json:"owned_by"`
		ShutdownDate string `json:"shutdown_date"`
	} `json:"data"`
}

func (r *RespModels) Get(burl, apikey string) {
	client := &http.Client{
		Timeout: time.Minute,
	}

	req, err := http.NewRequest("GET", burl+"/models", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+apikey)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal("staus isn't ok")
	}

	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		log.Fatal(err)
	}
}

type ReqEmbeddings struct {
	Prompt         string `json:"input"`
	Model          string `json:"model"`
	Dimensions     int    `json:"dimensions"`
	EncodingFormat string `json:"encoding_format"` // Can only be "float" or "base64"
}

type RespEmbeddins struct {
	Model  string `json:"model"`
	Object string `json:"object"`
	Usage  struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
	Data struct {
		Index     int    `json:"index"`
		Object    string `json:"object"`
		Embedding any    `json:"embedding"`
	} `json:"data"`
}

func main() {
	BaseURL := os.Getenv("OMNIROUTE_BASE_URL")
	if BaseURL == "" {
		fmt.Println("base URL not set")
		BaseURL = "http://localhost:20128/v1"
	}

	APIKey := os.Getenv("OMNIROUTE_API_KEY")
	if APIKey == "" {
		fmt.Println("no API key given")
		os.Exit(1)
	}

	var models RespModels
	models.Get(BaseURL, APIKey)

	fmt.Println("There are", len(models.Data), "models available")

	options := []struct { // the functions are stored here for now
		desc     string
		function func()
	}{
		{desc: "List all models", function: func() {
			for _, d := range models.Data {
				fmt.Println(d.Id)
			}
		}},
		{desc: "Embed (turn a phrase into tokens) (not working)", function: func() {
			fmt.Println("List embedding modesl? y/N: ")
			var ans byte
			fmt.Scanf("%c", &ans)
			if ans == 'y' {
				for i, d := range models.Data {
					if strings.Contains(d.Id, "embedding") {
						fmt.Println(strconv.FormatInt(int64(i), 10)+".", d.Id)
					}
				}
				fmt.Println()
			}
			var req ReqEmbeddings
			fmt.Print("Enter index of model to use: ")
			var modelIndex int
			fmt.Scanf("%d", &modelIndex)
			fmt.Println()
			req.Model = models.Data[modelIndex].Id
			fmt.Print("Enter num of dimensions: ")
			fmt.Scanf("%d", &req.Dimensions)
			fmt.Println()
			fmt.Print("Enter your prompt: ")
			fmt.Scanf("%s", &req.Prompt)
			fmt.Println()
			req.EncodingFormat = "float"
		}},
		{desc: "Exit", function: func() {
			os.Exit(0)
		}},
	}

	for {
		fmt.Println()
		for i, s := range options {
			fmt.Println(strconv.FormatInt(int64(i+1), 10)+".", s.desc)
		}
		fmt.Print("What do you want to do? : ")

		var ans string
		fmt.Scanf("%s", &ans)
		fmt.Println()
		option, err := strconv.Atoi(ans)
		if err != nil || (option < 1 || option > len(options)) {
			fmt.Println("please enter valid number")
			continue
		}

		options[option-1].function()
	}
}

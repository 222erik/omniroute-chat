package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	tenv "github.com/muesli/termenv"
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

func (r *RespModels) Execute(output *tenv.Output) {
	leng := len(r.Data)
	for i := range leng {
		fmt.Println(r.Data[leng-i-1].Id)
	}
	time.Sleep(1000) // just so it doesnt exit immediatelly
	s := output.String("press enter when done")
	fmt.Println(s.Foreground(output.Color("#FFDDDD")))
	fmt.Scanf("%v", nil)
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

	output := tenv.NewOutput(os.Stdout)
	output.ClearScreen()

	options := []string{
		"List all models",
		"Chat with a model",
		"Embed something",
		"Exit",
	}

	for {
		for i, s := range options {
			fmt.Println(strconv.FormatInt(int64(i+1), 10)+".", s)
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

		output.ClearScreen()
		ExecuteCommand(option, output, BaseURL, APIKey)
		output.ClearScreen()
	}
}

func ExecuteCommand(c int, output *tenv.Output, baseurl, apikey string) {
	switch c {
	case 1:
		var models RespModels
		models.Get(baseurl, apikey)
		models.Execute(output)
	case 2:
		var messages []ChatMessage
		var model string
		resoning := "none"
		stream := false
		temperature := 0.0
		toolusage := "none"
		verbosity := "low"

		for model == "" {
			fmt.Print("Enter model ID: ")
			fmt.Scanf("%s", &model)
			fmt.Println()
		}

		var req ChatCompletionReq
		req.Model = model
		req.ReasoningEffort = resoning
		req.Stream = stream
		req.Temperature = float32(temperature)
		req.ToolUsage = toolusage
		req.Verbosity = verbosity
		req.Messages = messages

		fmt.Println("type 'exit' to quit")
		scanner := bufio.NewScanner(os.Stdin)
		for {
			s := output.String(">>> ")
			fmt.Print(s.Foreground(output.Color("#FFDDDD")))
			var prompt string
			scanner.Scan()
			prompt = scanner.Text()
			if prompt == "" {
				continue
			}
			if prompt == "exit" {
				break
			}
			messages = append(messages, ChatMessage{Role: "user", Content: prompt})
			req.Messages = messages

			resp, err := req.Send(baseurl, apikey)
			if err != nil {
				errstr := output.String(err.Error())
				fmt.Println(errstr.Foreground(output.Color("#FF8888")))
				continue
			}

			choice := resp.Choices[0]
			respstr := output.String(choice.Content)
			fmt.Println(respstr.Foreground(output.Color("#CCCCFF")))
			fmt.Println("done because:", choice.FinishReason)
			fmt.Println("usage:")
			fmt.Println("prompt:", resp.Usage.PromptTokens, "tokens")
			fmt.Println("completion:", resp.Usage.CompletionsTokens, "tokens")
			fmt.Println("total:", resp.Usage.TotalTokens, "tokens")
			fmt.Println()

			messages = append(messages, choice.ChatMessage)
		}
	case 3:
		var model string
		dimensions := 200
		format := "float"

		for model == "" {
			fmt.Print("Enter model ID: ")
			fmt.Scanf("%s", &model)
			fmt.Println()
		}

		var req EmbeddingsReq
		req.Model = model
		req.Dimensions = dimensions
		req.EncodingFormat = format

		fmt.Println("enter text to embed")
		scanner := bufio.NewScanner(os.Stdin)
		s := output.String(">>> ")
		fmt.Print(s.Foreground(output.Color("#FFDDDD")))
		var prompt string
		scanner.Scan()
		prompt = scanner.Text()
		req.Prompt = prompt

		resp, err := req.Send(baseurl, apikey)
		if err != nil {
			errstr := output.String(err.Error())
			fmt.Println(errstr.Foreground(output.Color("#FF8888")))
		}

		embedding := resp.Data[0]
		for i, val := range embedding.Embedding {
			respstr := output.String(strconv.FormatFloat(float64(val), 'f', -1, 32))
			fmt.Println(i, respstr.Foreground(output.Color("#CCCCFF")))
		}

		fmt.Println("usage:")
		fmt.Println("prompt:", resp.Usage.PromptTokens, "tokens")
		fmt.Println("total:", resp.Usage.TotalTokens, "tokens")
		fmt.Println("press enter to continue")
		fmt.Scanf("%d")
	case 4:
		os.Exit(0)
	}
}

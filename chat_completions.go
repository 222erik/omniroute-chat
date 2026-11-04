package main

type chatDevSystemUserFuncMessage struct {
	Role    string `json:"role"` // system, developer, user, function
	Content string `json:"content"`
	Name    string `json:"name"` // your name
}

type chatAssistantMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type ReqChat struct {
	Model       string                         `json:"model"`
	Stream      bool                           `json:"stream"`      // show the model's response while it's generating
	Temperature float32                        `json:"temperature"` // the higher the more random
	ToolUsage   string                         `json:"took_choice"` // none, auto, required, or type:function
	Verbosity   string                         `json:"verbosity"`   // low, meduim, high
	Messages    []chatDevSystemUserFuncMessage `json:"messages"`
}

type RespChat struct {
	Id    string `json:"id"`
	Usage struct {
		CompletionsTokens int `json:"completion_tokens"`
		PromptTokens      int `json:"prompt_tokens"`
		TotalTokens       int `json:"total_tokens"`
	}
	Choices []struct { // every choice is a generated message my the model
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	}
}

package utils

import (
	"fmt"
	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
)

type GPTTokenCount interface {
	NumTokensFromMessages(messages []openai.ChatCompletionMessage, model string) (numTokens int)
}

func NumTokensFromMessages(messages []openai.ChatCompletionMessage, model string) (numTokens int) {
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		err = fmt.Errorf("EncodingForModel: %v", err)
		fmt.Println(err)
		return
	}

	var tokensPerMessage int
	var tokensPerName int
	if model == "gpt-3.5-turbo-0301" || model == "gpt-3.5-turbo" {
		tokensPerMessage = 4
		tokensPerName = -1
	} else if model == "gpt-4-0314" || model == "gpt-4" {
		tokensPerMessage = 3
		tokensPerName = 1
	} else {
		fmt.Println("Warning: model not found. Using cl100k_base encoding.")
		tokensPerMessage = 3
		tokensPerName = 1
	}

	for _, message := range messages {
		numTokens += tokensPerMessage
		numTokens += len(tkm.Encode(message.Content, nil, nil))
		numTokens += len(tkm.Encode(message.Role, nil, nil))
		if message.Name != "" {
			numTokens += tokensPerName
		}
	}
	numTokens += 3
	return numTokens
}

func NumTokensFromText(text string, model string) (numTokens int) {
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		err = fmt.Errorf("EncodingForModel: %v", err)
		fmt.Println(err)
		return
	}

	var tokensPerMessage int
	if model == "gpt-3.5-turbo-0301" || model == "gpt-3.5-turbo" {
		tokensPerMessage = 4
	} else if model == "gpt-4-0314" || model == "gpt-4" {
		tokensPerMessage = 3
	} else {
		fmt.Println("Warning: model not found. Using cl100k_base encoding.")
		tokensPerMessage = 3
	}

	numTokens += tokensPerMessage
	numTokens += len(tkm.Encode(text, nil, nil))
	numTokens += 3
	return numTokens
}

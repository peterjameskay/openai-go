package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	// Adding info
	clientKey := os.Getenv("OPENAI_SECRET_KEY")
	client := openai.NewClient(clientKey)
	ctx := context.Background()
	keepAskingQuestions := true
	var conversationHistory []openai.ChatCompletionMessage

	for keepAskingQuestions {
		fmt.Println("\n" + `Please ask the bot a question. Type "exit" to exit.`)
		reader := bufio.NewReader(os.Stdin)
		userResponse, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("userResponse error: " + err.Error())
			return
		}

		conversationHistory = append(conversationHistory, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: userResponse,
		})

		req := openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: conversationHistory,
			Stream:   true,
		}

		stream, err := client.CreateChatCompletionStream(ctx, req)
		if err != nil {
			fmt.Printf("ChatCompletionStream error: %v\n", err)
			return
		}
		defer stream.Close()

		fmt.Printf("\nBot Response: ")
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				fmt.Println()
				break
			}

			if err != nil {
				fmt.Printf("\nStream error: %v\n", err)
				return
			}

			fmt.Printf(response.Choices[0].Delta.Content)
		}

		if strings.Compare(userResponse, "exit\n") == 0 {
			keepAskingQuestions = false
		}
	}
}

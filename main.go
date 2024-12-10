package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/openai/openai-go"
)

func main() {
	if len(os.Args) < 2 {
		println("Usage: send a query to 01 via typing something or cat a file")
		os.Exit(1)
	}

	client := openai.NewClient()
	input := readInput()
	query := os.Args[1]
	if len(input) > 0 {
		query = query + ": " + input
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	chatCompletion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(os.Args[1] + ": " + input),
		}),
		Model: openai.F(openai.ChatModelO1Mini),
	})
	if err != nil {
		log.Fatalf("Failed to get chat completion: %v ", err)
	}
	println(chatCompletion.Choices[0].Message.Content)
}

func readInput() string {
	info, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking stdin: %v\n", err)
		os.Exit(1)
	}
	// Check if stdin is a pipe or a regular file
	if info.Mode()&os.ModeCharDevice == 0 {
		reader := bufio.NewReader(os.Stdin)
		// lineNumber := 1
		b, err := io.ReadAll(reader)
		if err != nil {
			println("Error reading input")
			os.Exit(1)
		}
		return string(b)
	}
	return ""
}

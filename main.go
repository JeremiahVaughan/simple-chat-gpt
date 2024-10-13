package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	openai "github.com/sashabaranov/go-openai"
)

var chatClient *openai.Client

func init() {
	chatClient = openai.NewClient(os.Getenv("OPENAI_API_KEY"))
}

func main() {
	ctx := context.Background()
	p := tea.NewProgram(
		model{content: `eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good eat good`},
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		log.Fatalf("could not run program:", err)
	}
	err := submitChatMessage(ctx)
	if err != nil {
		log.Fatalf("error, when submitChatMessage() for main(). Error: %v", err)
	}
}

func submitChatMessage(ctx context.Context) error {
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT4oMini,
		MaxTokens: 4096,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "tell me a dad joke",
			},
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "no I won't do that",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "why not?",
			},
		},
		Stream: true,
	}

	stream, err := chatClient.CreateChatCompletionStream(ctx, req)
	if err != nil {
		log.Fatalf("error, when creating chat completion stream for main(). Error: %v\n", err)
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return fmt.Errorf("\nStream error: %v\n", err)
		}
		fmt.Printf(response.Choices[0].Delta.Content)
	}
}

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	openai "github.com/sashabaranov/go-openai"
)

var chatClient *openai.Client
var ctx context.Context

func init() {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		log.Fatal("your API key is missing, please add OPENAI_API_KEY as an environmental variable")
	}
	chatClient = openai.NewClient(key)
}

func main() {
	ctx = context.Background()
	p := tea.NewProgram(
		initModel(),
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		log.Fatalf("could not run program. Error: %v", err)
	}
}

func initModel() model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = ""
	ta.CharLimit = 1024

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	m := model{
		textarea:        ta,
		displayMessages: []string{},
		err:             nil,
	}
	m.resetSpinner()
	return m
}
func submitChatMessage(ctx context.Context, sendMessages []string) error {
	chatRequest := make([]openai.ChatCompletionMessage, len(sendMessages))
	for i := range sendMessages {
		var role string
		if isEven(i) {
			role = openai.ChatMessageRoleUser
		} else {
			role = openai.ChatMessageRoleSystem
		}
		chatRequest[i] = openai.ChatCompletionMessage{
			Role:    role,
			Content: sendMessages[i],
		}
	}
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT4oMini,
		MaxTokens: 4096,
		Messages:  chatRequest,
		Stream:    true,
	}

	stream, err := chatClient.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return fmt.Errorf("error, when creating chat completion stream for submitChatMessage(). Error: %v", err)
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return fmt.Errorf("error, when streaming. Error: %v", err)
		}
		streamingResponse <- response.Choices[0].Delta.Content
	}
}

func isEven(i int) bool {
	return i%2 == 0
}

func (m *model) resetSpinner() {
	s := spinner.New()
	s.Spinner = spinner.Moon
	m.spinner = s
}

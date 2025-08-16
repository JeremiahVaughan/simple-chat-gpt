package main

import (
	"fmt"
	"strings"
    "strconv"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// You generally won't need this unless you're processing stuff with
// complicated ANSI escape sequences. Turn it on if you notice flickering.
//
// Also keep in mind that high performance rendering only works for programs
// that use the full size of the terminal. We're enabling that below with
// tea.EnterAltScreen().
const useHighPerformanceRenderer = false

var (
	pageStyle   = lipgloss.NewStyle()
	senderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("201"))
	titleStyle  = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
	streamingResponse   chan string = make(chan string, 1000)
	loadingFinished     chan error  = make(chan error)
	loading             bool
	pasteBufferClosed   chan bool     = make(chan bool)
	messagesReadyToSend chan []string = make(chan []string, 1)
)

type model struct {
	recordedMessages []string
	sendMessages     []string
	displayMessages  []string
	ready            bool
	viewport         viewport.Model
	textarea         textarea.Model
	err              error
	spinner          spinner.Model
	currentRequest   string
	currentResponse  string
    tokensUsed int
    selectedAiModel string
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) headerView() string {
    text := fmt.Sprintf( "%s tokens used: %s / %s",
        m.selectedAiModel,
        formatNumberWithCommas(m.tokensUsed),
        formatNumberWithCommas(aiModelMap[m.selectedAiModel]),
    )
	title := titleStyle.Render(text)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func formatNumberWithCommas(num int) string {                         
    str := strconv.Itoa(num)                                          
    if len(str) < 4 {
        return str
    }
    reverseStr := reverse(str)                                        
    var sb strings.Builder                                            
    for i, c := range reverseStr {                                    
        if i > 0 && i%3 == 0 {                                        
            sb.WriteString(",")                                       
        }                                                             
        sb.WriteRune(c)                                               
    }                                                                 
    return reverse(sb.String())                                       
}                                                                     
                                                                      
func reverse(s string) string {                                       
    reversed := []rune(s)                                             
    for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {          
        reversed[i], reversed[j] = reversed[j], reversed[i]           
    }                                                                 
    return string(reversed)                                           
}                                                                     

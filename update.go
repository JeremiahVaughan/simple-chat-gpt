package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	chatMode := m.textarea.Focused() && !loading
	var enterChatMode bool
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// we are assuming the user has seen the error message if they are deciding to push more keys after it appears
		m.err = nil

		msgString := msg.String()
		switch msgString {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+n":
			if !loading {
				m.viewport.SetContent("")
				m.recordedMessages = []string{}
				m.sendMessages = []string{}
				m.displayMessages = []string{}
				m.textarea.Reset()
				m.textarea.Focus()
			}
			enterChatMode = true
		case "tab":
			if chatMode {
				m.textarea.Blur()
			} else if !loading {
				enterChatMode = true
			}
		case "enter":
			if loading {
				return m, nil
			}
			v := m.textarea.Value()
			if v == "" {
				return m, nil
			}
			m.currentResponse = ""
			m.currentRequest = v
			m.sendMessages = make([]string, len(m.recordedMessages)+1)
			for i, msg := range m.recordedMessages {
				m.sendMessages[i] = msg
			}
			m.sendMessages[len(m.recordedMessages)] = m.currentRequest
			m.displayMessages = make([]string, len(m.sendMessages)+1)
			for i, msg := range m.sendMessages {
				var dm string
				if isEven(i) {
					dm = senderStyle.Render(msg)
				} else {
					dm = msg
				}
				m.displayMessages[i] = dm
			}
			m.viewport.SetContent(generateViewportContent(m.displayMessages, messageDelimiter, m.viewport.Width))
			m.textarea.Reset()
			m.viewport.GotoBottom()
			loading = true
			go func() {
				err := submitChatMessage(ctx, m.sendMessages)
				if err != nil {
					loadingFinished <- fmt.Errorf("error, when submitChatMessage() for update(). Error: %v", err)
					return
				}
				loadingFinished <- nil
			}()
			return m, m.spinner.Tick
		case "i", "a", "s":
			if !chatMode {
				enterChatMode = true
			}
		case "ctrl+u":
			if !chatMode {
				newOffset := m.viewport.YOffset - (m.viewport.Height / 2)
				if newOffset < 0 {
					newOffset = 0
				}
				m.viewport.SetYOffset(newOffset)
			}
		case "ctrl+d":
			if !chatMode {
				newOffset := m.viewport.YOffset + (m.viewport.Height / 2)
				if newOffset > m.viewport.TotalLineCount() {
					newOffset = m.viewport.TotalLineCount()
				}
				m.viewport.SetYOffset(newOffset)
			}
		case "g":
			if !chatMode {
				m.viewport.GotoTop()
			}
		case "G":
			if !chatMode {
				m.viewport.GotoBottom()
			}
		case "k":
			if !chatMode && !m.viewport.AtTop() {
				m.viewport.SetYOffset(m.viewport.YOffset - 1)
			}
		case "j":
			if !chatMode && !m.viewport.AtBottom() {
				m.viewport.SetYOffset(m.viewport.YOffset + 1)
			}
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight
		viewportWidth := 70
		viewportHeight := msg.Height - (verticalMarginHeight + m.textarea.Height())
		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(viewportWidth, viewportHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(generateViewportContent(m.displayMessages, messageDelimiter, m.viewport.Width))
			m.ready = true
			rightMarginWidth := (msg.Width - viewportWidth) / 2
			pageStyle = pageStyle.MarginLeft(rightMarginWidth)

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
			m.textarea.SetWidth(viewportWidth)

		} else {
			m.viewport.Width = viewportWidth
			m.viewport.Height = viewportHeight
			m.viewport.SetContent(generateViewportContent(m.displayMessages, messageDelimiter, m.viewport.Width))
		}

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	case spinner.TickMsg:
		m.currentResponse = drainAllResponses(m.currentResponse, streamingResponse)
		if m.currentResponse != "" {
			m.displayMessages[len(m.displayMessages)-1] = m.currentResponse
			m.viewport.SetContent(generateViewportContent(m.displayMessages, messageDelimiter, m.viewport.Width))
		}
		select {
		case err := <-loadingFinished:
			m.resetSpinner()
			loading = false
			enterChatMode = true
			if err == nil {
				m.recordedMessages = append(
					m.recordedMessages,
					m.currentRequest,
					m.currentResponse,
				)
			}
		default:
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	if chatMode {
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}
	if enterChatMode {
		m.textarea.Focus()
		newBlink := cursor.Blink()
		m.textarea, cmd = m.textarea.Update(newBlink)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func drainAllResponses(currentResponse string, allResponses chan string) string {
	for {
		select {
		case resp := <-allResponses:
			currentResponse += resp
		default:
			return currentResponse
		}
	}
}

func generateViewportContent(displayMessages []string, messageDelimiter string, viewportWidth int) string {
	for i, msg := range displayMessages {
		displayMessages[i] = applyWordWrap(msg, viewportWidth)
	}
	return "\n" + strings.Join(displayMessages, messageDelimiter)
}

func applyWordWrap(msg string, viewportWidth int) string {
	blocks := strings.Split(msg, "\n")
	out := []string{}
	for _, block := range blocks {
		formatted, cnt := replaceIndent(block)
		wrapped := overflowWrap(formatted, viewportWidth)
		out = append(out, repairIndent(wrapped, cnt))
	}

	return strings.Join(out, "\n")
}

func replaceIndent(s string) (string, int) {
	cnt := 0
	for _, char := range s {
		if char != ' ' {
			break
		}
		cnt += 1
	}
	return strings.Repeat("c", cnt) + s[cnt:], cnt
}

func repairIndent(s string, cnt int) string {
	return strings.Repeat(" ", cnt) + s[cnt:]
}

/*
Assumes input is a trimmed string with only a single white space between words
*/
func overflowWrap(s string, maxWidth int) string {
	if maxWidth < 1 {
		panic("assertion failed: maxWidth must be at least 1")
	}
	words := strings.Fields(s)
	lines := []string{}
	var line strings.Builder
	for _, word := range words {
		if line.Len()+1+len(word) >= maxWidth {
			if line.Len() > 0 {
				lines = append(lines, line.String())
				line = strings.Builder{}
			}
			for len(word) >= maxWidth {
				left, right := breakWord(word, maxWidth)
				lines = append(lines, left)
				word = right
			}
			if len(word) > 0 {
				line.WriteString(word)
			}
		} else {
			if line.Len() > 0 {
				line.WriteRune(' ')
			}
			line.WriteString(word)
		}
	}
	if line.Len() > 0 {
		lines = append(lines, line.String())
	}
	return strings.Join(lines, "\n")
}

func breakWord(word string, maxWidth int) (string, string) {
	return word[:maxWidth], word[maxWidth:]
}

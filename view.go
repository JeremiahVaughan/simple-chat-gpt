package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var messageDelimiter = "\n\n"

func (m model) View() string {
	if m.err != nil {
		return getErrorStyle(m.err.Error())
	}

	if !m.ready {
		return "\n  Initializing..."
	}
	var textareaDisplay string
	if loading {
		textareaDisplay = fmt.Sprintf("%s Chat Jipity is speaking...", m.spinner.View())
	} else {
		textareaDisplay = m.textarea.View()
	}
	pageContent := fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		m.headerView(m.tokensUsed, maxTokens),
		m.viewport.View(),
		m.footerView(),
		textareaDisplay,
	)
	return pageStyle.Render(pageContent)
}

func getErrorStyle(errMsg string) string {
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true).Width(80).MarginLeft(4)
	return fmt.Sprintf("\n\n%v", errorStyle.Render(errMsg))
}

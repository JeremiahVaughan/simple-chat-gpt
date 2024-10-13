package main

import "fmt"

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	pageContent := fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		m.headerView(),
		m.viewport.View(),
		m.footerView(),
		m.textarea.View(),
	)
	return pageStyle.Render(pageContent)
}

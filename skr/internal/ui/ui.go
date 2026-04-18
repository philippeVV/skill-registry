package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	accent  = lipgloss.Color("6")  // cyan
	success = lipgloss.Color("2")  // green
	warn    = lipgloss.Color("3")  // yellow
	danger  = lipgloss.Color("1")  // red
	muted   = lipgloss.Color("8")  // gray

	// Styles
	Title   = lipgloss.NewStyle().Bold(true).Foreground(accent)
	Success = lipgloss.NewStyle().Foreground(success)
	Warning = lipgloss.NewStyle().Foreground(warn)
	Error   = lipgloss.NewStyle().Foreground(danger)
	Muted   = lipgloss.NewStyle().Foreground(muted)
	Bold    = lipgloss.NewStyle().Bold(true)

	// Badges
	BadgeModified = lipgloss.NewStyle().Foreground(warn).Render("[modified]")
	BadgeSystem   = lipgloss.NewStyle().Foreground(muted).Render("[system]")
)

// PackageTable prints a formatted table of packages.
func PackageTable(headers []string, rows [][]string) string {
	if len(rows) == 0 {
		return Muted.Render("  No results.")
	}

	// Compute column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	var b strings.Builder

	// Header
	for i, h := range headers {
		if i > 0 {
			b.WriteString("  ")
		}
		b.WriteString(Bold.Render(padRight(h, widths[i])))
	}
	b.WriteString("\n")

	// Separator
	for i, w := range widths {
		if i > 0 {
			b.WriteString("  ")
		}
		b.WriteString(Muted.Render(strings.Repeat("─", w)))
	}
	b.WriteString("\n")

	// Rows
	for _, row := range rows {
		for i, cell := range row {
			if i > 0 {
				b.WriteString("  ")
			}
			if i < len(widths) {
				b.WriteString(padRight(cell, widths[i]))
			} else {
				b.WriteString(cell)
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}

// InfoBlock prints a labeled set of key-value pairs.
func InfoBlock(pairs [][]string) string {
	maxKey := 0
	for _, p := range pairs {
		if len(p[0]) > maxKey {
			maxKey = len(p[0])
		}
	}

	var b strings.Builder
	for _, p := range pairs {
		key := Bold.Render(padRight(p[0], maxKey))
		b.WriteString(fmt.Sprintf("  %s  %s\n", key, p[1]))
	}
	return b.String()
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

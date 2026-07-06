package render

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"github.com/allanmedeiros71/criminalsay/internal/quotes"
)

var (
	sidebarStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	quoteStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	authorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	dimStyle     = lipgloss.NewStyle().Faint(true)
)

const (
	sidebarBlock = '\u258c'
	sidebarPlain = '|'
	metaMarker   = " · Criminal Minds"
)

func buildAttribution(q quotes.Quote) (authorPart, metaPart string) {
	if q.Author != q.Character && q.Character != "" {
		authorPart = q.Author + ", cited by " + q.Character
	} else {
		authorPart = q.Author
	}
	metaPart = metaMarker
	if q.Season != 0 && q.Episode != 0 {
		metaPart += fmt.Sprintf(" · S%dE%d", q.Season, q.Episode)
		if q.EpisodeTitle != "" {
			metaPart += fmt.Sprintf(" (%s)", q.EpisodeTitle)
		}
	}
	return authorPart, metaPart
}

func sidebarChar(colorEnabled bool) rune {
	if colorEnabled {
		return sidebarBlock
	}
	return sidebarPlain
}

func styleSidebar(ch rune, enabled bool) string {
	if enabled {
		return sidebarStyle.Render(string(ch))
	}
	return string(ch)
}

func renderAttributionLine(line string, colorEnabled bool) string {
	if !colorEnabled {
		return line
	}
	idx := strings.Index(line, metaMarker)
	if idx < 0 {
		return authorStyle.Render(line)
	}
	return authorStyle.Render(line[:idx]) + dimStyle.Render(line[idx:])
}

func formatQuoteLine(inner string, colorEnabled bool) string {
	ch := sidebarChar(colorEnabled)
	prefix := quoteIndent + styleSidebar(ch, colorEnabled) + sidebarGap
	if !colorEnabled {
		return prefix + inner
	}
	openIdx := strings.Index(inner, `"`)
	closeIdx := strings.LastIndex(inner, `"`)
	if openIdx >= 0 && closeIdx > openIdx {
		return prefix + inner[:openIdx] + quoteStyle.Render(inner[openIdx:closeIdx+1]) + inner[closeIdx+1:]
	}
	return prefix + quoteStyle.Render(inner)
}

func buildQuoteBlock(q quotes.Quote, colorEnabled bool) string {
	innerLines := wrapQuoteText(q.Quote)

	var quoteLines []string
	for _, inner := range innerLines {
		quoteLines = append(quoteLines, formatQuoteLine(inner, colorEnabled))
	}

	authorPart, metaPart := buildAttribution(q)
	fullAttr := authorPart + metaPart
	attrLines := wrapAttribution(fullAttr)
	var renderedAttr []string
	for _, al := range attrLines {
		renderedAttr = append(renderedAttr, attrIndent+renderAttributionLine(al, colorEnabled))
	}

	return strings.Join(quoteLines, "\n") + "\n\n" + strings.Join(renderedAttr, "\n")
}

// Quote formats q as Estilo 4 terminal output and returns the string.
// Does not print or detect TTY — caller passes colorEnabled from main.
func Quote(q quotes.Quote, colorEnabled bool) string {
	if colorEnabled {
		prev := lipgloss.ColorProfile()
		lipgloss.SetColorProfile(termenv.TrueColor)
		defer lipgloss.SetColorProfile(prev)
	}
	return buildQuoteBlock(q, colorEnabled)
}

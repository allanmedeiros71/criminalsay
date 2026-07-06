package render

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

const (
	blockWidth  = 52
	quoteIndent = "  "
	attrIndent  = "       "
	sidebarGap  = " "
)

// linePrefixWidth is quoteIndent + sidebar + sidebarGap (4 display columns).
const linePrefixWidth = 4

func wrapWords(text string, maxWidth int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var current strings.Builder
	currentWidth := 0

	flush := func() {
		if current.Len() > 0 {
			lines = append(lines, current.String())
			current.Reset()
			currentWidth = 0
		}
	}

	for _, word := range words {
		wordWidth := ansi.StringWidth(word)
		spaceWidth := 0
		if current.Len() > 0 {
			spaceWidth = 1
		}
		if currentWidth+spaceWidth+wordWidth > maxWidth {
			flush()
		}
		if current.Len() > 0 {
			current.WriteByte(' ')
			currentWidth++
		}
		current.WriteString(word)
		currentWidth += wordWidth
	}
	flush()
	return lines
}

func wrapQuoteText(text string) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{`""`}
	}

	// Single-line budget: 52 - prefix(4) - open(1) - close(1) = 46
	innerSingle := blockWidth - linePrefixWidth - 2
	innerFirst := blockWidth - linePrefixWidth - 1
	innerMiddle := blockWidth - linePrefixWidth
	innerLast := blockWidth - linePrefixWidth - 1

	var innerLines []string
	var current strings.Builder
	currentWidth := 0
	limit := innerFirst

	flush := func() {
		if current.Len() > 0 {
			innerLines = append(innerLines, current.String())
			current.Reset()
			currentWidth = 0
		}
	}

	for i, word := range words {
		wordWidth := ansi.StringWidth(word)
		spaceWidth := 0
		if current.Len() > 0 {
			spaceWidth = 1
		}
		if currentWidth+spaceWidth+wordWidth > limit {
			flush()
			if len(innerLines) == 1 {
				limit = innerLast
			} else {
				limit = innerMiddle
			}
		}
		if current.Len() > 0 {
			current.WriteByte(' ')
			currentWidth++
		}
		current.WriteString(word)
		currentWidth += wordWidth

		// After first flush, subsequent lines use middle/last limits.
		if i == len(words)-1 && len(innerLines) == 0 && currentWidth <= innerSingle {
			limit = innerSingle
		}
	}
	flush()

	if len(innerLines) == 1 {
		return []string{`"` + innerLines[0] + `"`}
	}

	innerLines[0] = `"` + innerLines[0]
	last := len(innerLines) - 1
	innerLines[last] = innerLines[last] + `"`

	return innerLines
}

func wrapAttribution(text string) []string {
	return wrapWords(text, blockWidth-ansi.StringWidth(attrIndent))
}

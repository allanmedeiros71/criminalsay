package render_test

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"

	"github.com/allanmedeiros71/criminalsay/internal/quotes"
	"github.com/allanmedeiros71/criminalsay/internal/render"
)

func assertNoANSI(t *testing.T, s string) {
	t.Helper()
	if strings.Contains(s, "\x1b[") {
		t.Errorf("output contains ANSI escape sequences: %q", s)
	}
}

func rossiFixture() quotes.Quote {
	return quotes.Quote{
		Quote:        "It's alchemy. Alchemy turns common metals into precious ones.",
		Author:       "David Rossi",
		Character:    "David Rossi",
		Season:       9,
		Episode:      13,
		EpisodeTitle: "The Road Home",
	}
}

func TestQuote_noColor(t *testing.T) {
	out := render.Quote(rossiFixture(), false)
	assertNoANSI(t, out)
	if !strings.Contains(out, "|") {
		t.Error("expected pipe sidebar in no-color output")
	}
	if strings.Contains(out, "\u258c") {
		t.Error("no-color output must not contain block sidebar")
	}
	if !strings.Contains(out, "David Rossi") || !strings.Contains(out, "S9E13") {
		t.Errorf("expected author and episode in output: %q", out)
	}
}

func TestQuote_wrapWidth(t *testing.T) {
	word := "alchemy"
	var words []string
	for i := 0; i < 30; i++ {
		words = append(words, word)
	}
	q := quotes.Quote{
		Quote:     strings.Join(words, " "),
		Author:    "David Rossi",
		Character: "David Rossi",
		Season:    9,
		Episode:   13,
	}
	out := render.Quote(q, false)
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "|") || strings.Contains(line, "\u258c") {
			if w := ansi.StringWidth(line); w > 52 {
				t.Errorf("line width %d > 52: %q", w, line)
			}
		}
	}
}

func TestQuote_externalAuthor(t *testing.T) {
	q := quotes.Quote{
		Quote:     "Test quote.",
		Author:    "Marcus Aurelius",
		Character: "Spencer Reid",
		Season:    1,
		Episode:   1,
	}
	out := render.Quote(q, false)
	if !strings.Contains(out, "Marcus Aurelius, cited by Spencer Reid") {
		t.Errorf("expected cited-by attribution: %q", out)
	}
}

func TestQuote_sameAuthor(t *testing.T) {
	out := render.Quote(rossiFixture(), false)
	if strings.Contains(out, "cited by") {
		t.Error("same author/character must not include cited by")
	}
	if strings.Count(out, "David Rossi") < 1 {
		t.Error("expected author name in attribution")
	}
}

func TestQuote_spacing(t *testing.T) {
	out := render.Quote(rossiFixture(), false)
	if !strings.Contains(out, "\n\n") {
		t.Error("expected blank line between quote block and attribution")
	}
}

func TestQuote_sidebar(t *testing.T) {
	out := render.Quote(rossiFixture(), false)
	parts := strings.Split(out, "\n\n")
	quoteBlock := parts[0]
	for _, line := range strings.Split(quoteBlock, "\n") {
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "  |") && !strings.HasPrefix(line, "  \u258c") {
			t.Errorf("quote line missing sidebar prefix: %q", line)
		}
	}
}

func TestQuote_optionalFields(t *testing.T) {
	q := quotes.Quote{
		Quote:     "Short.",
		Author:    "Author Only",
		Character: "",
		Season:    0,
		Episode:   0,
	}
	out := render.Quote(q, false)
	if strings.Contains(out, "S0E") {
		t.Error("zero season/episode must omit SxEy segment")
	}
	if strings.Contains(out, "cited by") {
		t.Error("empty character must omit cited-by clause")
	}
}

func TestQuote_color(t *testing.T) {
	out := render.Quote(rossiFixture(), true)
	if !strings.Contains(out, "\x1b[") {
		t.Error("color output should contain ANSI sequences")
	}
}

func TestQuote_noBox(t *testing.T) {
	boxChars := []rune{'\u250c', '\u2510', '\u2514', '\u2518', '\u2502', '\u2500'}
	out := render.Quote(rossiFixture(), true)
	for _, ch := range boxChars {
		if strings.ContainsRune(out, ch) {
			t.Errorf("output contains box-drawing char %q", ch)
		}
	}
}

func TestQuote_quotes(t *testing.T) {
	out := render.Quote(rossiFixture(), false)
	if !strings.Contains(out, `"`) {
		t.Error("output must wrap quote text in double quotes")
	}
}

func TestQuote_attribution(t *testing.T) {
	out := render.Quote(rossiFixture(), false)
	if !strings.Contains(out, " · Criminal Minds · ") {
		t.Error("attribution must contain Criminal Minds separator")
	}
	if !strings.Contains(out, "S9E13") {
		t.Error("attribution must contain episode code")
	}
}

func TestQuote_episodeTitle(t *testing.T) {
	out := render.Quote(rossiFixture(), false)
	if !strings.Contains(out, "Road Home") {
		t.Errorf("episode title must appear in attribution, got: %q", out)
	}
}

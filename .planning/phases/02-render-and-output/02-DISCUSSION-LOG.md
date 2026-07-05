# Phase 2: Render and Output - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-07-05
**Phase:** 2-Render and Output
**Areas discussed:** Attribution line, Word wrap & layout, Optional fields, Color & compatibility

---

## Attribution Line

| Option | Description | Selected |
|--------|-------------|----------|
| cited_by | `<Author>, cited by <Character> · Criminal Minds · SxEy` | ✓ |
| author_only | Author name only, ignore character | |
| character_first | Highlight CM character who spoke | |

| Option | Description | Selected |
|--------|-------------|----------|
| author_only_same | Single name when author == character | ✓ |
| show_both | Redundant dual display | |

| Option | Description | Selected |
|--------|-------------|----------|
| title_paren | Include `(Episode Title)` after SxEy when present | ✓ |
| no_title | SxEy only | |
| title_only_if_short | Conditional inclusion | |

| Option | Description | Selected |
|--------|-------------|----------|
| author_yellow | Only author segment yellow; rest dim | ✓ |
| all_yellow | Entire line yellow | |
| author_character_yellow | Author + character yellow | |

**User's choice:** External quotes use "cited by" format; same-author shows name once; episode title in parens when present; yellow on author portion only.

---

## Word Wrap & Layout

| Option | Description | Selected |
|--------|-------------|----------|
| fixed_52 | 52 useful columns | ✓ |
| fixed_60 | 60 columns | |
| adaptive | Terminal width | |

| Option | Description | Selected |
|--------|-------------|----------|
| indent_2 | 2 spaces before sidebar | ✓ |
| indent_4 | 4 spaces | |
| indent_0 | No indent | |

| Option | Description | Selected |
|--------|-------------|----------|
| sidebar_each | ▌ on every wrapped line | ✓ |
| sidebar_first | ▌ on first line only | |

| Option | Description | Selected |
|--------|-------------|----------|
| quotes_wrap | Quotes wrap inside block | ✓ |
| quotes_per_line | Quotes per line | |

**User's choice:** Fixed 52 cols, 2-space indent, sidebar on all lines, quotes wrap inside block.

---

## Optional Fields

| Option | Description | Selected |
|--------|-------------|----------|
| omit_silent | Omit empty character silently | ✓ |
| omit_unknown | Show "Unknown character" | |
| treat_as_author | Treat empty as author==character | |

| Option | Description | Selected |
|--------|-------------|----------|
| omit_zero | Omit SxEy when season/episode is 0 | ✓ |
| show_s0e0 | Show S0E0 | |
| show_partial | Partial episode code | |

| Option | Description | Selected |
|--------|-------------|----------|
| omit_empty | Omit parens when title empty | ✓ |
| show_unknown | Show (Unknown episode) | |

| Option | Description | Selected |
|--------|-------------|----------|
| wrap_attr | Wrap long attribution at 52 cols | ✓ |
| truncate_attr | Ellipsis truncation | |
| no_wrap_attr | Allow overflow | |

**User's choice:** Graceful omission for empty fields; wrap long attribution lines.

---

## Color & Compatibility

| Option | Description | Selected |
|--------|-------------|----------|
| lipgloss_auto | main passes colorEnabled from IsTerminal | ✓ |
| explicit_all | Explicit NO_COLOR + TTY + pipe checks | |
| lipgloss_only | Trust lipgloss only | |

| Option | Description | Selected |
|--------|-------------|----------|
| plain_ascii | Replace ▌ with \| in no-color mode | ✓ |
| strip_ansi | Keep ▌, strip ANSI only | |

| Option | Description | Selected |
|--------|-------------|----------|
| white_bright | Quote text white/bright | ✓ |
| default_fg | Terminal default | |
| bold_only | Bold without color change | |

| Option | Description | Selected |
|--------|-------------|----------|
| red_sidebar | Red ▌ in color mode | ✓ |
| red_dim | Dim red | |
| accent_match | Yellow sidebar | |

**User's choice:** lipgloss + IsTerminal/NO_COLOR in main; `|` instead of `▌` without color; white quote text; red sidebar.

**Note:** D-18 (`|` replacement) deviates from literal COLOR-02 requirement text — captured as intentional user override in CONTEXT.md.

---

## Claude's Discretion

- lipgloss helper structure and wrap implementation details
- Test decomposition for render package

## Deferred Ideas

None.

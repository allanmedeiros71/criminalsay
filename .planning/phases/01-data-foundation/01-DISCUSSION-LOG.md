# Phase 1: Data Foundation - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-07-05
**Phase:** 1-Data Foundation
**Areas discussed:** Character field semantics, Quotes curation approach, main.go stub behavior

---

## Character Field Semantics

**Q1: When a CM character originates a quote (not citing an outside author), what goes in each field?**

| Option | Description | Selected |
|--------|-------------|----------|
| author=character name, character empty | Original quotes have author set to character name, character left blank | |
| Both fields filled, same value | author and character both set to the same value (e.g., "David Rossi") | ✓ |
| author=character name, character='David Rossi' | Always populate character; same value when original | |

**User's choice:** Both fields filled with same value
**Notes:** When Rossi says something original, `author: "David Rossi"`, `character: "David Rossi"`.

---

**Q2: When author ≠ character, what attribution line should render?**

| Option | Description | Selected |
|--------|-------------|----------|
| Just the original author | Show "Blaise Pascal · Criminal Minds · S..." | |
| Character name only | Show "Spencer Reid · Criminal Minds · S..." | |
| Both: original + cited-by | Show "Blaise Pascal, cited by Reid · Criminal Minds · S..." | ✓ |

**User's choice:** Both — full attribution showing original author and CM character who cited it
**Notes:** Phase 2 rendering decision, but quotes.json must populate both fields accurately.

---

## Quotes Curation Approach

**Q1: How should 15+ real Criminal Minds quotes get into quotes.json?**

| Option | Description | Selected |
|--------|-------------|----------|
| Claude pre-populates from training knowledge | Claude fills quotes.json; user reviews/edits afterward | ✓ |
| You provide the quotes manually | quotes.json starts as stub; user pastes real quotes | |
| Stub for tests, real quotes added later | 3-5 placeholder entries for tests; real quotes before shipping | |

**User's choice:** Claude pre-populates from training knowledge
**Notes:** User will review and validate accuracy of episode/season attributions.

---

**Q2: Which characters/seasons should the initial quotes prioritize?**

| Option | Description | Selected |
|--------|-------------|----------|
| Spread across all main characters | Mix of Reid, Hotch, Morgan, Rossi, Garcia, Prentiss, JJ | ✓ |
| Focus on Reid and Rossi | Quote-heavy characters for quality density | |
| You decide the best spread | Claude picks most memorable regardless of balance | |

**User's choice:** Spread across all main characters
**Notes:** Representative of the full BAU team across multiple seasons.

---

## main.go Stub Behavior

**Q1: What should running 'criminalsay' do at the end of Phase 1?**

| Option | Description | Selected |
|--------|-------------|----------|
| Print raw quote text (unformatted) | fmt.Println(q.Quote + " — " + q.Author) | ✓ |
| Print JSON of the selected quote | json.Marshal or fmt.Printf("%+v\n", q) | |
| Just exit 0 silently | Load+Random called; result discarded; exits cleanly | |

**User's choice:** Print raw quote text (unformatted)
**Notes:** Replaced entirely by Phase 2's render package. Acts as visible confirmation the data layer works.

---

## Claude's Discretion

- Exact episode and season numbers for individual quotes (Claude selects best-known episodes with confident attribution; user reviews)
- Number of quotes above 15-quote minimum (Claude may include up to 25)

## Deferred Ideas

- Quote language field (PT-BR translation) — out of scope for v1
- Character filter flag (`--character`) — v2
- Season filter flag (`--season`) — v2
- Adaptive terminal width — Phase 2 concern

package quotes

import (
	"encoding/json"
	"math/rand/v2"
)

// Quote represents a single Criminal Minds quote record.
// Fields match the data/quotes.json schema exactly.
type Quote struct {
	Quote        string `json:"quote"`
	Author       string `json:"author"`
	Character    string `json:"character"`
	Season       int    `json:"season"`
	Episode      int    `json:"episode"`
	EpisodeTitle string `json:"episodeTitle"`
}

// Load parses JSON-encoded quote data and returns the slice.
// data should be the embedded bytes from main.go.
// Returns a non-nil error if JSON is malformed.
// Returns ([]Quote{}, nil) for a valid empty array — caller must check length.
func Load(data []byte) ([]Quote, error) {
	var qs []Quote
	if err := json.Unmarshal(data, &qs); err != nil {
		return nil, err
	}
	return qs, nil
}

// Random returns a uniformly random element from qs.
// Panics if qs is empty — caller must ensure len(qs) > 0.
func Random(qs []Quote) Quote {
	return qs[rand.IntN(len(qs))]
}

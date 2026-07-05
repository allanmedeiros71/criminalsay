package quotes_test

import (
	"testing"

	"github.com/allanmedeiros71/criminalsay/internal/quotes"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
		wantLen int
	}{
		{
			name:    "valid single quote",
			input:   []byte(`[{"quote":"Test","author":"A","character":"A","season":1,"episode":1,"episodeTitle":"Pilot"}]`),
			wantErr: false,
			wantLen: 1,
		},
		{
			name:    "valid multiple quotes",
			input:   []byte(`[{"quote":"Q1","author":"A1","character":"A1","season":1,"episode":1,"episodeTitle":""},{"quote":"Q2","author":"A2","character":"A2","season":2,"episode":3,"episodeTitle":""}]`),
			wantErr: false,
			wantLen: 2,
		},
		{
			name:    "empty array",
			input:   []byte(`[]`),
			wantErr: false,
			wantLen: 0,
		},
		{
			name:    "malformed JSON",
			input:   []byte(`not json`),
			wantErr: true,
		},
		{
			name:    "missing optional fields",
			input:   []byte(`[{"quote":"Q","author":"A","character":"","season":1,"episode":1,"episodeTitle":""}]`),
			wantErr: false,
			wantLen: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := quotes.Load(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Load() error = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr && len(got) != tc.wantLen {
				t.Errorf("Load() len = %d, want %d", len(got), tc.wantLen)
			}
		})
	}
}

func TestRandom_ReturnsItemFromSlice(t *testing.T) {
	qs := []quotes.Quote{
		{Quote: "A", Author: "X"},
		{Quote: "B", Author: "Y"},
		{Quote: "C", Author: "Z"},
	}
	q := quotes.Random(qs)
	found := false
	for _, item := range qs {
		if item.Quote == q.Quote {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Random() returned item not in input slice: %+v", q)
	}
}

func TestRandom_Distribution(t *testing.T) {
	// Run enough iterations to confirm non-constant selection.
	qs := []quotes.Quote{
		{Quote: "A"}, {Quote: "B"}, {Quote: "C"},
		{Quote: "D"}, {Quote: "E"}, {Quote: "F"},
	}
	seen := make(map[string]bool)
	for i := 0; i < 200; i++ {
		q := quotes.Random(qs)
		seen[q.Quote] = true
	}
	if len(seen) < 2 {
		t.Errorf("Random() appears non-random: only %d unique values in 200 calls", len(seen))
	}
}

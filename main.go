package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/allanmedeiros71/criminalsay/internal/quotes"
)

//go:embed data/quotes.json
var quoteData []byte

func main() {
	qs, err := quotes.Load(quoteData)
	if err != nil {
		fmt.Fprintln(os.Stderr, "criminalsay: failed to load quotes:", err)
		os.Exit(1)
	}
	if len(qs) == 0 {
		fmt.Fprintln(os.Stderr, "criminalsay: no quotes available")
		os.Exit(1)
	}
	q := quotes.Random(qs)
	fmt.Println(q.Quote + " — " + q.Author)
}

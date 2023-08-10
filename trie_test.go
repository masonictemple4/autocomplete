package autocomplete

import (
	"fmt"
	"testing"
)

func TestTrie(t *testing.T) {
	trie := newTrie()

	words := []string{"bike", "bike path", "bicycle repair", "pool", "beach", "waterfront", "dog park", "resteraunts"}

	for _, word := range words {
		trie.Insert(word)
	}

	// Test ListContents.

	contents := trie.ListContents()

	if len(contents) != len(words) {
		t.Errorf("Expected %d words, got %d", len(words), len(contents))
	}

	fmt.Printf("The contents: %v\n", contents)

	results := trie.Autocomplete("bi")
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	fmt.Printf("The results: %v\n", results)

}

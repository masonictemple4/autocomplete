package autocomplete

import (
	"fmt"
	"testing"
)

func TestTrie(t *testing.T) {
	trie := newTrie()

	words := []string{"hello", "help", "helium", "helicopter", "helipad", "heaven"}

	for _, word := range words {
		trie.Insert(word)
	}

	// Test ListContents.

	contents := trie.ListContents()

	if len(contents) != len(words) {
		t.Errorf("Expected %d words, got %d", len(words), len(contents))
	}

	fmt.Printf("The contents: %v\n", contents)

}

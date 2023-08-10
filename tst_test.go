package autocomplete

import (
	"fmt"
	"testing"
)

func TestTernarySearchTree(t *testing.T) {
	t.Run("empty setup", func(t *testing.T) {
		tree := newTernarySearchTree("")

		contents := tree.ListContents()
		if len(contents) > 0 {
			t.Errorf("Expected %d words, got %d", 0, len(contents))
		}

		fmt.Printf("The contents: %v\n", contents)
	})

	t.Run("basic setup", func(t *testing.T) {
		words := []string{"bike", "bike path", "bicycle repair", "pool", "beach", "waterfront", "dog park", "resteraunts"}

		tree := newTernarySearchTree("")

		for _, word := range words {
			tree.Insert(word)
		}

		// Test ListContents.
		contents := tree.ListContents()
		if len(contents) != len(words) {
			t.Errorf("Expected %d words, got %d", len(words), len(contents))
		}

		fmt.Printf("The contents: %v\n", contents)

		results := tree.Autocomplete("bi")
		if len(results) != 3 {
			t.Errorf("Expected 3 results, got %d", len(results))
		}

	})

}

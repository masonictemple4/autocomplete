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
		words := []string{"help", "helium", "helicopter", "helipad", "heaven"}

		tree := newTernarySearchTree("hello")

		for _, word := range words {
			tree.Insert(word)
		}

		// Test ListContents.
		contents := tree.ListContents()
		if len(contents) != len(words)+1 {
			t.Errorf("Expected %d words, got %d", len(words)+1, len(contents))
		}

		fmt.Printf("The contents: %v\n", contents)
	})

}

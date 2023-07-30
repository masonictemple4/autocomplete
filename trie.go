package autocomplete

import "sync"

// Make sure we implement the auto completer
var _ autocompleter = (*trie)(nil)

type trieNode struct {
	// Using rune for future extensibility
	children map[rune]*trieNode
	isEnd    bool
}

type trie struct {
	Root *trieNode

	mu sync.RWMutex
}

func newTrie() *trie {
	return &trie{
		Root: &trieNode{children: make(map[rune]*trieNode)},
	}
}

func (t *trie) Insert(word string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.Root == nil {
		t.Root = &trieNode{children: make(map[rune]*trieNode)}
	}

	curr := t.Root

	for _, r := range word {
		if _, ok := curr.children[r]; !ok {
			curr.children[r] = &trieNode{children: make(map[rune]*trieNode)}
		}
		curr = curr.children[r]
	}

	curr.isEnd = true
}

func (t *trie) Autocomplete(prefix string) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var results []string

	curr := t.Root

	// loop through prefix and find the last node of the prefix.
	for _, r := range prefix {
		// return empty results if we encounter a letter not in the prefix path in the trie.
		if _, ok := curr.children[r]; !ok {
			return results
		}
		curr = curr.children[r]
	}

	// Need to search on the last node to find all children.
	t.findAllChildren(curr, prefix, &results)

	return results
}

// This is also known as dfs.
func (t *trie) findAllChildren(node *trieNode, prefix string, results *[]string) {
	// if node is end we need to make sure to update results with the
	// prefix which is the full word.
	if node.isEnd {
		*results = append(*results, prefix)
		return
	}

	for r, child := range node.children {
		// since we're going to have to search through all the child's children
		// and all their children might as well just call ourselves with the child node.
		t.findAllChildren(child, prefix+string(r), results)
	}
}

func (t *trie) Contains(word string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	curr := t.Root

	for _, r := range word {
		if _, ok := curr.children[r]; !ok {
			// we don't have that character in this chain.
			return false
		}
		// move to the next child node
		curr = curr.children[r]
	}
	// is this node marked as the end? If not technically the word doesn't exist.
	return curr.isEnd
}

func (t *trie) ListContents() []string {
	var results []string

	if t.Root == nil {
		return results
	}

	curr := t.Root
	for r, child := range curr.children {
		t.findAllChildren(child, string(r), &results)
	}

	return results
}

// Make the root empty, removing all references to the old data.
func (t *trie) Clear() {
	t.Root = &trieNode{children: make(map[rune]*trieNode)}
}

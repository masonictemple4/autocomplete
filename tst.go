package autocomplete

import "sync"

var _ autocompleter = (*ternarysearchtree)(nil)

type tstNode struct {
	Char             rune
	Left, Mid, Right *tstNode
	IsEnd            bool
}

type ternarysearchtree struct {
	Root *tstNode

	mu sync.RWMutex
}

func nNewTSTNode(char rune) *tstNode {
	return &tstNode{Char: char, IsEnd: false}
}

func newTernarySearchTree(word string) *ternarysearchtree {
	tst := &ternarysearchtree{}
	if word == "" {
		return tst
	}

	tst.Insert(word)
	return tst
}

func (t *ternarysearchtree) Insert(word string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Root = t.insert(t.Root, word, 0)
}

func (t *ternarysearchtree) insert(node *tstNode, word string, index int) *tstNode {
	char := rune(word[index])

	if node == nil {
		node = nNewTSTNode(char)
	}

	if char < node.Char {
		node.Left = t.insert(node.Left, word, index)
	} else if char > node.Char {
		node.Right = t.insert(node.Right, word, index)
	} else if index < len(word)-1 {
		// if the char is equal/not less than or greater than node char
		// we know we're in the mid, now we need to make sure that we still have
		// characters left in the word. So we set mid, and increment the index
		node.Mid = t.insert(node.Mid, word, index+1)
	} else {
		node.IsEnd = true
	}

	return node
}

func (t *ternarysearchtree) Contains(word string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	node := t.contains(t.Root, word, 0)
	return node != nil && node.IsEnd
}

func (t *ternarysearchtree) contains(node *tstNode, word string, index int) *tstNode {
	char := rune(word[index])

	if node == nil {
		return nil
	}

	if char < node.Char {
		return t.contains(node.Left, word, index)
	} else if char > node.Char {
		return t.contains(node.Right, word, index)
	} else if index < len(word)-1 {
		return t.contains(node.Mid, word, index+1)
	} else {
		return node
	}

}

func (t *ternarysearchtree) Autocomplete(prefix string) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var results []string
	node := t.getPrefixNode(t.Root, prefix, 0)
	if node == nil {
		return results
	}

	// middle node continues a word. So we know that every
	// word in the subtree of the middle child of this node
	// is a valid completion of the prefix.
	t.collect(node.Mid, prefix, &results)

	return results
}

func (t *ternarysearchtree) getPrefixNode(node *tstNode, prefix string, index int) *tstNode {
	// recursive so make sure to check first
	if node == nil {
		return nil
	}

	char := rune(prefix[index])

	if char < node.Char {
		return t.getPrefixNode(node.Left, prefix, index)
	} else if char > node.Char {
		return t.getPrefixNode(node.Right, prefix, index)
	} else if index < len(prefix)-1 {
		return t.getPrefixNode(node.Mid, prefix, index+1)
	} else {
		return node
	}
}

// dfs, also in order traversal (left, parent, middle, right)
func (t *ternarysearchtree) collect(node *tstNode, prefix string, results *[]string) {
	// recursive so return early.
	if node == nil {
		return
	}

	t.collect(node.Left, prefix, results)
	if node.IsEnd {
		*results = append(*results, prefix+string(node.Char))
	}
	t.collect(node.Mid, prefix+string(node.Char), results)
	t.collect(node.Right, prefix, results)

}

func (t *ternarysearchtree) ListContents() []string {
	var results []string

	t.collect(t.Root, "", &results)

	return results
}

// Make the root empty, removing all references to the old data.
func (t *ternarysearchtree) Clear() {
	t.Root = &tstNode{}
}

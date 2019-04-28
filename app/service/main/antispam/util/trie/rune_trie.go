package trie

func (trie *RuneTrie) find(contents []rune, accu []rune) (string, interface{}) {
	if trie.value != nil {
		return string(accu), trie.value
	}
	if len(contents) == 0 {
		return "", nil
	}
	k := contents[0]
	tt, ok := trie.children[k]
	if !ok {
		return "", nil
	}
	return tt.find(contents[1:], append(accu, k))
}

func (trie *RuneTrie) Find(content string, accu string) (string, interface{}) {
	rcontent := []rune(content)
	if len(rcontent) == 0 {
		return "", nil
	}
	k, v := trie.find(rcontent, []rune(accu))
	if v == nil {
		return trie.Find(string(rcontent[1:]), accu)
	}
	return k, v
}

/*
* the belowing codes come from "https://github.com/dghubble/trie"
 */
type RuneTrie struct {
	value    interface{}
	children map[rune]*RuneTrie
}

func NewRuneTrie() *RuneTrie {
	return &RuneTrie{
		children: make(map[rune]*RuneTrie),
	}
}

func (trie *RuneTrie) Get(key string) interface{} {
	node := trie
	for _, r := range key {
		node = node.children[r]
		if node == nil {
			return nil
		}
	}
	return node.value
}

func (trie *RuneTrie) Put(key string, value interface{}) bool {
	node := trie
	for _, r := range key {
		child := node.children[r]
		if child == nil {
			child = NewRuneTrie()
			node.children[r] = child
		}
		node = child
	}
	isNewVal := node.value == nil
	node.value = value
	return isNewVal
}

func (trie *RuneTrie) Delete(key string) bool {
	path := make([]nodeRune, len(key))
	node := trie
	for i, r := range key {
		path[i] = nodeRune{r: r, node: node}
		node = node.children[r]
		if node == nil {
			return false
		}
	}
	node.value = nil
	if node.isLeaf() {
		for i := len(key) - 1; i >= 0; i-- {
			parent := path[i].node
			r := path[i].r
			delete(parent.children, r)
			if parent.value != nil || !parent.isLeaf() {
				break
			}
		}
	}
	return true
}

func (trie *RuneTrie) Walk(walker WalkFunc) error {
	return trie.walk("", walker)
}

type nodeRune struct {
	node *RuneTrie
	r    rune
}

func (trie *RuneTrie) walk(key string, walker WalkFunc) error {
	if trie.value != nil {
		walker(key, trie.value)
	}
	for r, child := range trie.children {
		err := child.walk(key+string(r), walker)
		if err != nil {
			return err
		}
	}
	return nil
}

func (trie *RuneTrie) isLeaf() bool {
	return len(trie.children) == 0
}

type WalkFunc func(key string, value interface{}) error

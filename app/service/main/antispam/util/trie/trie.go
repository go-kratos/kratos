package trie

type Trier interface {
	Put(key string, value interface{}) bool
	Find(in string, acc string) (key string, value interface{})
	Get(key string) interface{}
}

func New() Trier {
	return NewRuneTrie()
}

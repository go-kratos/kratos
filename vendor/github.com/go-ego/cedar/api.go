package cedar

// Status reports the following statistics of the cedar:
//	keys:		number of keys that are in the cedar,
//	nodes:		number of trie nodes (slots in the base array) has been taken,
//	size:			the size of the base array used by the cedar,
//	capacity:		the capicity of the base array used by the cedar.
func (da *Cedar) Status() (keys, nodes, size, capacity int) {
	for i := 0; i < da.Size; i++ {
		n := da.Array[i]
		if n.Check >= 0 {
			nodes++
			if n.Value >= 0 {
				keys++
			}
		}
	}
	return keys, nodes, da.Size, da.Capacity
}

// Jump travels from a node `from` to another node
// `to` by following the path `path`.
// For example, if the following keys were inserted:
//	id	key
//	19	abc
//	23	ab
//	37	abcd
// then
//	Jump([]byte("ab"), 0) = 23, nil		// reach "ab" from root
//	Jump([]byte("c"), 23) = 19, nil			// reach "abc" from "ab"
//	Jump([]byte("cd"), 23) = 37, nil		// reach "abcd" from "ab"
func (da *Cedar) Jump(path []byte, from int) (to int, err error) {
	for _, b := range path {
		if da.Array[from].Value >= 0 {
			return from, ErrNoPath
		}
		to = da.Array[from].base() ^ int(b)
		if da.Array[to].Check != from {
			return from, ErrNoPath
		}
		from = to
	}
	return to, nil
}

// Key returns the key of the node with the given `id`.
// It will return ErrNoPath, if the node does not exist.
func (da *Cedar) Key(id int) (key []byte, err error) {
	for id > 0 {
		from := da.Array[id].Check
		if from < 0 {
			return nil, ErrNoPath
		}
		if char := byte(da.Array[from].base() ^ id); char != 0 {
			key = append(key, char)
		}
		id = from
	}
	if id != 0 || len(key) == 0 {
		return nil, ErrInvalidKey
	}
	for i := 0; i < len(key)/2; i++ {
		key[i], key[len(key)-i-1] = key[len(key)-i-1], key[i]
	}
	return key, nil
}

// Value returns the value of the node with the given `id`.
// It will return ErrNoValue, if the node does not have a value.
func (da *Cedar) Value(id int) (value int, err error) {
	value = da.Array[id].Value
	if value >= 0 {
		return value, nil
	}
	to := da.Array[id].base()
	if da.Array[to].Check == id && da.Array[to].Value >= 0 {
		return da.Array[to].Value, nil
	}
	return 0, ErrNoValue
}

// Insert adds a key-value pair into the cedar.
// It will return ErrInvalidValue, if value < 0 or >= ValueLimit.
func (da *Cedar) Insert(key []byte, value int) error {
	if value < 0 || value >= ValueLimit {
		return ErrInvalidValue
	}
	p := da.get(key, 0, 0)
	*p = value
	return nil
}

// Update increases the value associated with the `key`.
// The `key` will be inserted if it is not in the cedar.
// It will return ErrInvalidValue, if the updated value < 0 or >= ValueLimit.
func (da *Cedar) Update(key []byte, value int) error {
	p := da.get(key, 0, 0)

	// key was not inserted
	if *p == ValueLimit {
		*p = value
		return nil
	}

	// key was inserted before
	if *p+value < 0 || *p+value >= ValueLimit {
		return ErrInvalidValue
	}
	*p += value
	return nil
}

// Delete removes a key-value pair from the cedar.
// It will return ErrNoPath, if the key has not been added.
func (da *Cedar) Delete(key []byte) error {
	// if the path does not exist, or the end is not a leaf,
	// nothing to delete
	to, err := da.Jump(key, 0)
	if err != nil {
		return ErrNoPath
	}

	if da.Array[to].Value < 0 {
		base := da.Array[to].base()
		if da.Array[base].Check == to {
			to = base
		}
	}

	for to > 0 {
		from := da.Array[to].Check
		base := da.Array[from].base()
		label := byte(to ^ base)

		// if `to` has sibling, remove `to` from the sibling list, then stop
		if da.Ninfos[to].Sibling != 0 || da.Ninfos[from].Child != label {
			// delete the label from the child ring first
			da.popSibling(from, base, label)
			// then release the current node `to` to the empty node ring
			da.pushEnode(to)
			break
		}
		// otherwise, just release the current node `to` to the empty node ring
		da.pushEnode(to)
		// then check its parent node
		to = from
	}
	return nil
}

// Get returns the value associated with the given `key`.
// It is equivalent to
//		id, err1 = Jump(key)
//		value, err2 = Value(id)
// Thus, it may return ErrNoPath or ErrNoValue,
func (da *Cedar) Get(key []byte) (value int, err error) {
	to, err := da.Jump(key, 0)
	if err != nil {
		return 0, err
	}
	return da.Value(to)
}

// PrefixMatch returns a list of at most `num` nodes
// which match the prefix of the key.
// If `num` is 0, it returns all matches.
// For example, if the following keys were inserted:
//	id	key
//	19	abc
//	23	ab
//	37	abcd
// then
//	PrefixMatch([]byte("abc"), 1) = [ 23 ]				// match ["ab"]
//	PrefixMatch([]byte("abcd"), 0) = [ 23, 19, 37]
// match ["ab", "abc", "abcd"]
func (da *Cedar) PrefixMatch(key []byte, num int) (ids []int) {
	for from, i := 0, 0; i < len(key); i++ {
		to, err := da.Jump(key[i:i+1], from)
		if err != nil {
			break
		}
		if _, err := da.Value(to); err == nil {
			ids = append(ids, to)
			num--
			if num == 0 {
				return
			}
		}
		from = to
	}
	return
}

// PrefixPredict returns a list of at most `num` nodes
// which has the key as their prefix.
// These nodes are ordered by their keys.
// If `num` is 0, it returns all matches.
// For example, if the following keys were inserted:
//	id	key
//	19	abc
//	23	ab
//	37	abcd
// then
//	PrefixPredict([]byte("ab"), 2) = [ 23, 19 ]			// predict ["ab", "abc"]
//	PrefixPredict([]byte("ab"), 0) = [ 23, 19, 37 ]
// predict ["ab", "abc", "abcd"]
func (da *Cedar) PrefixPredict(key []byte, num int) (ids []int) {
	root, err := da.Jump(key, 0)
	if err != nil {
		return
	}
	for from, err := da.begin(root); err == nil; from, err = da.next(from, root) {
		ids = append(ids, from)
		num--
		if num == 0 {
			return
		}
	}
	return
}

func (da *Cedar) begin(from int) (to int, err error) {
	for c := da.Ninfos[from].Child; c != 0; {
		to = da.Array[from].base() ^ int(c)
		c = da.Ninfos[to].Child
		from = to
	}
	if da.Array[from].base() > 0 {
		return da.Array[from].base(), nil
	}
	return from, nil
}

func (da *Cedar) next(from int, root int) (to int, err error) {
	c := da.Ninfos[from].Sibling
	for c == 0 && from != root && da.Array[from].Check >= 0 {
		from = da.Array[from].Check
		c = da.Ninfos[from].Sibling
	}
	if from == root {
		return 0, ErrNoPath
	}
	from = da.Array[da.Array[from].Check].base() ^ int(c)
	return da.begin(from)
}

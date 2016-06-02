/*

Symbol table implementation.

*/

package hufio

import (
	"github.com/icza/huffman"
	"sort"
)

const (
	newValue huffman.ValueType = 1<<31 - 1 - iota // Value representing a new value
	eofValue                                      // Value representing end of data
)

// symbols manages the symbol table and their frequencies.
type symbols struct {
	// leaves of the Huffman tree, symbols previously encountered.
	// Kept sorted by Node.Count, descendant (so new nodes can simply be appended)!
	leaves []*huffman.Node

	root *huffman.Node // Root of the Huffman tree.

	valueMap map[huffman.ValueType]*huffman.Node // Map from value to Node
}

// newSymbols creates a new symbols.
func newSymbols() *symbols {
	// initial leaves: 2 nodes (newValue and eofValue) with count=1, and a high capacity
	leaves := make([]*huffman.Node, 2, 258)
	leaves[0] = &huffman.Node{Value: newValue, Count: 1}
	leaves[1] = &huffman.Node{Value: eofValue, Count: 1}

	valueMap := make(map[huffman.ValueType]*huffman.Node, cap(leaves))
	for _, v := range leaves {
		valueMap[v.Value] = v
	}

	s := &symbols{leaves: leaves, valueMap: valueMap}
	// Reader needs the Huffman tree ready right away, so build it:
	s.rebuildTree()

	return s
}

// update updates the symbol table by incrementing the occurrence count of the specified Node.
func (s *symbols) update(node *huffman.Node) {
	// We have to keep leaves sorted!
	// So first find node in the leaves slice using binary search
	// (remember: leaves is sorted by Node.Count descendant)
	ls, count := s.leaves, node.Count
	idx := sort.Search(len(ls), func(i int) bool { return ls[i].Count <= count })
	idx2 := idx // Store it for later use
	// idx points to the first (lowest) Node having count.
	// There might be more nodes with the same count, find our node:
	for ; ls[idx] != node; idx++ {
	}
	// If there are more nodes with the same count, our node must be switched
	// with the node having same count and the lowest index
	// (so slice remains sorted after incrementing the count of our node).
	if idx2 != idx {
		ls[idx2], ls[idx] = ls[idx], ls[idx2]
	}

	node.Count++

	s.rebuildTree()
}

// insert inserts an encountered new symbol.
func (s *symbols) insert(symbol huffman.ValueType) {
	node := &huffman.Node{Value: symbol, Count: 1}
	// leaves is sorted descendant, so we could simply append.
	// But last are the newValue and eofValue which never increase, so we insert before them:
	ls := s.leaves
	ls = append(ls, ls[len(ls)-1])
	ls[len(ls)-2] = ls[len(ls)-3]
	ls[len(ls)-3] = node
	s.leaves = ls

	s.valueMap[node.Value] = node

	s.rebuildTree()
}

// rebuildTree rebuilds the Huffman tree.
func (s *symbols) rebuildTree() {
	// huffman.BuildSorted() modifies the slice, so make a copy:
	// leaves is sorted descendant, so fill backward:
	j := len(s.leaves)
	ls := make([]*huffman.Node, j)
	for _, v := range s.leaves {
		j--
		ls[j] = v
	}

	// Note: Writer doesn't need the returned root, as codes are read out by stepping through Node.Parent
	// but Reader needs it.
	s.root = huffman.BuildSorted(ls)
}

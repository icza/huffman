/*

Symbol table implementation.

*/

package hufio

import (
	"github.com/icza/huffman"
	"sort"
)

const (
	newValue    huffman.ValueType   = 1<<31 - 1 - iota // Value representing a new value
	eofValue                                           // Value representing end of data
	extraValues = iota                                 // Number of extra, custom values
	maxValues   = 256 + extraValues                    // Max values: number of bytes + extra values
)

// win is a sliding window buffer, the base of the symbol table.
type win struct {
	buf    []huffman.ValueType // Content of the window buffer
	pos    int                 // Position in the buffer
	filled bool                // Tells if the buffer has been filled
}

// store stores the next symbol, and slides the window if it is already filled.
func (w *win) store(symbol huffman.ValueType) {
	w.buf[w.pos] = symbol
	w.pos++
	if w.pos == len(w.buf) {
		w.pos = 0
		if !w.filled {
			w.filled = true
		}
	}
}

// symbols manages the symbol table and their frequencies.
type symbols struct {
	// leaves of the Huffman tree, symbols previously encountered.
	// Kept sorted by Node.Count, descendant (so new nodes can simply be appended)!
	leaves []*huffman.Node

	buffer []*huffman.Node // Reusable buffer to pass when building the Huffman tree

	root *huffman.Node // Root of the Huffman tree.

	valueMap map[huffman.ValueType]*huffman.Node // Map from value to Node

	win *win // The window buffer, nil if no window buffer is used
}

// newSymbols creates a new symbols.
func newSymbols(o *Options) *symbols {
	// initial leaves: 2 nodes (newValue and eofValue) with count=1, and a high capacity
	leaves := make([]*huffman.Node, extraValues, maxValues)
	leaves[0] = &huffman.Node{Value: newValue, Count: 1}
	leaves[1] = &huffman.Node{Value: eofValue, Count: 1}

	valueMap := make(map[huffman.ValueType]*huffman.Node, cap(leaves))
	for _, v := range leaves {
		valueMap[v.Value] = v
	}

	s := &symbols{leaves: leaves, valueMap: valueMap, buffer: make([]*huffman.Node, 0, maxValues)}
	if o.WinSize > 0 {
		s.win = &win{buf: make([]huffman.ValueType, o.WinSize)}
	}

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
	idx := sort.Search(len(ls)-extraValues, func(i int) bool { return ls[i].Count <= count })
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

	s.updateWin(node.Value)

	s.rebuildTree()
}

// updateWin updates the window: slides it if it is already filled, and stores the currently handled symbol.
func (s *symbols) updateWin(symbol huffman.ValueType) {
	if s.win == nil {
		return
	}

	if s.win.filled {
		// Handle symbol shifting out of the window buffer:
		node := s.valueMap[s.win.buf[s.win.pos]]
		// We have to keep leaves sorted!
		// So first find node in the leaves slice using binary search
		// (remember: leaves is sorted by Node.Count descendant)
		ls, count := s.leaves, node.Count
		idx := sort.Search(len(ls)-extraValues, func(i int) bool { return ls[i].Count <= count })
		// idx points to the first (lowest) Node having count.
		// There might be more nodes with the same count, find our node:
		for ; ls[idx] != node; idx++ {
		}

		if count > 1 {
			// If there are more nodes with the same count, our node must be switched
			// with the node having same count and the highest index
			// (so slice remains sorted after decrementing the count of our node).
			idx2 := idx + 1
			for ; idx2 < len(ls)-extraValues && ls[idx2].Count == count; idx2++ {
			}
			if idx2 = idx2 - 1; idx2 != idx {
				ls[idx2], ls[idx] = ls[idx], ls[idx2]
			}
			node.Count--
		} else {
			// Count will decrease to zero: remove node
			// Note: this is optional. If we would not remove it, algorithm would still work as-is.
			// Consider putting this to an option.
			s.leaves = append(ls[:idx], ls[idx+1:]...)
			// Also remove from valueMap:
			delete(s.valueMap, node.Value)
		}
	}

	s.win.store(symbol)
}

// insert inserts an encountered new symbol.
func (s *symbols) insert(symbol huffman.ValueType) {
	node := &huffman.Node{Value: symbol, Count: 1}
	// leaves is sorted descendant, so we could simply append.
	// But extra values are at the end never increase, so we insert before them:
	ls := s.leaves
	// extend by 1 for the new node
	ls = ls[:len(ls)+1]
	// Copy extra values to the end (higher by 1)
	copy(ls[len(ls)-extraValues:], ls[len(ls)-extraValues-1:])
	// And insert the new node
	ls[len(ls)-extraValues-1] = node
	s.leaves = ls

	s.valueMap[node.Value] = node

	s.updateWin(node.Value)

	s.rebuildTree()
}

// rebuildTree rebuilds the Huffman tree.
func (s *symbols) rebuildTree() {
	// huffman.BuildSorted() modifies the slice, so make a copy:
	// leaves is sorted descendant, so fill backward:
	j := len(s.leaves)
	ls := s.buffer[:j]
	for _, v := range s.leaves {
		j--
		ls[j] = v
	}

	// Note: Writer doesn't need the returned root, as codes are read out by stepping through Node.Parent
	// but Reader needs it.
	s.root = huffman.BuildSorted(ls)
}

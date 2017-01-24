/*

Options for creating Huffman Readers and Writers.

*/

package hufio

// Options wraps options for creating Huffman Readers and Writers.
// Zero value for a field means to use the default value for that field.
type Options struct {

	// WinSize specifies the size of the sliding window that is used to manage
	// the symbol table.
	// 0 means to use a sliding window with the default size (2048 bytes / symbols).
	// Negative values mean not to use a sliding window, that is, symbol table is
	// calculated based on all previously encountered symbols.
	WinSize int
}

// checkOptions returns a new Options where "missing" fields (with zero value) are set to default values.
// The passed options is not modified.
// It is allowed to pass nil, which is treated as the zero value of Options.
func checkOptions(o *Options) *Options {
	o2 := new(Options)
	if o != nil {
		*o2 = *o
	}

	if o2.WinSize == 0 {
		o2.WinSize = 2048
	}

	return o2
}

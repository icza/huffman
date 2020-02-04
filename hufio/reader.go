/*

Huffman code Reader implementation.

*/

package hufio

import (
	"io"

	"github.com/icza/bitio"
	"github.com/icza/huffman"
)

// Reader is the Huffman reader implementation.
// It also implements io.ByteReader.
type Reader struct {
	*symbols
	br *bitio.Reader
}

// NewReader returns a new Reader using the specified io.Reader as the input (source),
// with the default Options.
func NewReader(in io.Reader) *Reader {
	return NewReaderOptions(in, nil)
}

// NewReaderOptions returns a new Reader using the specified io.Reader as the input (source)
// with the specified Options.
//
// Note: Options are not transmitted internally! The Reader will only be able to properly decode the stream
// created by a Writer if the same Options is used both at the Reader and Writer.
// Transmitting the Options has to be done manually if needed.
func NewReaderOptions(in io.Reader, o *Options) *Reader {
	o = checkOptions(o)
	return &Reader{symbols: newSymbols(o), br: bitio.NewReader(in)}
}

// Read decompresses up to len(p) bytes from the source.
func (r *Reader) Read(p []byte) (n int, err error) {
	for i := range p {
		if p[i], err = r.ReadByte(); err != nil {
			return i, err
		}
	}
	return len(p), nil
}

// ReadByte decompresses a single byte.
func (r *Reader) ReadByte() (b byte, err error) {
	// Read Huffman code
	br := r.br
	node := r.root
	for node.Left != nil { // read until we reach a leaf
		var right bool
		if right, err = br.ReadBool(); err != nil {
			return
		} else if right {
			node = node.Right
		} else {
			node = node.Left
		}
	}

	switch node.Value {
	case newValue:
		if b, err = br.ReadByte(); err != nil {
			return
		}
		r.insert(huffman.ValueType(b))
		return
	case eofValue:
		return 0, io.EOF
	default:
		r.update(node)
		return byte(node.Value), nil
	}
}

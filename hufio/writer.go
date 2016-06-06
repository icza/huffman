/*

Huffman code Writer implementation.

*/

package hufio

import (
	"github.com/icza/bitio"
	"github.com/icza/huffman"
	"io"
)

// Writer is the Huffman writer interface.
// Must be closed in order to properly send EOF.
type Writer interface {
	// Writer is an io.Writer and io.Closer.
	// Close closes the Huffman writer, properly sending EOF.
	// If the underlying io.Writer implements io.Closer,
	// it will be closed after sending EOF.
	// Write writes the compressed form of p to the underlying io.Writer.
	// The compressed bytes are not necessarily flushed until the Writer is closed.
	io.WriteCloser

	// Writer is also an io.ByteWriter.
	io.ByteWriter
}

// writer is the Huffman writer implementation.
type writer struct {
	*symbols
	bw bitio.Writer
}

// NewWriter returns a new Writer using the specified io.Writer as the output,
// with the default Options.
func NewWriter(out io.Writer) Writer {
	return NewWriterOptions(out, nil)
}

// NewWriterOptions returns a new Writer using the specified io.Writer as the output,
// with the specified Options.
//
// Note: Options are not transmitted internally! The Reader will only be able to properly decode the stream
// created by a Writer if the same Options is used both at the Reader and Writer.
// Transmitting the Options has to be done manually if needed.
func NewWriterOptions(out io.Writer, o *Options) Writer {
	o = checkOptions(o)
	return &writer{symbols: newSymbols(o), bw: bitio.NewWriter(out)}
}

// Write implements io.Writer.
func (w *writer) Write(p []byte) (n int, err error) {
	for i, v := range p {
		if err = w.WriteByte(v); err != nil {
			return i, err
		}
	}
	return len(p), nil
}

// WriteByte implements io.ByteWriter.
func (w *writer) WriteByte(b byte) (err error) {
	value := huffman.ValueType(b)
	node := w.valueMap[value]

	if node == nil {
		// New value, write out newValue's Huffman code
		if err = w.bw.WriteBits(w.valueMap[newValue].Code()); err != nil {
			return
		}
		// ...and the new value
		if err = w.bw.WriteByte(b); err != nil {
			return
		}
		w.insert(value)
	} else {
		// Write out node's Huffman code
		if err = w.bw.WriteBits(node.Code()); err != nil {
			return
		}
		w.update(node)
	}
	return
}

// Close implements io.Closer.
func (w *writer) Close() (err error) {
	// If there were any data, write out eofValue
	if len(w.leaves) > 2 {
		// Write out eofValue's Huffman code
		if err = w.bw.WriteBits(w.valueMap[eofValue].Code()); err != nil {
			return
		}
	}
	return w.bw.Close()
}

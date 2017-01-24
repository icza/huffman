# huffman

[![Build Status](https://travis-ci.org/icza/huffman.svg?branch=master)](https://travis-ci.org/icza/huffman)
[![GoDoc](https://godoc.org/github.com/icza/huffman?status.svg)](https://godoc.org/github.com/icza/huffman)
[![Go Report Card](https://goreportcard.com/badge/github.com/icza/huffman)](https://goreportcard.com/report/github.com/icza/huffman)
[![codecov](https://codecov.io/gh/icza/huffman/branch/master/graph/badge.svg)](https://codecov.io/gh/icza/huffman)

[Huffman coding](https://en.wikipedia.org/wiki/Huffman_coding) implementation in Go
(Huffman tree, Symbol table, Huffman Reader + Writer).

### Huffman Tree

Use the `Build()` function to build a Huffman tree. Use the `Print()` function to print Huffman codes
of all leaves of a tree (for verification).

Example:

	leaves := []*Node{
		{Value: ' ', Count: 20},
		{Value: 'a', Count: 40},
		{Value: 'm', Count: 10},
		{Value: 'l', Count: 7},
		{Value: 'f', Count: 8},
		{Value: 't', Count: 15},
	}
	root := Build(leaves)
	Print(root)

Output:

	'a': 0
	'm': 100
	'l': 1010
	'f': 1011
	't': 110
	' ': 111

### Huffman Reader and Writer

[![GoDoc](https://godoc.org/github.com/icza/huffman/hufio?status.svg)](https://godoc.org/github.com/icza/huffman/hufio)

The `hufio` package implements a Huffman `Reader` and `Writer`. You may use these to transmit Huffman code of your data.

This `Reader` and `Writer` internally manages a Symbol Table (the frequency of encountered symbols, updated dynamically).
The `Writer` computes and sends the Huffman code of the data, the `Reader` receives the Huffman code and "reconstructs"
the original data based on that.

The implementation uses a _sliding window_ which is used to manage the symbol table.
The sliding window is optional, that is, if no window is used, the symbol table is calculated based on
all previously encountered symbols.

`Writer` + `Reader` example:

	buf := &bytes.Buffer{}
	w := NewWriter(buf)
	if _, err := w.Write([]byte("Testing Huffman Writer + Reader.")); err != nil {
		log.Panicln("Failed to write:", err)
	}
	if err := w.Close(); err != nil {
		log.Panicln("Failed to close:", err)
	}

	r := NewReader(bytes.NewReader(buf.Bytes()))
	if data, err := ioutil.ReadAll(r); err != nil {
		log.Panicln("Failed to read:", err)
	} else {
		log.Println("Read:", string(data))
	}

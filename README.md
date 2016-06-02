# huffman

[![GoDoc](https://godoc.org/github.com/icza/huffman?status.svg)](https://godoc.org/github.com/icza/huffman)

Huffman coding implementation in Go (Huffman tree, Symbol table, Huffman Reader + Writer).

Use the `Build()` function to build a Huffman tree. Use the `Print()` function to print a Huffman tree (for debugging purposes).

### Huffman Tree

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

The `hufio` package implements a Huffman `Reader` and `Writer`. You may use these to transmit Huffman code of your data.

This `Reader` and `Writer` internally manages a Symbol Table (the frequency of encountered symbols, updated dynamically).
The `Writer` computes and sends the Huffman code of the data, the `Reader` receives the Huffman code and "reconstructs"
the original data based on that.


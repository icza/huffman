/*

Package hufio implements a Huffman Reader and Writer which transmits the Huffman codes of data.

This Reader and Writer internally manages a Symbol Table (the frequency of encountered symbols, updated dynamically).
The Writer computes and sends the Huffman code of the data, the Reader receives the Huffman code and "reconstructs"
the original data based on that.

*/
package hufio

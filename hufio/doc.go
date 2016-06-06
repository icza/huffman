/*

Package hufio implements a Huffman Reader and Writer which transmits the Huffman codes of data.

This Reader and Writer internally manages a Symbol Table (the frequency of encountered symbols, updated dynamically).
The Writer computes and sends the Huffman code of the data, the Reader receives the Huffman code and "reconstructs"
the original data based on that.

The implementation uses a sliding window which is used to manage the symbol table.
The sliding window is optional, that is, if no window is used, the symbol table is calculated based on
all previously encountered symbols.

Writer + Reader example:

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


*/
package hufio

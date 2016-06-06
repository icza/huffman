package hufio

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const dataSize = 9000

func TestOptions(t *testing.T) {
	data := []byte("testing, testing ttttttttttttt")
	cases := []struct {
		name string
		data []byte
		o    *Options
	}{
		{"Options Default", data, &Options{}},
		{"Options nil", data, nil},
		{"Options [WinSize= 3]", data, &Options{WinSize: 3}},
		{"Options [WinSize= 1]", data, &Options{WinSize: 1}},
		{"Options [WinSize=-1]", data, &Options{WinSize: -1}},
	}

	for _, v := range cases {
		testWriteAndRead(v.name, v.data, t, v.o)
	}
}

func TestRandomDigits(t *testing.T) {
	data := make([]byte, dataSize)
	for i := range data {
		data[i] = '0' + byte(rand.Int31n(10))
	}
	testWriteAndRead("Random Digits [0..9]", data, t)
}

func TestRandomLetters(t *testing.T) {
	data := make([]byte, dataSize)
	for i := range data {
		data[i] = 'a' + byte(rand.Int31n(26))
	}
	testWriteAndRead("Random Letters [a..z]", data, t)
}

func TestRandomBytes(t *testing.T) {
	data := make([]byte, dataSize)
	for i := range data {
		data[i] = byte(rand.Int31n(256))
	}
	testWriteAndRead("Random Bytes [0..255]", data, t)
}

func TestFiles(t *testing.T) {
	fnames := []string{
		"wiki_huffman.html_",
		"wiki_huffman.zip",
	}
	for _, fname := range fnames {
		data, err := ioutil.ReadFile("_test_files/" + fname)
		if err != nil {
			t.Error("Can't read input:", err)
		}
		testWriteAndRead(fname, data, t)
	}
}

func testWriteAndRead(testCase string, data []byte, t *testing.T, os ...*Options) {
	buf := &bytes.Buffer{}

	var w Writer
	if len(os) > 0 {
		w = NewWriterOptions(buf, os[0])
	} else {
		w = NewWriter(buf)
	}
	if _, err := w.Write(data); err != nil {
		t.Error("Failed to write:", err)
	}
	if err := w.Close(); err != nil {
		t.Error("Failed to close:", err)
	}
	outs := len(buf.Bytes())
	fmt.Printf("%-22s: Writer Input: %6d bytes, Output: %6d, ratio: %6.2f %%, %.2f bit/symbol\n",
		testCase, len(data), outs, float64(outs)/float64(len(data))*100, float64(outs)*8.0/float64(len(data)))

	var r Reader
	if len(os) > 0 {
		r = NewReaderOptions(bytes.NewReader(buf.Bytes()), os[0])
	} else {
		r = NewReader(bytes.NewReader(buf.Bytes()))
	}
	data2, err := ioutil.ReadAll(r)
	if err != nil {
		t.Error("Failed to read:", err)
	}
	if !bytes.Equal(data, data2) {
		t.Error("Decoded doesn't match original!", len(data2))
	}
}

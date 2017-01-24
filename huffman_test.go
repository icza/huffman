package huffman

import (
	"fmt"
	"testing"
)

func TestBuild(t *testing.T) {
	fmt.Println()
	t1, t2 := Build(nil), Build([]*Node{})
	if t1 != nil || t2 != nil {
		t.Errorf("Got: %v, %v, want: nil, nil", t1, t2)
	}

	leaves := []*Node{ // From "this is an example of a huffman tree"
		{Value: ' ', Count: 7},
		{Value: 'a', Count: 4},
		{Value: 'e', Count: 4},
		{Value: 'f', Count: 3},
		{Value: 'h', Count: 2},
		{Value: 'i', Count: 2},
		{Value: 'm', Count: 2},
		{Value: 'n', Count: 2},
		{Value: 's', Count: 2},
		{Value: 't', Count: 2},
		{Value: 'l', Count: 1},
		{Value: 'o', Count: 1},
		{Value: 'p', Count: 1},
		{Value: 'r', Count: 1},
		{Value: 'u', Count: 1},
		{Value: 'x', Count: 1},
	}
	root := Build(leaves)
	if root == nil {
		t.Errorf("Got: %v, want: not nil", root)
	}
	Print(root)
}

func TestBuild2(t *testing.T) {
	fmt.Println()
	t1, t2 := Build(nil), Build([]*Node{})
	if t1 != nil || t2 != nil {
		t.Errorf("Got: %v, %v, want: nil, nil", t1, t2)
	}

	leaves := []*Node{
		{Value: ' ', Count: 20},
		{Value: 'a', Count: 40},
		{Value: 'l', Count: 7},
		{Value: 'm', Count: 10},
		{Value: 'f', Count: 8},
		{Value: 't', Count: 15},
	}
	root := Build(leaves)
	if root == nil {
		t.Errorf("Got: %v, want: not nil", root)
	}
	Print(root)

	type code struct {
		r    uint64
		bits byte
	}

	expected := map[ValueType]code{
		'a': {0x0, 1},  // 0
		'm': {0x04, 3}, // 100
		'l': {0x0a, 4}, // 1010
		'f': {0x0b, 4}, // 1011
		't': {0x06, 3}, // 110
		' ': {0x07, 3}, // 111
	}
	for _, leave := range leaves {
		if leave.Left != nil {
			continue // Not leaf
		}
		r, bits := leave.Code()
		if exp, got := (code{r, bits}), expected[leave.Value]; exp != got {
			t.Errorf("Got: %v, want: %v (leave: '%c')", got, exp, leave.Value)
		}
	}
}

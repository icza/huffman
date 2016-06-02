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
}

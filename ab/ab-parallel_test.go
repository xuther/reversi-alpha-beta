package ab

import (
	"log"
	"testing"
)

func TestABMulti(t *testing.T) {
	//test tree
	/*
		a -> b,c
		b -> d, e
		d -> 3, 5
		e -> 6, 9
		c -> f, g
		f -> 1, 2
		g -> 0, -1
	*/

	three := testNode{3, []Node{}, 0, "three"}
	five := testNode{5, []Node{}, 0, "five"}
	six := testNode{6, []Node{}, 0, "six"}
	nine := testNode{9, []Node{}, 0, "nine"}
	one := testNode{1, []Node{}, 0, "one"}
	two := testNode{2, []Node{}, 0, "two"}
	zero := testNode{0, []Node{}, 0, "zero"}
	negone := testNode{-1, []Node{}, 0, "negone"}
	d := testNode{0, []Node{three, five}, 1, "d"}
	e := testNode{0, []Node{six, nine}, 1, "e"}
	f := testNode{0, []Node{one, two}, 1, "f"}
	g := testNode{0, []Node{zero, negone}, 1, "g"}
	b := testNode{0, []Node{d, e}, -1, "b"}
	c := testNode{0, []Node{f, g}, -1, "c"}
	a := testNode{0, []Node{b, c}, 1, "a"}

	val, path := StartSearchMulti(a, 3, -100000, 100000)
	log.Printf("Val: %v, Path: %v", val, path)
}

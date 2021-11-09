package main

import (
	"fmt"
	"testing"
)

type my struct {
	name string
}

type you = my

func (m *my) foo() {

}

func TestTypeTrans(t *testing.T) {

	a := my{name: "zeta"}

	b := you{}
	b.foo()
	// b :=

	// fmt.Println(a, you(a))

	a.name = "chow"
	fmt.Println(a, you(a))

	// b.name = "fred"
	// fmt.Println(a, b)

}

package main

import "os"

func main() {
	Foo(2, 4, 8, 16, 32, 64)
}

func Foo(a, b, c, d, e, f int) int {
	g := a + b - c*d/e%f
	return g
}

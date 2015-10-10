package main

import "os"

func main() {
	a := 32
	b := 64
	c := a + b
	d := b - c
	e := c * d
	f := d / e
	g := e % f
	os.Exit(g)
}

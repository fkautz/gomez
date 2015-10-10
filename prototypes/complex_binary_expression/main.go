package main

import "os"

func main() {
	a := 2
	b := 4
	c := 8
	d := 16
	e := 32
	f := 64
	g := a + b - c*d/e%f + 128 - 256*512/1024%2048
	os.Exit(g)
}

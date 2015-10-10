package main

func main() {
	x := 32
	y := 64
	if x > y {
		x = add(x, y)
	}
	z := add(x, y)
	write(z)
}

func add(a, b int) int {
	return a + b
}

func write(a int) int {
	return a
}

package main

func main() {
	x := 32
	y := 64
	if x > y {
		x = add(x, y)
	}
	y, x = swap(x, y)
	z := y / x
	write(z)
}

func add(a, b int) int {
	return a + b
}

func swap(a, b int) (int, int) {
	return b, a
}

func write(a int) int {
	return a
}

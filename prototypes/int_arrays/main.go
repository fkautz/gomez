package main

func main() {
	var foo [7]int
	for i := 0; i < 7; i = i + 1 {
		Foo(foo[i])
	}
}

func Foo(a int) int {
	return a
}

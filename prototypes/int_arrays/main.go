package main

func main() {
	var foo [7]int
	count := 0
	//	for i := 0; i < len(foo); i = i + 1 {
	//		count = foo[i]
	//	}
	for i := 0; i < 7; i = i + 1 {
		Foo(foo[i])
	}
	//	len(count)
}

func Foo(i int) int {
	return i
}

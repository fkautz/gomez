package main

func main() {
	Foo(10, 20)
}

func Foo(sum int, sum2 int) int {
	for i := 0; i < sum2; i = i + 1 {
		sum = sum + sum2
	}
	sum = Add(sum)
	return sum
}

func Add(a int) int {
	return a + a
}

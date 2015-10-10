default: main

main: optimized_main.o
	gcc -o main optimized_main.o

optimized_main.ll: main.ll
	opt -O3 -S main.ll > optimized_main.ll

optimized_main.o: optimized_main.ll
	llc -filetype=obj optimized_main.ll

main.ll: main.go
	tigerbat compile -i main.go -o main.ll

clean: 
	rm -rf main *.o *.ll

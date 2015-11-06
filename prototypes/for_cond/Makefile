default: main

main.ll: main.go
	gomez compile -i main.go -o main.ll

optimized_main.ll: main.ll
	opt -O3 -S main.ll > optimized_main.ll

optimized_main.o: optimized_main.ll
	llc -filetype=obj optimized_main.ll

main: optimized_main.o
	gcc -o main optimized_main.o

clean: 
	rm -rf main *.o *.ll

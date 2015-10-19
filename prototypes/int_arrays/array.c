#include <stdio.h>
int main() {
  int foo[10];
  foo[0] = 5280;
  printf("%d", foo[7]);
}

int Foo(int foo[10], int index) {
  return foo[index];
}

// File to test extracting globals and external references.
#include <stdio.h>

extern int i;

void foo() {
}
static void bar() {
}
extern int baz();

int i = 7;
extern int j;

int main(int argc, char **argv) {
    int k;
    printf("hello\n");
    foo();
    baz(); // external
    bar();
    baz(); // external
    k = i;
    k = j; // external
    k = j; // external
}

# Fibonacci: compute fib(9) = 55
# x1 = fib(n-1), x2 = fib(n-2), x3 = temp, x4 = counter, x5 = limit

addi x1, x0, 0      # fib(n-1) = 1
addi x2, x0, 1      # fib(n-2) = 0
addi x4, x0, 0      # counter = 0
addi x5, x0, 9      # limit = 9 - looping from 0

loop:
    add x3, x1, x2   # temp = fib(n-2) + fib(n-1)
    add x1, x2, x0   # fib(n-2) = fib(n-1)
    add x2, x3, x0   # fib(n-1) = temp
    addi x4, x4, 1   # counter++
    blt x4, x5, loop # if counter < limit, loop
    addi x1, x2, 0   # move result to x1

# result in x1 = 55 is what I am after

# Should be something like this:
# loops: 1, 2, 3, 4, 5, 6, 7, 8, 9
#  0, 1, 1, 2, 3, 5, 8,13,21,34,55

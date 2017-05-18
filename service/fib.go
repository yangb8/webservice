package service

var cache map[int]int = map[int]int{}

func fib(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}

	if cache[n] == 0 {
		cache[n] = fib(n-2) + fib(n-1)
	}
	return cache[n]
}

package main

import (
	"fmt"
	"math/rand"
)

func count_trailing_zeros(n uint32) uint32 {
	if n == 0 {
		return 32 // all zeros
	} else {
		n = (n ^ (n - 1)) >> 1
		var c uint32
		for c = 0; n > 0; c++ {
			n >>= 1
		}
		return c
	}
}

func loglog(numbers []uint32) int {
	trailing_zeros := make([]uint32, len(numbers))

	for i, n := range numbers {
		trailing_zeros[i] = count_trailing_zeros(n)
	}

	var max uint32
	for _, n := range trailing_zeros {
		if n > max {
			max = n
		}
	}

	return 1 << (max + 1)
}

func hloglog(numbers []uint32) int {
	const num_buckets = 32
    var buckets [num_buckets]uint32

	for _, n := range numbers {
		// 5 bits  -> use to pick bucket
		bi := n >> 27 // bucket index

		// 27 bits -> count trailing zeros
		v := count_trailing_zeros(n)
		if v > 27 {
			v = 27 // don't count leading five bits
		}
		v += 1

		// put it in the bucket
		if buckets[bi] < v {
			buckets[bi] = v
		}

	}

	// combine buckets to make a guess
	sum := 0.0
	for _, b := range buckets {
		sum += 1.0 / float64(int(1) << b)
	}

	return int((0.697 * 32 * 32)/float64(sum) + 0.5)
}

func run_test(i int) float64 {
		numbers := make([]uint32, i)
		distinct_numbers := make(map[uint32]bool)
		for j := range numbers {
			numbers[j] = rand.Uint32()
			distinct_numbers[numbers[j]] = true
		}
		correct := len(distinct_numbers)
		loglog_guess := loglog(numbers)
		hloglog_guess := hloglog(numbers)
		fmt.Printf("%d\t\t%d\t\t%d\t\t%d\t\t%0.2f\t\t%0.2f\n", i, loglog_guess, hloglog_guess,
			       correct, float64(loglog_guess-correct)/float64(correct),
			       float64(hloglog_guess-correct)/float64(correct))

		return float64(hloglog_guess-correct)/float64(correct)

}

func main() {
	rand.Seed(0)

	fmt.Println("Count\t\tLogLog\t\tHLogLog\t\tCorrect\t\tLLError\t\tHLLError")
	for i := 100; i < 10000000; i *= 10 {
		sum := 0.0
		for j := 0; j < 10; j++ {
			sum += run_test(i)
		}
		fmt.Println("\n\t\tAverage HLL Error:", sum/10)
	}

}

package main

import "fmt"

var nums = []float64{29.895, 32.438, 36.416, 46.653, 50.042, 52.722, 54.625, 56.152, 57.973}

func main() {
	var total = 0.00
	numLen := len(nums)

	for i := 1; i < numLen; i++ {
		x := nums[i] - nums[i-1]
		total += x
	}
	fmt.Println("Target: 28.078")
	fmt.Println("Have: ", total)
}

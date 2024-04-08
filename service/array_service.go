package service

import (
	"context"
	"fmt"
)

type ArrayService interface {
	ReverseArray(ctx context.Context, requestId string) (arr []int)
	RotationArray(ctx context.Context, requestId string) (arr []int)
	RearrangeArray(ctx context.Context, requestId string) (arr []int)
	RangeSumArray(ctx context.Context, requestId string) (arr []int)
	RangeWithUpdateArray(ctx context.Context, requestId string) (arr []int)
	// sparsetable pending
	MetricArray1(ctx context.Context, requestId string) (arr []int)
	MetricArray2(ctx context.Context, requestId string) (arr []int)
	MultiplyMatrix(ctx context.Context, requestId string) (arr [][]int)
	KadanesAlgorithm(ctx context.Context, requestId string) (result int)
	DutchNationalFlagAlgorithm(ctx context.Context, requestId string) (arr []int)
}

type ArrayServiceImplementation struct {
}

func NewArrayService() ArrayService {
	return &ArrayServiceImplementation{}
}

func (service *ArrayServiceImplementation) ReverseArray(ctx context.Context, requestId string) (arr []int) {
	arr = []int{5, 6, 7, 4, 8, 9, 3, 8}
	fmt.Println(arr)
	middle := (len(arr) / 2) - 1
	for i := 0; i < len(arr); i++ {
		right := arr[len(arr)-i-1]
		arr[len(arr)-i-1] = arr[i]
		arr[i] = right
		if middle == i {
			break
		}
	}
	fmt.Println(arr)
	return
}

func (service *ArrayServiceImplementation) RotationArray(ctx context.Context, requestId string) (arr []int) {
	arr = []int{1, 2, 3, 4, 5}
	fmt.Println(arr)
	lastElement := arr[len(arr)-1]
	for i := len(arr) - 1; i > 0; i-- {
		arr[i] = arr[i-1]
	}
	arr[0] = lastElement
	fmt.Println(arr)
	return
}

func (service *ArrayServiceImplementation) RearrangeArray(ctx context.Context, requestId string) (arr []int) {
	arr = []int{-1, -1, 6, 1, 9, 3, 2, -1, 4, -1}
	for i := 0; i < len(arr); i++ {
		if arr[i] >= 0 && arr[i] != i {
			right := arr[arr[i]]
			arr[arr[i]] = arr[i]
			arr[i] = right
		}
	}
	fmt.Println(arr)
	return
}

func (service *ArrayServiceImplementation) RangeSumArray(ctx context.Context, requestId string) (arr []int) {
	arr = []int{1, 2, 3, 4, 5}
	i := 1
	j := 3
	sum := 0
	for k := i; k <= j; k++ {
		sum += arr[k]
	}
	fmt.Println(sum)
	return
}

func (service *ArrayServiceImplementation) RangeWithUpdateArray(ctx context.Context, requestId string) (arr []int) {
	arr = []int{1, 2, 3, 4, 5}
	i := 1
	j := 3
	addNumber := 10
	for k := i; k <= j; k++ {
		arr[k] += addNumber
	}
	fmt.Println(arr)
	return
}

func (service *ArrayServiceImplementation) MetricArray1(ctx context.Context, requestId string) (arr []int) {
	arrgiven := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	for i := 0; i < len(arrgiven); i++ {
		if i%2 == 0 {
			for j := 0; j < len(arrgiven[i]); j++ {
				arr = append(arr, arrgiven[i][j])
			}
		} else {
			for j := len(arrgiven[i]) - 1; j >= 0; j-- {
				arr = append(arr, arrgiven[i][j])
			}
		}
	}
	return
}

// https://www.geeksforgeeks.org/print-matrix-in-zig-zag-fashion/
func (service *ArrayServiceImplementation) MetricArray2(ctx context.Context, requestId string) (arr []int) {
	arrgiven := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	rows := len(arrgiven)
	cols := len(arrgiven[0])
	for i := 0; i <= rows+cols-2; i++ {
		if i%2 == 0 {
			// Even sum (downward traversal)
			for j := min(i, rows-1); j >= max(0, i-cols+1); j-- {
				fmt.Printf("%d ", arrgiven[j][i-j])
				arr = append(arr, arrgiven[j][i-j])
			}
		} else {
			// Odd sum (upward traversal)
			for j := max(0, i-cols+1); j <= min(i, rows-1); j++ {
				fmt.Printf("%d ", arrgiven[j][i-j])
				arr = append(arr, arrgiven[j][i-j])
			}
		}
	}
	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// https://www.geeksforgeeks.org/c-program-multiply-two-matrices/
func (service *ArrayServiceImplementation) MultiplyMatrix(ctx context.Context, requestId string) (arr [][]int) {
	matric1 := [][]int{{1, 1}, {2, 2}, {3, 3}}
	fmt.Println("matric1:", matric1)
	matric2 := [][]int{{1, 1, 1}, {2, 2, 2}}
	fmt.Println(matric2)
	rm1 := len(matric1)
	// cm1 := len(matric1[0])
	rm2 := len(matric2)
	cm2 := len(matric2[0])
	for i := 0; i < rm1; i++ {
		var arr2 []int
		for j := 0; j < cm2; j++ {
			var result int
			for k := 0; k < rm2; k++ {
				result += matric1[i][k] * matric2[k][j]
			}
			// fmt.Println("result:", i, j, result)
			arr2 = append(arr2, result)
		}
		arr = append(arr, arr2)
	}
	return
}

// https://www.geeksforgeeks.org/largest-sum-contiguous-subarray/
func (service *ArrayServiceImplementation) KadanesAlgorithm(ctx context.Context, requestId string) (result int) {
	arrgiven := []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}
	var maxCurrent int
	var maxGlobal int

	maxCurrent = arrgiven[0]
	maxGlobal = arrgiven[0]

	for _, value := range arrgiven[1:] {
		maxCurrentValue := maxCurrent + value
		if value > maxCurrentValue {
			maxCurrent = value
		} else {
			maxCurrent = maxCurrentValue
		}

		if maxGlobal < maxCurrent {
			maxGlobal = maxCurrent
		}
	}
	result = maxGlobal
	return
}

// https://www.geeksforgeeks.org/sort-an-array-of-0s-1s-and-2s/
func (service *ArrayServiceImplementation) DutchNationalFlagAlgorithm(ctx context.Context, requestId string) (arr []int) {
	arr = []int{0, 1, 1, 0, 1, 2, 1, 2, 0, 0, 0, 1}
	fmt.Println(arr)
	var low int
	var mid int
	var high int
	low = 0
	mid = 0
	high = len(arr) - 1
	for mid <= high {
		switch arr[mid] {
		case 0:
			temp := arr[mid]
			arr[mid] = arr[low]
			arr[low] = temp
			low++
			mid++
			break
		case 1:
			mid++
			break
		case 2:
			temp := arr[mid]
			arr[mid] = arr[high]
			arr[high] = temp
			high--
			break
		}
	}
	fmt.Println(arr)
	return
}

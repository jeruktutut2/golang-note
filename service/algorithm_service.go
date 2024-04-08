package service

import (
	"context"
	"fmt"
	"math"
)

type AlgorithmService interface {
	LinearSearch(ctx context.Context, requestId string) (arrvalue int)
	BinarySearch(ctx context.Context, requestId string) (arrvalue int)
	InterpolationSearch(ctx context.Context, requestId string) (arrvalue float64)
	JumpSearch(ctx context.Context, requestId string) (arrvalue int)
	TernarySearch(ctx context.Context, requestId string) (arrvalue int)
}

type AlgorithmServiceImplementation struct {
}

func NewAlgorithService() AlgorithmService {
	return &AlgorithmServiceImplementation{}
}

func (service *AlgorithmServiceImplementation) LinearSearch(ctx context.Context, requestId string) (arrvalue int) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
arrLoop:
	for _, value := range arr {
		if 1 == value {
			arrvalue = value
			break arrLoop
		}
	}
	return
}

func (service *AlgorithmServiceImplementation) BinarySearch(ctx context.Context, requestId string) (arrvalue int) {
	arr := []int{1, 2, 3, 4, 5, 7, 9, 11, 15}
	target := 15
	low := 0
	high := len(arr) - 1
	for low <= high {
		median := (low + high) / 2
		if arr[median] < target {
			low = median + 1
		} else if arr[median] > target {
			high = median - 1
		} else {
			arrvalue = arr[median]
			return
		}
	}
	arrvalue = 0
	fmt.Println("not found")
	return

}

// https://www.youtube.com/watch?v=jp3AqTOje5k
// https://www.youtube.com/watch?v=iMVKo1vXVsw&t=18s
func (service *AlgorithmServiceImplementation) InterpolationSearch(ctx context.Context, requestId string) (arrvalue float64) {
	arr := []float64{10, 14, 19, 26, 27, 31, 33, 35, 42, 44}
	key := 31
	low := 0
	high := len(arr) - 1
	for low <= high {
		est := float64(low) + (float64(high-low)/(arr[high]-arr[low]))*(float64(key)-arr[low])
		if arr[int(est)] < float64(key) {
			low = int(est) + 1
		} else if arr[int(est)] > float64(key) {
			high = int(est) - 1
		} else {
			arrvalue = arr[int(est)]
			return
		}
	}
	return
}

// https://www.youtube.com/watch?v=Va2UraOqeHQ&list=PLZBqAUVp79AJ89-2N-TlO8qsbMJbQ9e_S&index=3
func (service *AlgorithmServiceImplementation) JumpSearch(ctx context.Context, requestId string) (arrvalue int) {
	arr := []int{3, 11, 23, 34, 48, 57, 72, 86, 92}
	key := 57
	start := 0
	end := math.Sqrt(float64(len(arr)))
	for arr[int(end)] < key && start < len(arr) {
		start = int(end)
		end = end + math.Sqrt(float64(len(arr)))
		if int(end) > (len(arr) - 1) {
			end = float64(len(arr))
		}
	}
	for i := start; i < int(end); i++ {
		if arr[i] == key {
			arrvalue = arr[i]
			return
		}
	}
	return
}

func (service *AlgorithmServiceImplementation) TernarySearch(ctx context.Context, requestId string) (arrvalue int) {
	arr := []int{7, 12, 19, 24, 31, 48, 51, 63, 67, 73, 81, 92}
	key := 19
	left := 0
	right := len(arr)
	for {
		if left > right {
			return
		}

		middle1 := left + ((right - left) / 3)
		middle2 := right - ((right - left) / 3)

		if arr[middle1] == key {
			arrvalue = arr[middle1]
			return
		} else if arr[middle1] > key {
			right = middle1 - 1
		} else if arr[middle2] == key {
			arrvalue = arr[middle2]
			return
		} else if arr[middle2] > key {
			left = middle1 + 1
			right = middle2 - 1
		} else {
			left = middle2 + 1
		}
	}
}

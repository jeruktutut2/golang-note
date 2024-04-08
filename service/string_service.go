package service

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type StringService interface {
	Substring1(ctx context.Context, requestId string) (numberOfContain int)
	Substring2(ctx context.Context, requestId string) (arr []int)
	Subsequence1(ctx context.Context, requestId string) (arr []string)
	Subsequence2(ctx context.Context, requestId string) int
	Rotation(ctx context.Context, requestId string) string
	BinaryString(ctx context.Context, requestId string) (s string)
	Palindrome(ctx context.Context, requestId string) (s string)
	LexicographicRackString(ctx context.Context, requestId string) (rank int)
	PatternSearching(ctx context.Context, requestId string) (arr []int)
}

type StringServiceImplementation struct {
}

func NewStringService() StringService {
	return &StringServiceImplementation{}
}

// https://www.geeksforgeeks.org/number-of-substrings-of-one-string-present-in-other/
func (service *StringServiceImplementation) Substring1(ctx context.Context, requestId string) (numberOfContain int) {
	s1 := "aab"
	s2 := "aaaab"

	for i := 0; i < len(s1); i++ {
		var s3 string
		for j := i; j < len(s1); j++ {
			s3 += string(s1[j])
			isContain := strings.Contains(s2, s3)
			if isContain {
				numberOfContain++
			}
		}
	}
	return
}

// https://www.geeksforgeeks.org/print-all-substring-of-a-number-without-any-conversion/ and chatgpt to convert from javascript to golang
func (service *StringServiceImplementation) Substring2(ctx context.Context, requestId string) (arr []int) {
	n := 12345

	// Calculate the total number of digits
	s := int(math.Log10(float64(n)))

	// 0.5 has been added because it will return double value like 99.556
	d := int(math.Pow(10, float64(s)) + 0.5)
	k := d

	for n > 0 {
		// Print all the numbers from starting position
		for d > 0 {
			number := n / d
			d = d / 10
			arr = append(arr, number)
		}

		// Update the number
		n = n % k

		// Update the number of digits
		k = k / 10
		d = k
	}
	return
}

// https://www.geeksforgeeks.org/given-number-find-number-contiguous-subsequences-recursively-add-9/
func (service *StringServiceImplementation) Subsequence1(ctx context.Context, requestId string) (arr []string) {
	s := "4189"
	var count int
	count = 0
	for i := 0; i < len(s); i++ {
		var snum string
		var sum int
		for j := i + 0; j < len(s); j++ {
			snum += string(s[j])
			num, _ := strconv.Atoi(string(s[j]))
			sum += num
			if sum%9 == 0 {
				count++
				arr = append(arr, snum)
			}
		}
	}
	fmt.Println("count:", count)
	return
}

// https://www.geeksforgeeks.org/maximum-number-of-removals-of-given-subsequence-from-a-string/
func (service *StringServiceImplementation) Subsequence2(ctx context.Context, requestId string) int {
	s := "ggkssk"
	var i int
	var g int
	var gk int
	var gks int
	i = 0
	g = 0
	gk = 0
	gks = 0

	for i = 0; i < len(s); i++ {
		if string(s[i]) == "g" {
			g++
		} else if string(s[i]) == "k" {
			if g > 0 {
				g--
				gk++
			}
		} else if string(s[i]) == "s" {
			if gk > 0 {
				gk--
				gks++
			}
		}
	}
	return gks
}

// https://www.geeksforgeeks.org/left-rotation-right-rotation-string-2/
func (service *StringServiceImplementation) Rotation(ctx context.Context, requestId string) string {
	s1 := "GeeksforGeeks"
	// s2 := "GeeksforGeeks"
	// fmt.Println("s1:", string(s1[2:]), string(s1[:2]), string(s1[2:])+string(s1[:2]))
	// fmt.Println("s2:", string(s2[:len(s2)-2]), string(s2[len(s2)-2:]), string(s2[len(s2)-2:])+string(s2[:len(s2)-2]))
	return string(s1[2:]) + string(s1[:2])
}

// https://www.geeksforgeeks.org/what-is-binary-string/
func (service *StringServiceImplementation) BinaryString(ctx context.Context, requestId string) (s string) {
	binarystr1 := "10101010"
	binarystr2 := "01010101"
	fmt.Println("Length of binary string 1:", len(binarystr1))
	fmt.Println("Concatenation of binary strings:", binarystr1+binarystr2)
	fmt.Println("Substring of binary string 1:", string(binarystr1[2:6]))
	fmt.Println("Prefix of binary string 1:", string(binarystr1[0:3]))
	fmt.Println("Suffix of binary string 2:", string(binarystr2[4:8]))

	var hammingDist int
	for i := 0; i < len(binarystr1); i++ {
		if string(binarystr1[i]) != string(binarystr2[i]) {
			hammingDist++
		}
	}
	fmt.Println("Hamming distance between binary strings 1 and 2:", hammingDist)

	var hasRegularLanguage bool
forLoop:
	for i := 0; i < len(binarystr1); i++ {
		if string(binarystr1[i]) == "0" {
			hasRegularLanguage = true
			break forLoop
		}
	}
	if !hasRegularLanguage {
		fmt.Println("Does binary string 1 have a regular language? Yes")
	} else {
		fmt.Println("Does binary string 1 have a regular language? No")
	}

	binarynum1, _ := strconv.ParseInt(binarystr1, 2, 0)
	binarynum2, _ := strconv.ParseInt(binarystr2, 2, 0)
	fmt.Printf("Binary addition of %d and %d: %b\n", binarynum1, binarynum2, binarynum1+binarynum2)
	return binarystr1 + binarystr2
}

// https://www.youtube.com/watch?v=DXQuiPKl79Y
func (service *StringServiceImplementation) Palindrome(ctx context.Context, requestId string) (s string) {
	word := "maam"
	wordLength := len(word)
	loopTime := wordLength / 2
	isPalindrome := true
palindromeLoop:
	for i := 0; i < loopTime; i++ {
		if string(word[i]) != string(word[(wordLength-1)-i]) {
			fmt.Println("not a palindromme:", string(word[i]), string(word[(wordLength-1)-i]))
			isPalindrome = false
			break palindromeLoop
		}
	}

	if isPalindrome {
		s = word + "is a palindrome"
	} else {
		s = word + "is not a palindrome"
	}
	return
}

// https://www.geeksforgeeks.org/lexicographic-rank-of-a-string/
func fact(n int) int {
	if n <= 1 {
		return 1
	}
	return n * fact(n-1)
}

func (service *StringServiceImplementation) LexicographicRackString(ctx context.Context, requestId string) (rank int) {
	s := "string"
	// var rank int
	rank = 1

	for i := 0; i < len(s); i++ {
		var count int

		for j := i + 1; j < len(s); j++ {
			if string(s[i]) > string(s[j]) {
				count++
			}
		}
		rank += count * fact(len(s)-i-1)
	}
	return
}

// https://www.geeksforgeeks.org/naive-algorithm-for-pattern-searching/
func (service *StringServiceImplementation) PatternSearching(ctx context.Context, requestId string) (arr []int) {
	pat := "AABAACAADAABAABA"
	txt := "AABA"

	l1 := len(pat)
	l2 := len(txt)

	var i int
	var j int
	i = 0
	j = l2 - 1

	for i = 0; j < l1; i, j = i+1, j+1 {
		if txt == pat[i:j+1] {
			// fmt.Println("Pattern found at index", i)
			arr = append(arr, i)
		}
	}
	return
}

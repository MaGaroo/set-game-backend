package setgame

import "sort"

func max(a []int) int {
	result := a[0]
	for _, value := range a {
		if value > result {
			result = value
		}
	}
	return result
}

func removeByPosition(slice []int, positions []int) []int {
	sort.Ints(positions)
	reverse(positions)
	for _, pos := range positions {
		if pos != len(slice)-1 {
			slice[pos], slice[len(slice)-1] = slice[len(slice)-1], slice[pos]
		}
		slice = slice[:len(slice)-1]
	}
	return slice
}

func reverse(slice []int) {
	for left, right := 0, len(slice)-1; left < right; left, right = left+1, right-1 {
		slice[left], slice[right] = slice[right], slice[left]
	}
}

func appendByPosition(dst []int, src []int, positions []int) []int {
	for _, pos := range positions {
		dst = append(dst, src[pos])
	}
	return dst
}

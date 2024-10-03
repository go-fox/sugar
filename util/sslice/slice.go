// Package sslice
// MIT License
//
// # Copyright (c) 2024 go-fox
// Author https://github.com/go-fox/sugar
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package sslice

// Contain check if the target value is in the slice or not.
//
//	@param slice []T
//	@param target T
//	@return bool
//	@player
func Contain[T comparable](slice []T, target T) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}

	return false
}

// ContainBy returns true if predicate function return true.
//
//	@param slice []T
//	@param predicate func(item T) bool
//	@return bool
//	@player
func ContainBy[T any](slice []T, predicate func(item T) bool) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}

	return false
}

// ContainSubSlice check if the slice contains a given subslice or not.
//
//	@param slice []T
//	@param subSlice []T
//	@return bool
//	@player
func ContainSubSlice[T comparable](slice, subSlice []T) bool {
	if len(subSlice) == 0 {
		return true
	}
	if len(slice) == 0 {
		return false
	}

	elementMap := make(map[T]struct{}, len(slice))
	for _, item := range slice {
		elementMap[item] = struct{}{}
	}

	for _, item := range subSlice {
		if _, ok := elementMap[item]; !ok {
			return false
		}
	}
	return true
}

// Chunk creates a slice of elements split into groups the length of size.
//
//	@param slice []T
//	@param size int
//	@return [][]T
//	@player
func Chunk[T any](slice []T, size int) [][]T {
	result := [][]T{}

	if len(slice) == 0 || size <= 0 {
		return result
	}

	for _, item := range slice {
		l := len(result)
		if l == 0 || len(result[l-1]) == size {
			result = append(result, []T{})
			l++
		}

		result[l-1] = append(result[l-1], item)
	}

	return result
}

// Union 去重
func Union[T comparable, Slice ~[]T](lists ...Slice) Slice {
	var capLen int

	for _, list := range lists {
		capLen += len(list)
	}

	result := make(Slice, 0, capLen)
	seen := make(map[T]struct{}, capLen)

	for i := range lists {
		for j := range lists[i] {
			if _, ok := seen[lists[i][j]]; !ok {
				seen[lists[i][j]] = struct{}{}
				result = append(result, lists[i][j])
			}
		}
	}

	return result
}

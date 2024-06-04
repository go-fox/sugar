// Package sarray
// MIT License
//
// # Copyright (c) 2024 sugar
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
package sarray

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"

	"github.com/go-fox/sugar/internal/json"
	"github.com/go-fox/sugar/internal/rwmutex"
	"github.com/go-fox/sugar/util/sconv"
)

// New new and returns an empty array.
//
//	@param safe ...bool
//	@return *Array[V]
//	@player
func New[V any](safe ...bool) *Array[V] {
	return &Array[V]{
		mu:   rwmutex.New(safe...),
		data: make([]V, 0),
	}
}

// NewFromSlice returns an array of specified slices
//
//	@param data []T
//	@param safe ...bool
//	@return *Array[T]
//	@player
func NewFromSlice[T any](data []T, safe ...bool) *Array[T] {
	return &Array[T]{
		mu:   rwmutex.New(safe...),
		data: data,
	}
}

// Array wraps map type `[]any` and provides more map features.
type Array[V any] struct {
	mu   *rwmutex.RWMutex
	data []V
}

// Iterator is alias of IteratorAsc.
//
//	@receiver s
//	@param f func(key int, v V) bool
//	@player
func (s *Array[V]) Iterator(f func(key int, v V) bool) {
	s.IteratorAsc(f)
}

// IteratorAsc iterates the array readonly in ascending order with given callback function `f`.
//
//	@receiver s
//	@param f func(key int, v V) bool returns true, then it continues iterating; or false to stop.
//	@player
func (s *Array[V]) IteratorAsc(f func(key int, v V) bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, datum := range s.data {
		if !f(i, datum) {
			return
		}
	}
}

// IteratorDesc iterates the array readonly in descending order with given callback function `f`.
//
//	@receiver s
//	@param f func(key int, v V) bool returns true, then it continues iterating; or false to stop.
//	@player
func (s *Array[V]) IteratorDesc(f func(key int, v V) bool) {
	s.mu.Lock()
	s.mu.Unlock()
	for i := len(s.data) - 1; i >= 0; i-- {
		if !f(i, s.data[i]) {
			return
		}
	}
}

// Push append elements at the tail
//
//	@receiver s
//	@param item V
//	@player
func (s *Array[V]) Push(item V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = append(s.data, item)
}

// Remove Removes the element at the specified position in this array
//
//	@receiver s
//	@param index int
//	@player
func (s *Array[V]) Remove(index int) (v V, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.remove(index)
}

// RemoveValue Removes the element at the specified value in this array
//
//	@receiver s
//	@param value V
//	@return int
//	@return bool
//	@player
func (s *Array[V]) RemoveValue(value V) (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	index := s.search(value)
	_, ok := s.remove(index)
	return index, ok
}

// Search finds the specified value and return the position of the value in the array
//
//	@receiver s
//	@param value V
//	@return int
//	@player
func (s *Array[V]) Search(value V) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.search(value)
}

// Get returns the element at the specified position in this array.
//
//	@receiver s
//	@param index int index of the element to return
//	@return V the element at the specified position in this array
//	@return bool
//	@player
func (s *Array[V]) Get(index int) (v V, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if index < 0 || index >= len(s.data) {
		return
	}
	return s.data[index], true
}

// Set Replaces the element at the specified position in this list with the specified element.
//
//	@receiver s
//	@param index int index of the element to replace
//	@param element V element to be stored at the specified position
//	@player
func (s *Array[V]) Set(index int, element V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[index] = element
}

// Insert inserts a value at the specified position, move the current value back
//
//	@receiver s
//	@param index int
//	@param element V
//	@player
func (s *Array[V]) Insert(index int, elements ...V) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if index < 0 || index >= len(s.data) {
		return fmt.Errorf("sarray.Insert: index out of bounds: %d", index)
	}
	rear := append([]V{}, s.data[index:]...)
	s.data = append(s.data[0:index], elements...)
	s.data = append(s.data, rear...)
	return nil
}

// Sort
//
//	@receiver s
//	@param less func(v1, v2 interface{}) bool
//	@return *Array[V]
//	@player
func (s *Array[V]) Sort(less func(v1, v2 interface{}) bool) *Array[V] {
	s.mu.Lock()
	defer s.mu.Unlock()
	sort.Slice(s.data, func(i, j int) bool {
		return less(s.data[i], s.data[j])
	})
	return s
}

// Slice returns the underlying data of array.
//
//	@receiver s
//	@return []V
//	@player
func (s *Array[V]) Slice() []V {
	if s.mu.IsSafe() {
		s.mu.RLock()
		defer s.mu.RUnlock()
		array := make([]V, len(s.data))
		copy(array, s.data)
		return array
	} else {
		return s.data
	}
}

// Clone returns a new array, which is a copy of current array.
//
//	@receiver s
//	@return *Array[V]
//	@player
func (s *Array[V]) Clone() *Array[V] {
	s.mu.RLock()
	array := make([]V, len(s.data))
	copy(array, s.data)
	s.mu.RUnlock()
	return NewFromSlice(array, s.mu.IsSafe())
}

// Clear deletes all items of current array.
//
//	@receiver s
//	@return *Array[V]
//	@player
func (s *Array[V]) Clear() *Array[V] {
	s.mu.Lock()
	if len(s.data) > 0 {
		s.data = make([]V, 0)
	}
	s.mu.Unlock()
	return s
}

// Contains checks whether a value exists in the array.
//
//	@receiver s
//	@param value V
//	@return bool
//	@player
func (s *Array[V]) Contains(value V) bool {
	return s.Search(value) != -1
}

// Unique uniques the array, clear repeated items.
//
//	@receiver s
//	@return *Array[V]
//	@player
func (s *Array[V]) Unique() *Array[V] {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.data) == 0 {
		return s
	}
	var (
		ok          bool
		temp        interface{}
		uniqueSet   = make(map[interface{}]struct{})
		uniqueArray = make([]V, 0, len(s.data))
	)
	for i := 0; i < len(s.data); i++ {
		temp = s.data[i]
		if _, ok = uniqueSet[temp]; ok {
			continue
		}
		uniqueSet[temp] = struct{}{}
		uniqueArray = append(uniqueArray, temp)
	}
	s.data = uniqueArray
	return s
}

// Reverse makes array with elements in reverse order.
//
//	@receiver s
//	@return *Array[V]
//	@player
func (s *Array[V]) Reverse() *Array[V] {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, j := 0, len(s.data)-1; i < j; i, j = i+1, j-1 {
		s.data[i], s.data[j] = s.data[j], s.data[i]
	}
	return s
}

// Join joins array elements with a string `sep`.
//
//	@receiver s
//	@param sep string
//	@return string
//	@player
func (s *Array[V]) Join(sep string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.data) == 0 {
		return ""
	}
	buffer := bytes.NewBuffer(nil)
	for k, v := range s.data {
		str, err := sconv.ToString(v)
		if err != nil {
			str = ""
		}
		buffer.WriteString(str)
		if k != len(s.data)-1 {
			buffer.WriteString(sep)
		}
	}
	return buffer.String()
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
//
//	@receiver s
//	@return []byte
//	@return error
//	@player
func (s *Array[V]) MarshalJSON() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return json.Marshal(s.data)
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
//
//	@receiver s
//	@param data []byte
//	@return error
//	@player
func (s *Array[V]) UnmarshalJSON(data []byte) error {
	if s.data == nil {
		s.data = make([]V, 0)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := json.UnmarshalUseNumber(data, &s.data); err != nil {
		return err
	}
	return nil
}

// Filter filters the array readonly with custom callback function `f`.
//
//	@receiver s
//	@param f func(index int, v V) bool
//	@return *Array[V]
//	@player
func (s *Array[V]) Filter(f func(index int, v V) bool) []V {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]V, 0)
	for i, datum := range s.data {
		b := f(i, datum)
		if b {
			result = append(result, datum)
		}
	}
	return result
}

// Size returns the number of elements in this array.
//
//	@receiver s
//	@return int the number of elements in this list
//	@player
func (s *Array[V]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

func (s *Array[V]) search(value V) int {
	result := -1
	for index, datum := range s.data {
		if reflect.DeepEqual(datum, value) {
			result = index
			break
		}
	}
	return result
}

func (s *Array[V]) remove(index int) (v V, ok bool) {
	if index < 0 || index >= len(s.data) {
		return
	}
	// Determine array boundaries when deleting to improve deletion efficiency.
	if index == 0 {
		value := s.data[0]
		s.data = s.data[1:]
		return value, true
	} else if index == len(s.data)-1 {
		value := s.data[index]
		s.data = s.data[:index]
		return value, true
	}
	value := s.data[index]
	s.data = append(s.data[:index], s.data[index+1:]...)
	return value, true
}

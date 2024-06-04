// Package smap
// MIT License
//
// # Copyright (c) 2024 go-fox
// Author https://github.com/go-fox/fox
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
package smap

import (
	"github.com/go-fox/sugar/internal/rwmutex"
)

// New new and returns an empty hash map.
//
//	@param safe bool is it used during concurrency
//	@return *Map[K
//	@return V]
//	@player
func New[K comparable, V any](safe ...bool) *Map[K, V] {
	return &Map[K, V]{
		mu:   rwmutex.New(safe...),
		data: make(map[K]V),
	}
}

// Map wraps map type `map[K comparable]V any` and provides more map features.
type Map[K comparable, V any] struct {
	mu   *rwmutex.RWMutex
	data map[K]V
}

// Iterator iterates the hash map readonly with custom callback function `f`.
//
//	@receiver s
//	@param f func(key K, value V) bool returns true, then it continues iterating; or false to stop.
//	@player
func (s *Map[K, V]) Iterator(f func(key K, value V) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k, v := range s.data {
		if !f(k, v) {
			return
		}
	}
}

// CopyMap returns a shallow copy of the underlying data of the hash map.
//
//	@receiver s
//	@return map[K]V
//	@player
func (s *Map[K, V]) CopyMap() map[K]V {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data := make(map[K]V, len(s.data))
	for k, v := range s.data {
		data[k] = v
	}
	return data
}

// Clone returns a new hash map with copy of current map data.
//
//	@receiver s
//	@param safe ...bool
//	@return *Map[K
//	@return V]
//	@player
func (s *Map[K, V]) Clone(safe ...bool) *Map[K, V] {
	return &Map[K, V]{
		mu:   rwmutex.New(safe...),
		data: s.CopyMap(),
	}
}

// Filter filters the hash map readonly with custom callback function `f`.
//
//	@receiver s
//	@param f func(k K, v V) bool returns true, the result append; or false ignore this value
//	@player
func (s *Map[K, V]) Filter(f func(k K, v V) bool) map[K]V {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[K]V)
	for k, v := range s.data {
		if f(k, v) {
			result[k] = v
		}
	}
	return result
}

// DeleteWith delete the hash map with custom callback function `f`.
//
//	@receiver s
//	@param f func(k K, v V) bool returns true, the value delete; or false ignore this value
//	@player
func (s *Map[K, V]) DeleteWith(f func(k K, v V) bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range s.data {
		if f(k, v) {
			delete(s.data, k)
		}
	}
}

// Set sets key-value to the hash map.
//
//	@receiver s
//	@param key K
//	@param value V
//	@player
func (s *Map[K, V]) Set(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Get returns the value by given `key`.
//
//	@receiver s
//	@param key K
//	@return v
//	@return ok
//	@player
func (s *Map[K, V]) Get(key K) (v V, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok = s.data[key]
	return
}

// Del delete value by `key`
//
//	@receiver s
//	@param key K
//	@return {}
//	@player
func (s *Map[K, V]) Del(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

// Package satomic
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
package satomic

import "sync/atomic"

// Value warp atomic.Value
type Value[T any] struct {
	value atomic.Value
}

// New newa atomic value
//
//	@return *Value
//	@player
func New[T any]() *Value[T] {
	return &Value[T]{
		value: atomic.Value{},
	}
}

// Get implements the interface Get for atomic.Value.
//
//	@receiver v
//	@return T
//	@player
func (v *Value[T]) Get() T {
	return v.value.Load().(T)
}

// Store implements the interface Store for atomic.Value.
//
//	@receiver v
//	@param val T
//	@player
func (v *Value[T]) Store(val T) {
	v.value.Store(val)
}

// Swap implements the interface Swap for atomic.Value.
//
//	@receiver v
//	@param new T
//	@return old
//	@player
func (v *Value[T]) Swap(new T) (old T) {
	old = v.value.Load().(T)
	v.value.Store(new)
	return old
}

// CompareAndSwap implements the interface CompareAndSwap for atomic.Value.
//
//	@receiver v
//	@param old T
//	@param new T
//	@return swapped
//	@player
func (v *Value[T]) CompareAndSwap(old, new T) (swapped bool) {
	return v.value.CompareAndSwap(old, new)
}

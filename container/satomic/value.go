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

import (
	"fmt"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/go-fox/sugar/util/sconv"
)

// Value warp atomic.Value
type Value[T any] struct {
	atomic.Value
}

// New newa atomic value
//
//	@return *Value
//	@player
func New[T any]() *Value[T] {
	return &Value[T]{
		Value: atomic.Value{},
	}
}

// Load implements the interface Get for atomic.Value.
//
//	@receiver v
//	@return T
//	@player
func (v *Value[T]) Load() T {
	return v.Value.Load().(T)
}

// Store implements the interface Store for atomic.Value.
//
//	@receiver v
//	@param val T
//	@player
func (v *Value[T]) Store(val T) {
	v.Value.Store(val)
}

// Swap implements the interface Swap for atomic.Value.
//
//	@receiver v
//	@param new T
//	@return old
//	@player
func (v *Value[T]) Swap(new T) (old T) {
	old = v.Value.Load().(T)
	v.Value.Store(new)
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
	return v.Value.CompareAndSwap(old, new)
}

// IsEmpty implements the interface IsZero for reflect.Value.
//
//	@receiver v
//	@return bool
//	@player
func (v *Value[T]) IsEmpty() bool {
	val := v.Value.Load()
	return reflect.ValueOf(val).IsZero()
}

// Bool get value and convert to Bool value
//
//	@receiver v
//	@return bool
//	@return error
//	@player
func (v *Value[T]) Bool() (bool, error) {
	return sconv.ToBool(v.Value.Load())
}

// String get value and convert to string value
//
//	@receiver v
//	@return string
//	@return error
//	@player
func (v *Value[T]) String() (string, error) {
	return sconv.ToString(v.Value.Load())
}

// Int get value and convert to int64 value
//
//	@receiver v
//	@return int64
//	@return error
//	@player
func (v *Value[T]) Int() (int64, error) {
	return sconv.ToInt(v.Value.Load())
}

// Float get value and convert to float64 value
//
//	@receiver v
//	@return uint64
//	@return error
//	@player
func (v *Value[T]) Float() (float64, error) {
	return sconv.ToFloat(v.Value.Load())
}

// Duration get value and convert to time.Duration value
//
//	@receiver v
//	@return time.Duration
//	@return error
//	@player
func (v *Value[T]) Duration() (time.Duration, error) {
	return sconv.ToDuration(v.Value.Load())
}

// StringSlice get value and convert to []string value
//
//	@receiver v
//	@return []string
//	@return error
//	@player
func (v *Value[T]) StringSlice() ([]string, error) {
	vals, ok := v.Value.Load().([]interface{})
	if !ok {
		return []string{}, v.typeAssertError()
	}
	ret := make([]string, 0, len(vals))
	for _, val := range vals {
		v, err := sconv.ToString(val)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v)
	}
	return ret, nil
}

// StringMap get value and convert to map[string]string value
//
//	@receiver v
//	@return map[string]string
//	@return error
//	@player
func (v *Value[T]) StringMap() (map[string]string, error) {

	vals, ok := v.Value.Load().(map[string]interface{})
	if !ok {
		return map[string]string{}, v.typeAssertError()
	}
	ret := make(map[string]string, len(vals))
	for key, val := range vals {
		s, err := sconv.ToString(val)
		if err != nil {
			return map[string]string{}, err
		}
		ret[key] = s
	}
	return ret, nil
}

func (v *Value[T]) typeAssertError() error {
	return fmt.Errorf("type assert to %v failed", reflect.TypeOf(v.Load()))
}

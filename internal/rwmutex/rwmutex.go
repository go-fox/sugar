// Package rwmutex
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
package rwmutex

import "sync"

// RWMutex ...
type RWMutex struct {
	safe bool
	mu   sync.RWMutex
}

// New new a RWMutex
//
//	@param safe ...bool is it used during concurrency
//	@return *RWMutex
//	@player
func New(safe ...bool) *RWMutex {
	var s bool
	if len(safe) > 0 {
		s = safe[0]
	}
	return &RWMutex{
		safe: s,
		mu:   sync.RWMutex{},
	}
}

// Lock locks mutex for writing.
//
//	@receiver r
//	@player
func (r *RWMutex) Lock() {
	if r.safe {
		r.mu.Lock()
	}
}

// Unlock unlocks mutex for writing.
//
//	@receiver r
//	@player
func (r *RWMutex) Unlock() {
	if r.safe {
		r.mu.Unlock()
	}
}

// RLock locks mutex for read
//
//	@receiver r
//	@player
func (r *RWMutex) RLock() {
	if r.safe {
		r.mu.RLock()
	}
}

// RUnlock unlocks mutex for read
//
//	@receiver r
//	@player
func (r *RWMutex) RUnlock() {
	if r.safe {
		r.mu.RUnlock()
	}
}

// IsSafe return this RWMutex safe
//
//	@receiver r
//	@return bool
//	@player
func (r *RWMutex) IsSafe() bool {
	return r.safe
}

package spool

import "sync"

// Pool is an object pooling
type Pool[T any] struct {
	pool  *sync.Pool
	reset func(T)
}

// New news an object pool
//
//	@param factory func() T
//	@param reset ...func(T) T
//	@return *Pool[T]
//	@player
func New[T any](factory func() T, reset ...func(T)) *Pool[T] {
	p := &sync.Pool{}
	if factory == nil {
		p.New = func() any {
			var x T
			return x
		}
	} else {
		p.New = func() any {
			return factory()
		}
	}
	v := &Pool[T]{pool: p}
	if len(reset) > 0 {
		v.reset = reset[0]
	}
	return v
}

// Get 获取
//
//	@receiver v
//	@return T
func (v *Pool[T]) Get() T {
	return v.pool.Get().(T)
}

// Put 归还
//
//	@receiver v
//	@param x T
func (v *Pool[T]) Put(x T) {
	if v.reset != nil {
		v.reset(x)
	}
	v.pool.Put(x)
}

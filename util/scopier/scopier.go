// Package scopier
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
package scopier

import "reflect"

// DeepClone deep cloning
//
//	@param src T
//	@return T
//	@player
func DeepClone[T any](src T) T {
	c := cloner{
		ptrs: map[reflect.Type]map[uintptr]reflect.Value{},
	}
	result := c.clone(reflect.ValueOf(src))
	if result.Kind() == reflect.Invalid {
		var zeroValue T
		return zeroValue
	}

	return result.Interface().(T)
}

// CloneStruct clone structure
//
//	@param src interface{}
//	@return T
//	@player
func CloneStruct[T any](src interface{}) T {
	c := cloner{
		ptrs: make(map[reflect.Type]map[uintptr]reflect.Value),
	}
	result := c.cloneStruct(reflect.ValueOf(src))
	if result.Kind() == reflect.Invalid {
		var zeroValue T
		return zeroValue
	}
	return result.Interface().(T)
}

// CloneSlice clone slice
//
//	@param src T
//	@return T
//	@player
func CloneSlice[T any](src T) T {
	c := cloner{
		ptrs: map[reflect.Type]map[uintptr]reflect.Value{},
	}
	result := c.cloneSlice(reflect.ValueOf(src))
	if result.Kind() == reflect.Invalid {
		var zeroValue T
		return zeroValue
	}

	return result.Interface().(T)
}

// CloneMap clone map
//
//	@param src T
//	@return T
//	@player
func CloneMap[T any](src T) T {
	c := cloner{
		ptrs: map[reflect.Type]map[uintptr]reflect.Value{},
	}
	result := c.cloneMap(reflect.ValueOf(src))
	if result.Kind() == reflect.Invalid {
		var zeroValue T
		return zeroValue
	}

	return result.Interface().(T)
}

type cloner struct {
	ptrs map[reflect.Type]map[uintptr]reflect.Value
}

// clone return a duplicate of passed item.
func (c *cloner) clone(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Invalid:
		return reflect.ValueOf(nil)

	// bool
	case reflect.Bool:
		return reflect.ValueOf(v.Bool())

	//int
	case reflect.Int:
		return reflect.ValueOf(int(v.Int()))
	case reflect.Int8:
		return reflect.ValueOf(int8(v.Int()))
	case reflect.Int16:
		return reflect.ValueOf(int16(v.Int()))
	case reflect.Int32:
		return reflect.ValueOf(int32(v.Int()))
	case reflect.Int64:
		return reflect.ValueOf(v.Int())

	// uint
	case reflect.Uint:
		return reflect.ValueOf(uint(v.Uint()))
	case reflect.Uint8:
		return reflect.ValueOf(uint8(v.Uint()))
	case reflect.Uint16:
		return reflect.ValueOf(uint16(v.Uint()))
	case reflect.Uint32:
		return reflect.ValueOf(uint32(v.Uint()))
	case reflect.Uint64:
		return reflect.ValueOf(v.Uint())

	// float
	case reflect.Float32:
		return reflect.ValueOf(float32(v.Float()))
	case reflect.Float64:
		return reflect.ValueOf(v.Float())

	// complex
	case reflect.Complex64:
		return reflect.ValueOf(complex64(v.Complex()))
	case reflect.Complex128:
		return reflect.ValueOf(v.Complex())

	// string
	case reflect.String:
		return reflect.ValueOf(v.String())

	// array
	case reflect.Array, reflect.Slice:
		return c.cloneSlice(v)

	// map
	case reflect.Map:
		return c.cloneMap(v)

	// Ptr
	case reflect.Ptr:
		return c.clonePtr(v)

	// struct
	case reflect.Struct:
		return c.cloneStruct(v)

	// func
	case reflect.Func:
		return v

	// interface
	case reflect.Interface:
		return c.clone(v.Elem())

	}

	return reflect.Zero(v.Type())
}

func (c *cloner) cloneSlice(v reflect.Value) reflect.Value {
	if v.IsNil() {
		return reflect.Zero(v.Type())
	}

	arr := reflect.MakeSlice(v.Type(), v.Len(), v.Len())

	for i := 0; i < v.Len(); i++ {
		val := c.clone(v.Index(i))

		if !val.IsValid() {
			continue
		}

		item := arr.Index(i)
		if !item.CanSet() {
			continue
		}

		item.Set(val.Convert(item.Type()))
	}

	return arr
}

func (c *cloner) cloneMap(v reflect.Value) reflect.Value {
	if v.IsNil() {
		return reflect.Zero(v.Type())
	}

	clonedMap := reflect.MakeMap(v.Type())

	for _, key := range v.MapKeys() {
		value := v.MapIndex(key)
		clonedKey := c.clone(key)
		clonedValue := c.clone(value)

		if !isNillable(clonedKey) || !clonedKey.IsNil() {
			clonedKey = clonedKey.Convert(key.Type())
		}

		if (!isNillable(clonedValue) || !clonedValue.IsNil()) && clonedValue.IsValid() {
			clonedValue = clonedValue.Convert(value.Type())
		}

		if !clonedValue.IsValid() {
			clonedValue = reflect.Zero(clonedMap.Type().Elem())
		}

		clonedMap.SetMapIndex(clonedKey, clonedValue)
	}

	return clonedMap
}

func (c *cloner) clonePtr(v reflect.Value) reflect.Value {
	if v.IsNil() {
		return reflect.Zero(v.Type())
	}

	var newVal reflect.Value

	if v.Elem().CanAddr() {
		ptrs, exists := c.ptrs[v.Type()]
		if exists {
			if newVal, exists := ptrs[v.Elem().UnsafeAddr()]; exists {
				return newVal
			}
		}
	}

	newVal = c.clone(v.Elem())

	if v.Elem().CanAddr() {
		ptrs, exists := c.ptrs[v.Type()]
		if exists {
			if newVal, exists := ptrs[v.Elem().UnsafeAddr()]; exists {
				return newVal
			}
		}
	}

	clonedPtr := reflect.New(newVal.Type())
	clonedPtr.Elem().Set(newVal)

	return clonedPtr
}

func (c *cloner) cloneStruct(v reflect.Value) reflect.Value {
	clonedStructPtr := reflect.New(v.Type())
	clonedStruct := clonedStructPtr.Elem()

	if v.CanAddr() {
		ptrs := c.ptrs[clonedStructPtr.Type()]
		if ptrs == nil {
			ptrs = make(map[uintptr]reflect.Value)
			c.ptrs[clonedStructPtr.Type()] = ptrs
		}
		ptrs[v.UnsafeAddr()] = clonedStructPtr
	}

	for i := 0; i < v.NumField(); i++ {
		newStructValue := clonedStruct.Field(i)
		if !newStructValue.CanSet() {
			continue
		}

		clonedVal := c.clone(v.Field(i))
		if !clonedVal.IsValid() {
			continue
		}

		newStructValue.Set(clonedVal.Convert(newStructValue.Type()))
	}

	return clonedStruct
}

func isNillable(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Interface, reflect.Ptr, reflect.Func:
		return true
	}
	return false
}

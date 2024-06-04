// Package sconv
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
package sconv

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
)

// ToBool convert any to boolean.
//
//	@param v any
//	@return bool
//	@return error
//	@player
func ToBool(v any) (bool, error) {
	switch val := v.(type) {
	case bool:
		return val, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string:
		return strconv.ParseBool(fmt.Sprint(val))
	}
	return false, typeAssertError(v)
}

// ToBytes convert any to []byte.
//
//	@param value any
//	@return []byte
//	@return error
//	@player
func ToBytes(value any) ([]byte, error) {
	v := reflect.ValueOf(value)

	switch value.(type) {
	case int, int8, int16, int32, int64:
		number := v.Int()
		buf := bytes.NewBuffer([]byte{})
		buf.Reset()
		err := binary.Write(buf, binary.BigEndian, number)
		return buf.Bytes(), err
	case uint, uint8, uint16, uint32, uint64:
		number := v.Uint()
		buf := bytes.NewBuffer([]byte{})
		buf.Reset()
		err := binary.Write(buf, binary.BigEndian, number)
		return buf.Bytes(), err
	case float32:
		number := float32(v.Float())
		bits := math.Float32bits(number)
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, bits)
		return bytes, nil
	case float64:
		number := v.Float()
		bits := math.Float64bits(number)
		bytes := make([]byte, 8)
		binary.BigEndian.PutUint64(bytes, bits)
		return bytes, nil
	case bool:
		return strconv.AppendBool([]byte{}, v.Bool()), nil
	case string:
		return []byte(v.String()), nil
	case []byte:
		return v.Bytes(), nil
	default:
		newValue, err := json.Marshal(value)
		return newValue, err
	}
}

// ToChar convert any to []byte.
//
//	@param s string
//	@return []string
//	@player
func ToChar(s string) []string {
	c := make([]string, 0)
	if len(s) == 0 {
		c = append(c, "")
	}
	for _, v := range s {
		c = append(c, string(v))
	}
	return c
}

// ToInt convert any to int64.
//
//	@param val any
//	@return int64
//	@return error
//	@player
func ToInt(v any) (int64, error) {
	switch val := v.(type) {
	case int:
		return int64(val), nil
	case int8:
		return int64(val), nil
	case int16:
		return int64(val), nil
	case int32:
		return int64(val), nil
	case int64:
		return val, nil
	case uint:
		return int64(val), nil
	case uint8:
		return int64(val), nil
	case uint16:
		return int64(val), nil
	case uint32:
		return int64(val), nil
	case uint64:
		return int64(val), nil
	case float32:
		return int64(val), nil
	case float64:
		return int64(val), nil
	case json.Number:
		return val.Int64()
	case string:
		return strconv.ParseInt(val, 10, 64)
	}
	return 0, typeAssertError(v)
}

// ToString convert any to string
//
//	@param val any
//	@return string
//	@return error
//	@player
func ToString(v any) (string, error) {
	switch val := v.(type) {
	case string:
		return val, nil
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return fmt.Sprint(val), nil
	case []byte:
		return string(val), nil
	case fmt.Stringer:
		return val.String(), nil
	}
	return "", typeAssertError(v)
}

// ToFloat convert any to float64
//
//	@param v any
//	@return float64
//	@return error
//	@player
func ToFloat(v any) (float64, error) {
	switch val := v.(type) {
	case int:
		return float64(val), nil
	case int8:
		return float64(val), nil
	case int16:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case uint:
		return float64(val), nil
	case uint8:
		return float64(val), nil
	case uint16:
		return float64(val), nil
	case uint32:
		return float64(val), nil
	case uint64:
		return float64(val), nil
	case float32:
		return float64(val), nil
	case json.Number:
		return val.Float64()
	case float64:
		return val, nil
	case string:
		return strconv.ParseFloat(val, 64)
	}
	return 0.0, typeAssertError(v)
}

// ToDuration convert any to time.Duration
//
//	@param v any
//	@return time.Duration
//	@return error
//	@player
func ToDuration(v any) (time.Duration, error) {
	val, err := ToInt(v)
	if err != nil {
		return 0, err
	}
	return time.Duration(val), nil
}

// SliceToMap convert any to Map
//
//	@param array []T
//	@param iteratee func(T) (K, V)
//	@return map[K]V
//	@player
func SliceToMap[T any, K comparable, V any](array []T, iteratee func(T) (K, V)) map[K]V {
	result := make(map[K]V, len(array))
	for _, item := range array {
		k, v := iteratee(item)
		result[k] = v
	}

	return result
}

// MapToSlice convert any to time.Duration
//
//	@param aMap map[K]V
//	@param iteratee func(K, V) T
//	@return []T
//	@player
func MapToSlice[T any, K comparable, V any](aMap map[K]V, iteratee func(K, V) T) []T {
	result := make([]T, 0, len(aMap))

	for k, v := range aMap {
		result = append(result, iteratee(k, v))
	}

	return result
}

func typeAssertError(val any) error {
	return fmt.Errorf("type assert to %v failed", reflect.TypeOf(val))
}

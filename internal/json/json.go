// Package json
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
package json

import (
	"bytes"
	"encoding/json"
	"io"
)

// Marshal adapts to json/encoding Marshal API
//
//	@param val any
//	@return []byte
//	@return error
//	@player
func Marshal(val any) ([]byte, error) {
	return json.Marshal(val)
}

// Unmarshal adapts to json/encoding Unmarshal API
//
//	@param data []byte
//	@param val any
//	@return error
//	@player
func Unmarshal(data []byte, val any) error {
	return json.Unmarshal(data, val)
}

// MarshalIndent adapts to json/encoding MarshalIndent API
//
//	@param val any
//	@param prefix string
//	@param indent string
//	@return []byte
//	@return error
//	@player
func MarshalIndent(val any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(val, prefix, indent)
}

// UnmarshalUseNumber decodes the json data bytes to target interface using number option.
//
//	@param data []byte
//	@param val any
//	@return error
//	@player
func UnmarshalUseNumber(data []byte, val any) error {
	decoder := NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	return decoder.Decode(val)
}

// NewDecoder adapts to json/stream NewDecoder API.
//
//	@param reader io.Reader
//	@return *json.Decoder
//	@player
func NewDecoder(reader io.Reader) *json.Decoder {
	return json.NewDecoder(reader)
}

// Package surl
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
package surl

import "net/url"

// Scheme the scheme of url.URL
//
//	@param scheme string
//	@param isSecure bool
//	@return string
//	@player
func Scheme(scheme string, isSecure bool) string {
	if isSecure {
		return scheme + "s"
	}
	return scheme
}

// NewURL new an url.URL
//
//	@param scheme string
//	@param host string
//	@return *url.URL
//	@player
func NewURL(scheme, host string) *url.URL {
	return &url.URL{Scheme: scheme, Host: host}
}

// PickURL match specified scheme URL
//
//	@param urls []string
//	@param scheme string
//	@return string
//	@return error
//	@player
func PickURL(urls []string, scheme string) (string, error) {
	for _, e := range urls {
		u, err := url.Parse(e)
		if err != nil {
			return "", err
		}

		if u.Scheme == scheme {
			return u.Host, nil
		}
	}
	return "", nil
}

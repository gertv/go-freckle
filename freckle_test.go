// Copyright 2014 - anova r&d bvba. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package freckle

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const domain = "mydomain"
const token = "abcdefghijklmnopqrstuvwxyz"

func letsTestFreckle(ts *httptest.Server) Freckle {
	f := LetsFreckle(domain, token)
	f.Debug(true)
	f.base = ts.URL
	return f
}

func authenticated(t *testing.T, method, path string, fn func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, method, r.Method, "Should have been HTTP "+method)
		assert.Equal(t, path, r.URL.Path, "Should have been HTTP URL "+path)
		assert.Equal(t, domain, r.Header.Get("User-Agent"), "User-Agent header should have been set")
		assert.Equal(t, token, r.Header.Get("X-FreckleToken"), "X-FreckleToken header should have been set")

		fn(w, r)
	})
}

func response(body string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, body)
	}
}

func noContent() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	}
}

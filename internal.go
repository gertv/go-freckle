// Copyright 2014 - anova r&d bvba. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package freckle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type onResponse func([]byte, *http.Response) error

func (f Freckle) doHttpRequest(req *http.Request, fn onResponse) error {
	f.log("Request: HTTP %s %s", req.Method, req.URL)

	req.Header.Add("User-Agent", f.subdomain)
	req.Header.Add("X-FreckleToken", f.key)

	resp, err := f.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := body(resp)
	if err != nil {
		return err
	}

	f.log("Response: HTTP " + resp.Status)
	for key, value := range resp.Header {
		f.log("   %s: %s", key, value)
	}
	f.log("   %s", data)

	if resp.StatusCode >= 400 {
		return parseError(data, resp)
	}

	return fn(data, resp)
}

//
func (f Freckle) do(method, uri string, ps Parameters, is Inputs, fn onResponse) error {
	u := f.api(uri, ps)

	var b io.Reader
	if is != nil {
		data, err := json.Marshal(is)
		if err != nil {
			return err
		}
		b = bytes.NewReader(data)
		f.log("    %s", data)
	}

	req, err := http.NewRequest(method, u, b)
	if err != nil {
		return err
	}

	return f.doHttpRequest(req, fn)
}

// Try to parse the data into a FreckleError object
func parseError(data []byte, resp *http.Response) error {
	var result FreckleError
	err := json.Unmarshal(data, &result)
	if err == nil {
		return result
	} else {
		return err
	}
}

// Apply a slice of ParameterSetter functions to create a Parameters instance
func parameters(fns []ParameterSetter) Parameters {
	result := make(Parameters)
	for _, fn := range fns {
		fn(result)
	}
	return result
}

// Apply a slice of InputSetter functions to create an Inputs instance
func inputs(fns []InputSetter) Inputs {
	result := make(Inputs)
	for _, fn := range fns {
		fn(result)
	}
	return result
}

// Get the full API URL for a URL path and parameters
func (f Freckle) api(path string, ps Parameters) string {
	u := fmt.Sprintf("%s%s", f.base, path)
	if ps != nil && len(ps) > 0 {
		var v url.Values = make(url.Values)
		for key, value := range ps {
			v.Set(key, value)
		}
		u = fmt.Sprintf("%s?%s", u, v.Encode())
	}
	return u
}

// Extract body from the HTTP response
func body(resp *http.Response) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Simple internal logging method
func (f Freckle) log(msg string, data ...interface{}) {
	if f.debug {
		log.Printf("DEBUG: %s", fmt.Sprintf(msg, data...))
	}
}

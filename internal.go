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
	f.log("   %s", data)

	return fn(data, resp)
}

func (f Freckle) do(method, uri string, ps Parameters, is Inputs, fn onResponse) error {
	u := f.api(uri)
	if ps != nil {
		var v url.Values = make(url.Values)
		for key, value := range ps {
			v.Set(key, value)
		}
		u = fmt.Sprintf("%s?%s", u, v.Encode())
	}
	f.log("Request: HTTP %s %s", method, u)

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

func parameters(fns []ParameterSetter) Parameters {
	result := make(map[string]string)
	for _, fn := range fns {
		fn(result)
	}
	return result
}

func inputs(fns []InputSetter) Inputs {
	result := make(map[string]interface{})
	for _, fn := range fns {
		fn(result)
	}
	return result
}

func (f Freckle) api(path string) string {
	return fmt.Sprintf("%s%s", f.base, path)
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

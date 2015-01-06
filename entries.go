// Copyright 2014 - anova r&d bvba. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package freckle

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type EntriesAPI struct {
	freckle *Freckle
}

func (e EntriesAPI) ListEntries(fns ...ParameterSetter) (EntriesPage, error) {
	result := emptyEntriesPage(e.freckle)
	return result, e.freckle.do("GET", "/entries", parameters(fns), nil, result.onResponse)
}

func emptyEntriesPage(f *Freckle) EntriesPage {
	return EntriesPage{freckle: f}
}

func (p *EntriesPage) onResponse(data []byte, resp *http.Response) error {
	links := pagelinks(resp.Header.Get("Link"))
	var entries []Entry

	err := json.Unmarshal(data, &entries)
	p.links = links
	p.Entries = entries
	return err
}

func (e EntriesAPI) GetEntry(id int) (Entry, error) {
	var result Entry
	return result, e.freckle.do("GET", fmt.Sprintf("/entries/%d", id), nil, nil,
		func(data []byte, resp *http.Response) error {
			return json.Unmarshal(data, &result)
		})
}

func (e EntriesAPI) CreateEntry(date string, minutes int, fns ...InputSetter) (Entry, error) {
	is := inputs(fns)
	is["date"] = date
	is["minutes"] = minutes

	var result Entry
	return result, e.freckle.do("POST", "/entries", nil, is,
		func(output []byte, resp *http.Response) error {
			return json.Unmarshal(output, &result)
		})
}

func (e EntriesAPI) EditEntry(id int, fns ...InputSetter) (Entry, error) {
	var result Entry
	return result, e.freckle.do("PUT", fmt.Sprintf("/entries/%d", id), nil, inputs(fns),
		func(output []byte, resp *http.Response) error {
			return json.Unmarshal(output, &result)
		})
}

func (e EntriesAPI) MarkAsInvoiced(date string, id int) error {
	is := make(Inputs)
	is["date"] = date

	return e.freckle.do("PUT", fmt.Sprintf("/entries/%d/invoiced_outside_of_freckle", id), nil, is,
		func(output []byte, resp *http.Response) error {
			return nil
		})
}

func (e EntriesAPI) MarkMultipleAsInvoiced(date string, id ...int) error {
	is := make(Inputs)
	is["date"] = date
	is["entry_ids"] = id

	return e.freckle.do("PUT", "/entries/invoiced_outside_of_freckle", nil, is,
		func(output []byte, resp *http.Response) error {
			return nil
		})
}

func (e EntriesAPI) DeleteEntry(id int) error {
	return e.freckle.do("DELETE", fmt.Sprintf("/entries/%d", id), nil, nil,
		func(output []byte, resp *http.Response) error {
			return nil
		})
}

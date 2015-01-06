// Copyright 2014 - anova r&d bvba. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package freckle

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ProjectsAPI struct {
	freckle *Freckle
}

func (p ProjectsAPI) ListProjects(fns ...ParameterSetter) (ProjectsPage, error) {
	result := emptyProjectsPage(p.freckle)
	return result, p.freckle.do("GET", "/projects", parameters(fns), nil, result.onResponse)
}

func emptyProjectsPage(f *Freckle) ProjectsPage {
	return ProjectsPage{freckle: f}
}

func (p *ProjectsPage) onResponse(data []byte, resp *http.Response) error {
	links := pagelinks(resp.Header.Get("Link"))
	var projects []Project

	err := json.Unmarshal(data, &projects)
	p.links = links
	p.Projects = projects
	return err
}

func (p ProjectsAPI) GetProject(id int) (Project, error) {
	var result Project
	return result, p.freckle.do("GET", fmt.Sprintf("/projects/%d", id), nil, nil,
		func(data []byte, resp *http.Response) error {
			return json.Unmarshal(data, &result)
		})
}

func (p ProjectsAPI) CreateProject(name string, fns ...InputSetter) (Project, error) {
	is := inputs(fns)
	is["name"] = name

	var result Project
	return result, p.freckle.do("POST", "/projects", nil, is,
		func(data []byte, resp *http.Response) error {
			return json.Unmarshal(data, &result)
		})
}

func (p ProjectsAPI) GetEntries(id int) (EntriesPage, error) {
	result := emptyEntriesPage(p.freckle)
	return result, p.freckle.do("GET", fmt.Sprintf("/projects/%d/entries", id), nil, nil, result.onResponse)
}

func (p ProjectsAPI) GetInvoices(id int) ([]Invoice, error) {
	var result []Invoice
	return result, p.freckle.do("GET", fmt.Sprintf("/projects/%d/invoices", id), nil, nil,
		func(data []byte, resp *http.Response) error {
			return json.Unmarshal(data, &result)
		})
}

func (p ProjectsAPI) GetParticipants(id int) ([]Participant, error) {
	var result []Participant
	return result, p.freckle.do("GET", fmt.Sprintf("/projects/%d/participants", id), nil, nil,
		func(data []byte, resp *http.Response) error {
			return json.Unmarshal(data, &result)
		})
}

func (p ProjectsAPI) EditProject(id int, fns ...InputSetter) (Project, error) {
	var result Project
	return result, p.freckle.do("PUT", fmt.Sprintf("/projects/%d", id), nil, inputs(fns),
		func(data []byte, resp *http.Response) error {
			return json.Unmarshal(data, &result)
		})
}

func (p ProjectsAPI) MergeProject(target, toMerge int) error {
	is := make(Inputs)
	is["project_id"] = toMerge

	return p.freckle.do("PUT", fmt.Sprintf("/projects/%d/merge", target), nil, is,
		func(data []byte, resp *http.Response) error {
			return nil
		})
}

func (p ProjectsAPI) DeleteProject(id int) error {
	return p.freckle.do("DELETE", fmt.Sprintf("/projects/%d", id), nil, nil,
		func(data []byte, resp *http.Response) error {
			return nil
		})
}

func (p ProjectsAPI) ArchiveProject(id int) error {
	return p.freckle.do("PUT", fmt.Sprintf("/projects/%d/archive", id), nil, nil,
		func(data []byte, resp *http.Response) error {
			return nil
		})
}

func (p ProjectsAPI) UnarchiveProject(id int) error {
	return p.freckle.do("PUT", fmt.Sprintf("/projects/%d/unarchive", id), nil, nil,
		func(data []byte, resp *http.Response) error {
			return nil
		})
}

func (p ProjectsAPI) ArchiveMultipleProjects(ids ...int) error {
	is := make(Inputs)
	is["project_ids"] = ids

	return p.freckle.do("PUT", "/projects/archive", nil, is,
		func(data []byte, resp *http.Response) error {
			return nil
		})
}

func (p ProjectsAPI) UnarchiveMultipleProjects(ids ...int) error {
	is := make(Inputs)
	is["project_ids"] = ids

	return p.freckle.do("PUT", "/projects/unarchive", nil, is,
		func(data []byte, resp *http.Response) error {
			return nil
		})
}

func (p ProjectsAPI) DeleteMultipleProjects(ids ...int) error {
	is := make(Inputs)
	is["project_ids"] = ids

	return p.freckle.do("PUT", "/projects/delete", nil, is,
		func(data []byte, resp *http.Response) error {
			return nil
		})
}

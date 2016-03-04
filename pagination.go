// Copyright 2014 - anova r&d bvba. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package freckle

import (
	"net/http"
	"regexp"
)

const FirstPage = "first"
const LastPage = "last"
const NextPage = "next"
const PreviousPage = "prev"

// Is there a next page of entries?
func (p EntriesPage) HasNext() bool {
	return p.has(NextPage)
}

// Get the next page of entries
func (p EntriesPage) Next() (EntriesPage, error) {
	return p.fetch(NextPage)
}

// Is there a previous page of entries?
func (p EntriesPage) HasPrevious() bool {
	return p.has(PreviousPage)
}

// Get the previous page of entries
func (p EntriesPage) Previous() (EntriesPage, error) {
	return p.fetch(PreviousPage)
}

// Get the first page of entries
func (p EntriesPage) First() (EntriesPage, error) {
	return p.fetch(FirstPage)
}

// Get the last page of entries
func (p EntriesPage) Last() (EntriesPage, error) {
	return p.fetch(LastPage)
}

// Is there a next page of projects?
func (p ProjectsPage) HasNext() bool {
	return p.has(NextPage)
}

// Get the next page of projects
func (p ProjectsPage) Next() (ProjectsPage, error) {
	return p.fetch(NextPage)
}

// Is there a previous page of projects?
func (p ProjectsPage) HasPrevious() bool {
	return p.has(PreviousPage)
}

// Get the previous page of projects
func (p ProjectsPage) Previous() (ProjectsPage, error) {
	return p.fetch(PreviousPage)
}

// Get the first page of projects
func (p ProjectsPage) First() (ProjectsPage, error) {
	return p.fetch(FirstPage)
}

// Get the last page of projects
func (p ProjectsPage) Last() (ProjectsPage, error) {
	return p.fetch(LastPage)
}

// Get a channel to receive all entries. After all entries from the current
// page have been received, the next page will automatically be fetched.
func (p EntriesPage) AllEntries() chan Entry {
	result := make(chan Entry)
	go func() {
		p.push(result)
		close(result)
	}()
	return result
}

// Get a channel to receive all projects. Ater all projects from the current
// page have been receive, the next page will automatically be fetched.
func (p ProjectsPage) AllProjects() chan Project {
	result := make(chan Project)
	go func() {
		p.push(result)
		close(result)
	}()
	return result
}

// push all entries for current page and the next ones to the channel provided
func (p EntriesPage) push(c chan Entry) {
	for _, e := range p.Entries {
		c <- e
	}
	if p.HasNext() {
		next, err := p.Next()
		if err == nil {
			next.push(c)
		}
	}
}

// push all projects for current page and the next ones to the channel provided
func (p ProjectsPage) push(c chan Project) {
	for _, p := range p.Projects {
		c <- p
	}
	if p.HasNext() {
		next, err := p.Next()
		if err == nil {
			next.push(c)
		}
	}
}

// check if there is a page relative to the current one
func (p EntriesPage) has(id string) bool {
	_, ok := p.links[id]
	return ok
}

// check if there is a page relative to the current one
func (p ProjectsPage) has(id string) bool {
	_, ok := p.links[id]
	return ok
}

// fetch another entries page relative to the current one
func (p EntriesPage) fetch(id string) (EntriesPage, error) {
	f := p.freckle
	result := EntriesPage{freckle: f}

	req, err := http.NewRequest("GET", p.links[id], nil)
	if err != nil {
		return result, err
	}

	return result, f.doHttpRequest(req, result.onResponse)
}

// fetch another entries page relative to the current one
func (p ProjectsPage) fetch(id string) (ProjectsPage, error) {
	f := p.freckle
	result := emptyProjectsPage(f)

	req, err := http.NewRequest("GET", p.links[id], nil)
	if err != nil {
		return result, err
	}

	return result, f.doHttpRequest(req, result.onResponse)
}

// parse pagination links out of link header text
func pagelinks(header string) map[string]string {
	result := make(map[string]string)
	re := regexp.MustCompile("<(.*?)>; rel=\"(.*?)\"")
	for _, line := range re.FindAllStringSubmatch(header, -1) {
		rel := line[2]
		url := line[1]
		result[rel] = url
	}
	return result
}

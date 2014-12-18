// Copyright 2014 - anova r&d bvba. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This package implements a Go client for the Freckle V2 API.
//
// More information about the API itself can be found at
// http://developer.letsfreckle.com/v2/
package freckle

import (
	"net/http"
)

type Freckle struct {
	subdomain, key string
	debug          bool
	client         *http.Client
	base           string
}

// Start using the API here -
func LetsFreckle(subdomain, key string) Freckle {
	return Freckle{subdomain, key, false, &http.Client{}, "https://api.letsfreckle.com/v2"}
}

// Enable/disable debug mode. When debug mode is enabled,
// you will get additional logging showing the HTTP requests
// and responses
func (f *Freckle) Debug(d bool) {
	f.debug = d
}

// Configure a custom HTTP client (e.g. to configure a proxy server)
func (f *Freckle) Client(client *http.Client) {
	f.client = client
}

/*
Access the LetsFreckle v2 Entries API

more info at http://developer.letsfreckle.com/v2/entries
*/
func (f Freckle) EntriesAPI() EntriesAPI {
	return EntriesAPI{&f}
}

/*
Access the LetsFreckle v2 Projects API

more info at http://developer.letsfreckle.com/v2/projects
*/
func (f Freckle) ProjectsAPI() ProjectsAPI {
	return ProjectsAPI{&f}
}

// Get the error message for a Freckle API error
func (e FreckleError) Error() string {
	return e.Message
}

// Data type to represent values passed to a create, edit, ... calls
type Inputs map[string]interface{}

// Function to work with method values
type InputSetter func(Inputs)

// Data type to represents parameters on list/get/... calls
type Parameters map[string]string

// Function to work with API parameters
type ParameterSetter func(Parameters)

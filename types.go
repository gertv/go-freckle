// Copyright 2014 - anova r&d bvba. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package freckle

// Shared type definitions for API input/output

type Entry struct {
	Id          int            `json:"id,omitempty"`
	Date        string         `json:"date,omitempty"`
	User        Participant    `json:"user,omitempty"`
	Billable    bool           `json:"billable,omitempty"`
	Minutes     int            `json:"minutes,omitempty"`
	Description string         `json:"description,omitempty"`
	Project     ProjectSummary `json:"project,omitempty"`
	Tags        []Tag          `json:"tags,omitempty"`
	SourceUrl   string         `json:"source_url,omitempty"`
	InvoicedAt  string         `json:"invoiced_at,omitempty"`
	Invoice     Invoice        `json:"invoice,omitempty"`
	Import      Import         `json:"import,omitempty"`
	Url         string         `json:"url,omitempty"`
	CreatedAt   string         `json:"created_at,omitempty"`
	UpdatedAt   string         `json:"updated_at,omitempty"`
}

type Import struct {
	Id  int    `json:"id,omitempty"`
	Url string `json:"url,omitempty"`
}

type Invoice struct {
	Id     int     `json:"id,omitempty"`
	Number string  `json:"number,omitempty"`
	State  string  `json:"state,omitempty"`
	Total  float64 `json:"number,omitempty"`
	Url    string  `json:"url,omitempty"`
}

type Participant struct {
	Id              int    `json:"id,omitempty"`
	Email           string `json:"email,omitempty"`
	FirstName       string `json:"first_name,omitempty"`
	LastName        string `json:"last_name,omitempty"`
	ProfileImageUrl string `json:"profile_image_url,omitempty"`
	Url             string `json:"url,omitempty"`
}

type Project struct {
	Id                int           `json:"id,omitempty"`
	Name              string        `json:"name,omitempty"`
	BillingIncrement  int           `json:"billing_increment,omitempty"`
	Enabled           bool          `json:"enabled,omitempty"`
	Billable          bool          `json:"billable,omitempty"`
	Color             string        `json:"color,omitempty"`
	Url               string        `json:"url,omitempty"`
	Group             ProjectGroup  `json:"group, omitempty"`
	Minutes           int           `json:"minutes,omitempty"`
	BillableMinutes   int           `json:"billable_minutes,omitempty"`
	UnbillableMinutes int           `json:"unbillable_minutes,omitempty"`
	InvoicedMinutes   int           `json:"invoiced_minutes,omitempty"`
	RemainingMinutes  int           `json:"remaining_minutes,omitempty"`
	BudgetMinutes     int           `json:"budget_minutes,omitempty"`
	Import            Import        `json:"import,omitempty"`
	Invoices          []Invoice     `json:"invoices,omitempty"`
	Participants      []Participant `json:"participants,omitempty"`
	Entries           int           `json:"entries,omitempty"`
	EntriesUrl        string        `json:"entries_url,omitempty"`
	Expenses          int           `json:"expenses,omitempty"`
	ExpensesUrl       string        `json:"expenses_url,omitempty"`
	CreatedAt         string        `json:"created_at,omitempty"`
	UpdatedAt         string        `json:"updated_at,omitempty"`
}

type ProjectGroup struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

type ProjectSummary struct {
	Id               int    `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	BillingIncrement int    `json:"billing_increment,omitempty"`
	Enabled          bool   `json:"enabled,omitempty"`
	Billable         bool   `json:"billable,omitempty"`
	Color            string `json:"color,omitempty"`
	Url              string `json:"url,omitempty"`
}

type Tag struct {
	Id       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Billable bool   `json:"billable,omitempty"`
	Url      string `json:"url,omitempty"`
}

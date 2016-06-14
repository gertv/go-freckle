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

func TestListProjects(t *testing.T) {
	var ts *httptest.Server
	ts = httptest.NewServer(authenticated(t, "GET", "/projects", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Link", fmt.Sprintf("<%s%s?page=2>; rel=\"next\"", ts.URL, r.URL.Path))
		response(array_of_projects)(w, r)
	}))
	defer ts.Close()

	f := letsTestFreckle(ts)

	page, err := f.ProjectsAPI().ListProjects()
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, 1, len(page.Projects), "Should have one project")
	assert.True(t, page.HasNext(), "There should be a next page")
	project := page.Projects[0]
	assert.Equal(t, 1, len(project.Invoices), "Should have one invoice")
	invoice := project.Invoices[0]
	assert.Equal(t, "AA001", invoice.Reference, "Invoice reference mismatch")
	assert.Equal(t, 189.33, invoice.TotalAmount, "Invoice total_amount mismatch")

	page, err = page.Next()
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, 1, len(page.Projects), "Should have one project")
}

func TestListProjectsThroughChannel(t *testing.T) {
	page := 0
	var ts *httptest.Server
	ts = httptest.NewServer(authenticated(t, "GET", "/projects", func(w http.ResponseWriter, r *http.Request) {
		page += 1
		if page < 10 {
			w.Header().Set("Link", fmt.Sprintf("<%s%s?page=%d>; rel=\"next\"", ts.URL, r.URL.Path, page+1))
		}
		response(array_of_projects)(w, r)
	}))
	defer ts.Close()

	f := letsTestFreckle(ts)

	pp, err := f.ProjectsAPI().ListProjects()
	assert.Nil(t, err, "Error should be nil")
	projects := 0
	// reading through the channel should do 10 HTTP request, yielding 1 project each
	for _ = range pp.AllProjects() {
		projects += 1
	}
	assert.Equal(t, 10, projects, "Should have read 10 projects")
	assert.Equal(t, 10, page, "We should have read up to page 10")

}

func TestListProjectsWithParameters(t *testing.T) {
	ts := httptest.NewServer(authenticated(t, "GET", "/projects", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "true", r.URL.Query().Get("billable"))
		assert.Equal(t, "2014-12-18", r.URL.Query().Get("from"))
		response(array_of_projects)(w, r)
	}))
	defer ts.Close()

	f := letsTestFreckle(ts)

	page, err := f.ProjectsAPI().ListProjects(func(p Parameters) {
		p["billable"] = "true"
		p["from"] = "2014-12-18"
	})
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, 1, len(page.Projects), "Should have one project")
}

func TestListProjectsWithInvalidParameters(t *testing.T) {
	ts := httptest.NewServer(authenticated(t, "GET", "/projects", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		response(invalid_billing_code)(w, r)
	}))
	defer ts.Close()

	f := letsTestFreckle(ts)

	_, err := f.ProjectsAPI().ListProjects(func(p Parameters) {
		p["billable"] = "beer"
	})
	assert.NotNil(t, err, "Error should not be nil")
	if fe, ok := err.(FreckleError); ok {
		assert.Equal(t, fe.Message, "Validation Failed")
	} else {
		t.Errorf("Expected a FreckleError but got %s", err)
	}
}

func TestGetProject(t *testing.T) {
	ts := httptest.NewServer(authenticated(t, "GET", "/projects/37396", response(single_project)))
	defer ts.Close()

	f := letsTestFreckle(ts)

	_, err := f.ProjectsAPI().GetProject(37396)
	assert.Nil(t, err, "Error should be nil")
}

func TestCreateProject(t *testing.T) {
	ts := httptest.NewServer(authenticated(t, "POST", "/projects", response(single_project)))
	defer ts.Close()

	f := letsTestFreckle(ts)

	_, err := f.ProjectsAPI().CreateProject("Gear GmbH")
	assert.Nil(t, err, "Error should be nil")
}

func TestGetEntries(t *testing.T) {
	ts := httptest.NewServer(authenticated(t, "GET", "/projects/37396/entries", response(entries_for_project)))
	defer ts.Close()

	f := letsTestFreckle(ts)

	page, err := f.ProjectsAPI().GetEntries(37396)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, 1, len(page.Entries), "Should have one entry")
}

func TestGetInvoices(t *testing.T) {
	ts := httptest.NewServer(authenticated(t, "GET", "/projects/37396/invoices", response(invoices_for_project)))
	defer ts.Close()

	f := letsTestFreckle(ts)

	invoices, err := f.ProjectsAPI().GetInvoices(37396)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, 1, len(invoices), "Should have one invoice")
}

func TestGetParticipants(t *testing.T) {
	ts := httptest.NewServer(authenticated(t, "GET", "/projects/37396/participants", response(participants_for_project)))
	defer ts.Close()

	f := letsTestFreckle(ts)

	participants, err := f.ProjectsAPI().GetParticipants(37396)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, 1, len(participants), "Should have one participant")
}

func TestEditProject(t *testing.T) {
	ts := httptest.NewServer(authenticated(t, "PUT", "/projects/37396", response(single_project)))
	defer ts.Close()

	f := letsTestFreckle(ts)

	_, err := f.ProjectsAPI().EditProject(37396, func(i Inputs) {
		i["name"] = "New Name"
	})
	assert.Nil(t, err, "Error should be nil")
}

func TestMergeProject(t *testing.T) {
	ts := httptest.NewServer(authenticated(t, "PUT", "/projects/1234/merge", noContent()))
	defer ts.Close()

	f := letsTestFreckle(ts)

	err := f.ProjectsAPI().MergeProject(1234, 4567)
	assert.Nil(t, err, "Error should be nil")
}

func TestDeleteProject(t *testing.T) {
	ts := httptest.NewServer(authenticated(t, "DELETE", "/projects/1234", noContent()))
	defer ts.Close()

	f := letsTestFreckle(ts)

	err := f.ProjectsAPI().DeleteProject(1234)
	assert.Nil(t, err, "Error should be nil")
}

func TestArchiveProject(t *testing.T) {
	ts := httptest.NewServer(authenticated(t, "PUT", "/projects/1234/archive", noContent()))
	defer ts.Close()

	f := letsTestFreckle(ts)

	err := f.ProjectsAPI().ArchiveProject(1234)
	assert.Nil(t, err, "Error should be nil")
}

func TestUnarchiveProject(t *testing.T) {
	ts := httptest.NewServer(authenticated(t, "PUT", "/projects/1234/unarchive", noContent()))
	defer ts.Close()

	f := letsTestFreckle(ts)

	err := f.ProjectsAPI().UnarchiveProject(1234)
	assert.Nil(t, err, "Error should be nil")
}

func TestArchiveMultipleProjects(t *testing.T) {
	ts := httptest.NewServer(authenticated(t, "PUT", "/projects/archive", noContent()))
	defer ts.Close()

	f := letsTestFreckle(ts)

	err := f.ProjectsAPI().ArchiveMultipleProjects(1234, 4567)
	assert.Nil(t, err, "Error should be nil")
}

func TestUnarchiveMultipleProjects(t *testing.T) {
	ts := httptest.NewServer(authenticated(t, "PUT", "/projects/unarchive", noContent()))
	defer ts.Close()

	f := letsTestFreckle(ts)

	err := f.ProjectsAPI().UnarchiveMultipleProjects(1234, 4567)
	assert.Nil(t, err, "Error should be nil")
}

func TestDeleteMultipleProjects(t *testing.T) {
	ts := httptest.NewServer(authenticated(t, "PUT", "/projects/delete", noContent()))
	defer ts.Close()

	f := letsTestFreckle(ts)

	err := f.ProjectsAPI().DeleteMultipleProjects(1234, 4567)
	assert.Nil(t, err, "Error should be nil")
}

const array_of_projects = `[
  {
    "id": 37396,
    "name": "Gear GmbH",
    "billing_increment": 10,
    "enabled": true,
    "billable": true,
    "color": "#ff9898",
    "url": "https://api.letsfreckle.com/v2/projects/34580",
    "group": {
      "id": 3768,
      "name": "Sprockets, Inc.",
      "url": "https://api.letsfreckle.com/v2/project_groups/3768"
    },
    "minutes": 180,
    "billable_minutes": 120,
    "unbillable_minutes": 60,
    "invoiced_minutes": 120,
    "remaining_minutes": 630,
    "budget_minutes": 750,
    "import": {
      "id": 8910,
      "url": "https://api.letsfreckle.com/v2/imports/8910"
    },
    "invoices": [
      {
        "id": 12345678,
        "reference": "AA001",
        "invoice_date": "2013-07-09",
        "state": "unpaid",
        "total_amount": 189.33,
        "url": "https://api.letsfreckle.com/v2/invoices/12345678"
      }
    ],
    "participants": [
      {
        "id": 5538,
        "email": "john.test@test.com",
        "first_name": "John",
        "last_name": "Test",
        "profile_image_url": "https://api.letsfreckle.com/images/avatars/0000/0001/avatar.jpg",
        "url": "https://api.letsfreckle.com/v2/users/5538"
      }
    ],
    "entries": 0,
    "entries_url": "https://api.letsfreckle.com/v2/projects/34580/entries",
    "expenses": 0,
    "expenses_url": "https://api.letsfreckle.com/v2/projects/34580/expenses",
    "created_at": "2012-01-09T08:33:29Z",
    "updated_at": "2012-01-09T08:33:29Z"
  }
]`

const single_project = `{
  "id": 37396,
  "name": "Gear GmbH",
  "billing_increment": 10,
  "enabled": true,
  "billable": true,
  "color": "#ff9898",
  "url": "https://api.letsfreckle.com/v2/projects/34580",
  "group": {
    "id": 3768,
    "name": "Sprockets, Inc.",
    "url": "https://api.letsfreckle.com/v2/project_groups/3768"
  },
  "minutes": 180,
  "billable_minutes": 120,
  "unbillable_minutes": 60,
  "invoiced_minutes": 120,
  "remaining_minutes": 630,
  "budget_minutes": 750,
  "import": {
    "id": 8910,
    "url": "https://api.letsfreckle.com/v2/imports/8910"
  },
  "invoices": [
    {
      "id": 12345678,
      "reference": "AA001",
      "invoice_date": "2013-07-09",
      "state": "unpaid",
      "total_amount": 189.33,
      "url": "https://api.letsfreckle.com/v2/invoices/12345678"
    }
  ],
  "participants": [
    {
      "id": 5538,
      "email": "john.test@test.com",
      "first_name": "John",
      "last_name": "Test",
      "profile_image_url": "https://api.letsfreckle.com/images/avatars/0000/0001/avatar.jpg",
      "url": "https://api.letsfreckle.com/v2/users/5538"
    }
  ],
  "entries": 0,
  "entries_url": "https://api.letsfreckle.com/v2/projects/34580/entries",
  "expenses": 0,
  "expenses_url": "https://api.letsfreckle.com/v2/projects/34580/expenses",
  "created_at": "2012-01-09T08:33:29Z",
  "updated_at": "2012-01-09T08:33:29Z"
}`

const entries_for_project = `[
  {
    "id": 1,
    "date": "2012-01-09",
    "user": {
      "id": 5538,
      "email": "john.test@test.com",
      "first_name": "John",
      "last_name": "Test",
      "profile_image_url": "https://api.letsfreckle.com/images/avatars/0000/0001/avatar.jpg",
      "url": "https://api.letsfreckle.com/v2/users/5538"
    },
    "billable": true,
    "minutes": 60,
    "description": "freckle",
    "project": {
      "id": 37396,
      "name": "Gear GmbH",
      "billing_increment": 10,
      "enabled": true,
      "billable": true,
      "color": "#ff9898",
      "url": "https://api.letsfreckle.com/v2/projects/37396"
    },
    "tags": [
      {
        "id": 249397,
        "name": "freckle",
        "billable": true,
        "url": "https://api.letsfreckle.com/v2/tags/249397"
      }
    ],
    "source_url": "http://someapp.com/special/url/",
    "invoiced_at": "2012-01-10T08:33:29Z",
    "invoice": {
      "id": 12345678,
      "reference": "AA001",
      "invoice_date": "2013-07-09",
      "state": "unpaid",
      "total_amount": 189.33,
      "url": "https://api.letsfreckle.com/v2/invoices/12345678"
    },
    "import": {
      "id": 8910,
      "url": "https://api.letsfreckle.com/v2/imports/8910"
    },
    "url": "https://api.letsfreckle.com/v2/entries/1711626",
    "created_at": "2012-01-09T08:33:29Z",
    "updated_at": "2012-01-09T08:33:29Z"
  }
]`

const invoices_for_project = `[
  {
    "id": 26642,
    "state": "awaiting_payment",
    "number": "AB 0001",
    "invoice_date": "2013-07-09",
    "name": "Knockd, Freckle Support",
    "company_name": "John Test",
    "company_details": "1 Main Street\\r\\nMainsville, MA 11122",
    "recipient_details": "",
    "description": "",
    "footer": "",
    "show_hours": true,
    "show_details": false,
    "show_summaries": false,
    "customization": {
      "title": "Invoice",
      "date": "Date",
      "project": "Projects",
      "reference": "Invoice reference",
      "total_due": "Total amount due",
      "summary": "Summary",
      "work_time": "work time",
      "no_tax": "no tax",
      "tax": "tax",
      "subtotal": "subtotal",
      "total": "TOTAL",
      "report": "Report",
      "locale": "en-US",
      "currency_name": "",
      "currency_symbol": "$",
      "taxable_total": "Total taxable",
      "tax_total": "Total tax",
      "taxfree_total": "Total taxfree",
      "total_report": "TOTAL",
      "custom_css": null,
      "custom_html": "",
      "allow_paypal_invoice": true,
      "paypal_invoice_title": "",
      "paypal_currency_code": "USD",
      "paypal_address": "payment@test.com",
      "created_at": "2013-04-24T17:39:51Z",
      "updated_at": "2013-04-24T17:39:51Z"
    },
    "hours_calculation": {
      "calculation_method": "custom_hourly_rates",
      "custom_hourly_rates": [
        {
          "user": {
            "id": 5538,
            "email": "john.test@test.com",
            "first_name": "John",
            "last_name": "Test",
            "profile_image_url": "https://api.letsfreckle.com/images/avatars/0000/0001/avatar.jpg",
            "url": "https://api.letsfreckle.com/v2/users/5538"
          },
          "rate": 30.5,
          "hourly_rate_with_currency": "$30.50"
        }
      ]
    },
    "taxes": [
      {
        "id": 88292,
        "name": "Sales Tax",
        "percentage": 15.0
      }
    ],
    "amount_taxable": 100,
    "amount_taxfree": 0,
    "amount_tax_total": 0,
    "amount_total": 1,
    "amount_total_with_currency": "$1.00",
    "share_url": "https://apitest.letsfreckle.com/i/bqrnbojlbxqswtq9xla9uc40z",
    "payment": null,
    "payment_transactions": [
      {
        "description": "Notified that payment has been completed",
        "state": "paid",
        "payment_method": "paypal",
        "reference": "AP-AAAAABBBCCCCDDD111",
        "created_at": "2013-07-09T23:04:05Z",
        "updated_at": "2013-07-09T23:04:06Z"
      }
    ],
    "projects": [
      {
        "id": 37396,
        "name": "Gear GmbH",
        "billing_increment": 10,
        "enabled": true,
        "billable": true,
        "color": "#ff9898",
        "url": "https://api.letsfreckle.com/v2/projects/37396"
      }
    ],
    "entries": 0,
    "entries_url": "https://api.letsfreckle.com/v2/invoices/26642/entries",
    "expenses": 0,
    "expenses_url": "https://api.letsfreckle.com/v2/invoices/26642/expenses",
    "created_at": "2013-07-09T23:04:05Z",
    "updated_at": "2013-07-09T23:04:06Z",
    "from_address": null,
    "to_address": null,
    "tax_in_percent": null,
    "tax": null,
    "total": null,
    "subtotal": null
  }
]`

const participants_for_project = `[
  {
    "id": 5538,
    "email": "john.test@test.com",
    "first_name": "John",
    "last_name": "Test",
    "profile_image_url": "https://api.letsfreckle.com/images/avatars/0000/0001/avatar.jpg",
    "url": "https://api.letsfreckle.com/v2/users/5538",
    "state": "active",
    "role": "member",
    "participating_projects": 0,
    "participating_projects_url": "https://api.letsfreckle.com/v2/users/5538/participating_projects",
    "accessible_projects": 0,
    "accessible_projects_url": "https://api.letsfreckle.com/v2/users/5538/accessible_projects",
    "entries": 0,
    "entries_url": "https://api.letsfreckle.com/v2/users/5538/entries",
    "expenses": 0,
    "expenses_url": "https://api.letsfreckle.com/v2/users/5538/expenses",
    "add_project_access": "https://api.letsfreckle.com/v2/users/5538/project_access/add",
    "remove_project_access": "https://api.letsfreckle.com/v2/users/5538/project_access/remove",
    "created_at": "2010-06-09T20:44:57Z",
    "updated_at": "2010-06-09T20:44:57Z"
  }
]`

const invalid_billing_code = `{
	"errors":[{"code":"not_an_accepted_value","field":"billable","resource":"Project"}],
	"message":"Validation Failed"
}`

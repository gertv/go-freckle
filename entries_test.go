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

func TestListEntries(t *testing.T) {
	var ts *httptest.Server
	ts = httptest.NewServer(authenticated(t, "GET", "/entries", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Link", fmt.Sprintf("<%s%s?page=2>; rel=\"next\"", ts.URL, r.URL.Path))
		response(array_of_entries)(w, r)
	}))
	defer ts.Close()

	f := letsTestFreckle(ts)

	page, err := f.EntriesAPI().ListEntries()
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, 1, len(page.Entries), "Should have one entry")
	assert.True(t, page.HasNext(), "Should have a next page")
	assert.False(t, page.HasPrevious(), "Should not have a previous page")

	// Next() will just go back to the same server for now - just a first test for pagination
	page, err = page.Next()
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, 1, len(page.Entries), "Should have one entry")
	assert.True(t, page.HasNext(), "Should have a next page")
	assert.False(t, page.HasPrevious(), "Should not have a previous page")
}

func TestListEntriesThroughChannel(t *testing.T) {
	page := 0

	var ts *httptest.Server
	ts = httptest.NewServer(authenticated(t, "GET", "/entries", func(w http.ResponseWriter, r *http.Request) {
		page = page + 1
		if page < 10 {
			w.Header().Set("Link", fmt.Sprintf("<%s%s?page=%d>; rel=\"next\"", ts.URL, r.URL.Path, page+1))
		}
		response(array_of_entries)(w, r)
	}))
	defer ts.Close()

	f := letsTestFreckle(ts)

	items := 0
	ep, err := f.EntriesAPI().ListEntries()
	assert.Nil(t, err, "Error should be nil")
	// now let's read from the channel - it should go back to the server to fetch the next pages
	for _ = range ep.AllEntries() {
		items = items + 1
	}
	assert.Equal(t, 10, items, "Should have read 10 pages with 1 item each")
}

func TestGetEntry(t *testing.T) {
	ts := httptest.NewServer(authenticated(t, "GET", "/entries/1", response(single_entry)))
	defer ts.Close()

	f := letsTestFreckle(ts)

	_, err := f.EntriesAPI().GetEntry(1)
	assert.Nil(t, err, "Error should be nil")
}

func TestCreateEntry(t *testing.T) {
	ts := httptest.NewServer(authenticated(t, "POST", "/entries", response(single_entry)))
	defer ts.Close()

	f := letsTestFreckle(ts)

	_, err := f.EntriesAPI().CreateEntry("2014-12-18", 60, func(i Inputs) {
		i["description"] = "Very hard #support question"
	})
	assert.Nil(t, err, "Error should be nil")
}

func TestEditEntry(t *testing.T) {
	ts := httptest.NewServer(authenticated(t, "PUT", "/entries/1", response(single_entry)))
	defer ts.Close()

	f := letsTestFreckle(ts)

	_, err := f.EntriesAPI().EditEntry(1, func(i Inputs) {
		i["description"] = "Not so hard #support question"
	})
	assert.Nil(t, err, "Error should be nil")
}

func TestMarkAsInvoiced(t *testing.T) {
	ts := httptest.NewServer(authenticated(t, "PUT", "/entries/1/invoiced_outside_of_freckle", noContent()))
	defer ts.Close()

	f := letsTestFreckle(ts)

	err := f.EntriesAPI().MarkAsInvoiced("2014-12-18", 1)
	assert.Nil(t, err, "Error should be nil")
}

func TestMarkMultipleAsInvoiced(t *testing.T) {
	ts := httptest.NewServer(authenticated(t, "PUT", "/entries/invoiced_outside_of_freckle", noContent()))
	defer ts.Close()

	f := letsTestFreckle(ts)

	err := f.EntriesAPI().MarkMultipleAsInvoiced("2014-12-18", 1, 2, 3)
	assert.Nil(t, err, "Error should be nil")
}

func TestDeleteEntry(t *testing.T) {
	ts := httptest.NewServer(authenticated(t, "DELETE", "/entries/1", response(array_of_projects)))
	defer ts.Close()

	f := letsTestFreckle(ts)

	err := f.EntriesAPI().DeleteEntry(1)
	assert.Nil(t, err, "Error should be nil")
}

const array_of_entries = `[
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
      "number": "AA001",
      "state": "unpaid",
      "total": 189.33,
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

const single_entry = `{
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
    "number": "AA001",
    "state": "unpaid",
    "total": 189.33,
    "url": "https://api.letsfreckle.com/v2/invoices/12345678"
  },
  "import": {
    "id": 8910,
    "url": "https://api.letsfreckle.com/v2/imports/8910"
  },
  "url": "https://api.letsfreckle.com/v2/entries/1711626",
  "created_at": "2012-01-09T08:33:29Z",
  "updated_at": "2012-01-09T08:33:29Z"
}`

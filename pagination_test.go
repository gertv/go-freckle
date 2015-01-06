// Copyright 2014 - anova r&d bvba. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package freckle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagelinks(t *testing.T) {
	const header = `<https://apitest.letsfreckle.com/api/v2/users/?page=3&per_page=100>; rel="next",
  <https://apitest.letsfreckle.com/api/v2/users/?page=2&per_page=100>; rel="prev",
  <https://apitest.letsfreckle.com/api/v2/users/?page=1&per_page=100>; rel="first",
  <https://apitest.letsfreckle.com/api/v2/users/?page=50&per_page=100>; rel="last"`

	var links map[string]string = pagelinks(header)
	assert.Equal(t, "https://apitest.letsfreckle.com/api/v2/users/?page=3&per_page=100", links["next"])
	assert.Equal(t, "https://apitest.letsfreckle.com/api/v2/users/?page=2&per_page=100", links["prev"])
	assert.Equal(t, "https://apitest.letsfreckle.com/api/v2/users/?page=1&per_page=100", links["first"])
	assert.Equal(t, "https://apitest.letsfreckle.com/api/v2/users/?page=50&per_page=100", links["last"])
}

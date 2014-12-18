// Copyright 2014 - anova r&d bvba. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package freckle

// To get started with the API, use the LetsFreckle function
// and pass your domain and the Freckle V2 API token
func Example() {
	f := LetsFreckle("mycompany", "MyFreckleAPIV2Token")

	// once you have the Freckle object, just start using the API
	// through one of the ...API() functions
	f.EntriesAPI().ListEntries()
}

// Some functions, like create and edit functions, require
// additional input. This input can be provided through an
// anonymous function to add values to the Inputs object.
func ExampleInputSetter(f Freckle) {
	f.EntriesAPI().CreateEntry("2014-12-22", 60, func(i Inputs) {
		// here you can add addtional input data to your API call
		i["description"] = "My neat #development issue"
		i["project_name"] = "Customer Project"
	})
}

// Some functions, like the query functions, require
// additional parameters. These parameters can be provided through
// an anonymous function to add values to the Parameters object.
func ExampleParameterSetter(f Freckle) {
	f.EntriesAPI().ListEntries(func(p Parameters) {
		// here you can add addtional parameter data to your API call
		p["from"] = "2014-11-01"
		p["to"] = "2014-11-30"
	})
}

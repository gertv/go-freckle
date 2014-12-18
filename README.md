This project is a Go client to access the Freckle API V2.

Getting started
---------------

First, grab a copy of the client code

    go get github.com/gertv/go-freckle

Before you start using the Go client, you should have these
two bits of information at hand:
* your Freckle account subdomain
* a valid token for accessing the API V2.

Once you import the package, you can start using the API
through the `LetsFreckle(domain, token string)` function.

A minimal example of a command-line program that uses the
client would look like this:

```Go
package main

import (
  "github.com/gertv/go-freckle"
)

func main() {
  f := LetsFreckle("mycompany", "MyFreckleAPIV2Token")

  // once you have the Freckle object, just start using the API
  // through one of the ...API() functions
  f.EntriesAPI().ListEntries()
}
```


TODO
----

There's still some work to be done:

* Implement the other available resources of the V2 API
  * Timers
  * Expenses
  * Project Groups
  * ... (whatever else becomes availlable)
* Implement pagination for the `List...` methods
* Handling errors from the API
* Adding DSL methods for `Inputs` and `Parameters`

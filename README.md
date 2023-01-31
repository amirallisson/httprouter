# HTTP Routing Framework

A general library, called an HTTP Routing Framework,
that helps structure web applications based on patterns in end-user requests.
This general library supports a naming scheme for clients accessing resources
provided by a web application.

## API

The framework supports the following API:

```go
type HttpRouter struct {
    // the router itself
}

// Creates a new HttpRouter
func NewRouter() *HttpRouter

// Adds a new route to the HttpRouter
//
// Routes any request that matches the given `method` and `pattern` to a
// `handler`.
//
// `method`: should support arbitrary method strings, and least each of "POST",
//           "GET", "PUT", "DELETE". Method strings are case insensitive.
//
// `pattern`: patterns on the request _path_. Patterns can include arbitrary
//           "directories". Directories can include "captures" of the form
//           `:variable_name` that capture the actual directory value into a
//            HTTP query paramter. Leading and trailing '/' (slash) characters
//            are ignored.
//
// Example:
//
//   AddRoute("GET", "/users/:user/recent", RecentUserPosts)
//
// should map all GET requests with a path of the form "/users/*/recent" to the
// `RecentUserPosts` handler. It should populate the query parameter "user" with
// the value of the second directory.
//
// A request of the form "GET /users/cesar/recent HTTP/1.1" will call the
// RecentUserPosts with an `http.Request` with a `URL.RawQuery = "user=cesar"`
//
func (*HttpRouter) AddRoute(method string, pattern string, handler http.HandlerFunc)

// Conforms to the `http.Handler` interface
func (*HttpRouter) ServeHTTP(response http.ResponseWriter, request *http.Request)
```

You can run your unit tests with the command `go test`, which simply reports the
result of the test, and the reason for failure, if any, or you may add the `-v`
flag to see the verbose output of the unit tests.

For example, run the following from your top-level assignment 2 directory:
```bash
$ go test -v ./http_router
```
Equivalently, you may `cd` into the http_router directory and run the following:
```bash
$ go test -v
```
## Sample Application

The project also includes a sample application that uses the routing API.
The application is a simple microblogging application (similar to Twitter).
It uses an in-memory database to store users, threads, and messages, and
presents a JSON-based REST API for listing recent threads from a particular
user or all users the requesting user follows, posting new threads, responding
to threads, and creating new users.

### Running the application

The microblog client and server application can both be compiled using Go's
`go build` command. For example:

```bash
go build -o client ./microblog-client
go build -o server ./microblog-server
```

There is a simple Makefile that includes the two commands above.

To build the client and server programs, you can simply run the `make` command,
and the `make` utility will generate two executables, named `client` and `server`.


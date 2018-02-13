# Ersatz

> adjective
> (of a product) made or used as a substitute, typically an INFERIOR one, for something else.

Creates a stub API server using a simple configuration file you provide.

# Installation

```
go get github.com/peteclark-ft/ersatz
```

# Usage

First create a `fixtures.yml` file (see full example file [here](./_examples/example.yml)), to stub specific API calls.

For example, to stub an `/__health` API call, you can use the following configuration:

```
version: 1.0.0 # the ersatz fixtures version. Releases will be backwards compatible
fixtures:
   /__health: # the path of the stub
      get: # the http method
         status: 200 # the response status code to use
         headers: # a map of http response headers to values
            content-type: application/json
         body: # the body of the response - this will be serialised to match the content-type (if provided, otherwise json is the default)
            id: health-check
            ok: false
```

Similarly, to stub a `/__gtg` API call:

```
/__gtg:
   get:
      status: 200
      headers:
         content-type: plain/text; charset=US-ASCII
      body: OK # supports plaintext responses
```

Once you have created your file, start `ersatz` with:

```
ersatz --port 8080 ./fixtures.yml
```

# More Configuration Examples

Here is an example REST configuration for a `/{uuid}`, supporting `GET`, `PUT` and `DELETE`:

```
/85be197c-4fda-407b-8ae3-28bd81978616: # IMPORTANT: Paths must match exactly, as ersatz doesn't do anything clever with them. In this case, all other uuids will 404.
   get:
      status: 200
      body:
         content:
            title: Example Title
   put:
      status: 201
      body:
         message: Created ok
   delete:
      status: 204
```

Example redirect configuration to google.com:

```
/redirect-me:
   get:
      status: 301
      headers:
         Location: http://www.google.com
```

Example endpoint with expected `headers` and `queryParams`. Ersatz will respond with a `501 Unimplemented` to requests that do no match the configured expectations, which __should__ cause your sandbox tests to fail.

```
/expect:
   put:
      status: 200
      headers:
        x-returned-header: returned
      expectations:
        headers:
          x-expected-header: expected
        queryParams:
          param-name: expected-this-too
```

Ersatz also supports multiple expectations provided in an array; at least one expectation must much each request for the simulation to proceed, otherwise, it will return a `501 Unimplemented`.

```
/expect:
  put:
    status: 200
      headers:
        x-returned-header: returned
    expectations:
      - headers:
          x-expected-header: expected
      - queryParams:
          param-name: expected-this-too
```

# Why is Ersatz Useful?

* It's useful for local developer testing - you'd no longer need to point your local machine to real services in a test cluster.
* `fixtures.yml` files can be committed along with the codebase, so new developers can re-use your stubs to get up and running quickly.
* We can use it to simulate complex dependencies in CircleCI, allowing us to more easily test our OpenAPI files using DreddJS

# Road Map

* Support OpenAPI for more accurate stubs
* Comparisons between fixtures and the real API it is mocking

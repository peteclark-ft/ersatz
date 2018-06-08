# V1.0.0

Syntax guide for V1.0.0 ersatz fixtures configuration.

## Complete Syntax

* `version`: Must be `1.0.0`.
* `fixtures`: A map (key: endpoint path, value: Resource object), which contains the fixtures you wish to configure.

#### Resource Object

A map (key: HTTP Method, value: Response Object). Accepted HTTP Methods are `get | put | post | delete`. You must **not** specify the same HTTP Method twice, or the second will be overwritten.

#### Response Object

* **Required** `status`: The http status code to return in response.
* `headers`: Headers to return in the response. If `Content-Type` is set, this will dictate the format of the body. Supported content types are `application/json | text/plain | application/x-yaml`
* `body`: Polymorphic property, which supports values either of type string (should be used for `text/plain` responses) or of type Object, which will be serialised by default to JSON.
* `expectations`: Polymorphic property which supports passing either a single expectation, or multiple expectations in an array. Expectations check the request for provided `header` or `queryParam` values. If multiple expectations are provided, at least one set of expectations must pass for ersatz to proceed.

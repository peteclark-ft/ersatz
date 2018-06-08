# V2.0.0

Syntax guide for V2.0.0 ersatz fixtures configuration.

## Complete Syntax

* `version`: Must be `2.0.0`.
* `fixtures`: A map (key: endpoint path, value: Resource object), which contains the fixtures you wish to configure.

#### Resource Object

A map (key: HTTP Method, value: Either [Response Object](#response-object) or [Request Discriminator Object](#request-discriminator-object). Accepted HTTP Methods are `get | put | post | delete`. You must **not** specify the same HTTP Method twice, or the second will be overwritten.

Values may either be a single Response object, or many Request Discriminator objects. Request Discriminators are declared in an array, and allow you to specify different responses for different requests (discriminated by request properties other than the Path).

Discriminators are matched **in order**; if many discriminators match the same request, the **first** will be used.

#### Response Object

* **Required** `status`: The http status code to return in response.
* `headers`: Headers to return in the response. If `Content-Type` is set, this will dictate the format of the body. Supported content types are `application/json | text/plain | application/x-yaml`
* `body`: Polymorphic property, which supports values either of type string (should be used for `text/plain` responses) or of type Object, which will be serialised by default to JSON.

#### Request Discriminator Object

* **Required** `when`: Contains either `headers` or `queryParams` which are used to identify which response to use for the request.
   * `headers`: A map (key: string, value: string) of headers to look for the in the request.
   * `queryParams`: A map (key: string, value: string) of query parameters to look for the in the request.
* **Required** `response`: A [Response Object](#response-object) which will be used if the request matches the headers and query parameters specified.

Additionally, values included in the `when` statement can take the following formats:
* `${exists}`: Specifies that any value is acceptable for the header or query parameter, but it must be present.
* `${missing}`: Specifies that the value must not be present in the request.

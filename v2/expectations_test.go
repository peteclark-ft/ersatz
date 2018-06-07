package v2

import (
	"net/http/httptest"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoExpectationsReturnsTrue(t *testing.T) {
	e := Expectation{}

	r := httptest.NewRequest("GET", "/path", nil)
	ok := e.Validate(r)
	assert.True(t, ok)
}

func TestExpectHeadersSuccess(t *testing.T) {
	e := NewExpectation()
	e.Headers.Add("x-test", "im-here")

	r := httptest.NewRequest("GET", "/path", nil)
	r.Header.Add("x-test", "im-here")

	ok := e.Validate(r)
	assert.True(t, ok)
}

func TestExpectHeadersFailMissingHeader(t *testing.T) {
	e := NewExpectation()
	e.Headers.Add("x-test", "im-here")

	r := httptest.NewRequest("GET", "/path", nil)

	ok := e.Validate(r)
	assert.False(t, ok)
}

func TestExpectHeadersFailIncorrectHeader(t *testing.T) {
	e := NewExpectation()
	e.Headers.Add("x-test", "im-here")

	r := httptest.NewRequest("GET", "/path", nil)
	r.Header.Add("x-test", "im-not-right")

	ok := e.Validate(r)
	assert.False(t, ok)
}

func TestExpectQueryParamsSuccess(t *testing.T) {
	e := NewExpectation()
	e.QueryParams.Add("id", "expected-id")

	r := httptest.NewRequest("GET", "/path?id=expected-id", nil)

	ok := e.Validate(r)
	assert.True(t, ok)
}

func TestExpectQueryParamsFailsMissingParams(t *testing.T) {
	e := NewExpectation()
	e.QueryParams.Add("id", "expected-id")

	r := httptest.NewRequest("GET", "/path", nil)

	ok := e.Validate(r)
	assert.False(t, ok)
}

func TestExpectQueryParamsFailsIncorrectParam(t *testing.T) {
	e := NewExpectation()
	e.QueryParams.Add("id", "expected-id")

	r := httptest.NewRequest("GET", "/path?id=different-id", nil)

	ok := e.Validate(r)
	assert.False(t, ok)
}

const expectationTestYAML = `
headers:
   X-Expected-Header: expected
queryParams:
   expect: expected-this-too
`

func TestUnmarshalExpectation(t *testing.T) {
	e := Expectation{}
	err := yaml.Unmarshal([]byte(expectationTestYAML), &e)
	assert.NoError(t, err)

	assert.Equal(t, "expected", e.Headers.Get("x-expected-header"))
	assert.Equal(t, "expected-this-too", e.QueryParams.Get("expect"))
}

const expectationsTestYAML = `
headers:
   X-Expected-Header: expected
queryParams:
   expect: expected-this-too
`

func TestUnmarshalExpectationsWithNoArray(t *testing.T) {
	e := make(Expectations, 0)
	err := yaml.Unmarshal([]byte(expectationsTestYAML), &e)
	assert.NoError(t, err)

	require.Len(t, e, 1)

	assert.Equal(t, "expected", e[0].Headers.Get("x-expected-header"))
	assert.Equal(t, "expected-this-too", e[0].QueryParams.Get("expect"))
}

const multipleExpectationsTestYAML = `
- headers:
    X-Expected-Header: expected
  queryParams:
    expect: expected-this-too
- headers:
    x-another-expected-header: something-else
`

func TestUnmarshalExpectationsWithArray(t *testing.T) {
	e := make(Expectations, 0)
	err := yaml.Unmarshal([]byte(multipleExpectationsTestYAML), &e)
	assert.NoError(t, err)

	require.Len(t, e, 2)

	assert.Equal(t, "expected", e[0].Headers.Get("x-expected-header"))
	assert.Equal(t, "expected-this-too", e[0].QueryParams.Get("expect"))

	assert.Equal(t, "something-else", e[1].Headers.Get("x-another-expected-header"))
}

func TestExistsExpectation(t *testing.T) {
	e := Expectation{
		QueryParams: ExpectedQueryParams{
			map[string][]string{
				"not_me": []string{"${miss}"},
				"but_me": []string{"${exists}"},
			},
		},
	}
	r := httptest.NewRequest("GET", "/path?not_me=doesn't_matter_the_value&but_me=whatever", nil)
	assert.False(t, e.Validate(r))

	r = httptest.NewRequest("GET", "/path?but_me=whatever", nil)
	assert.True(t, e.Validate(r))
}

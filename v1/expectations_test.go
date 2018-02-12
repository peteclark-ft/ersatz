package v1

import (
	"net/http/httptest"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
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

const testYaml = `
headers:
   X-Expected-Header: expected
queryParams:
   expect: expected-this-too
`

func TestUnmarshalExpectation(t *testing.T) {
	e := Expectation{}
	err := yaml.Unmarshal([]byte(testYaml), &e)
	assert.NoError(t, err)

	assert.Equal(t, "expected", e.Headers.Get("x-expected-header"))
	assert.Equal(t, "expected-this-too", e.QueryParams.Get("expect"))
}

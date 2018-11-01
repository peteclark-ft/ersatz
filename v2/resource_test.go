package v2

import (
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
)

const resourceWithResponseTestYAML = `status: 200`

func TestResourceUnmarshal__WithResponse(t *testing.T) {
	r := Resource{}
	err := yaml.Unmarshal([]byte(resourceWithResponseTestYAML), &r)

	assert.NoError(t, err)
	assert.NotNil(t, r.Response)
	assert.Nil(t, r.Discriminators)
	assert.Equal(t, 200, r.Response.Status)
}

const resourceWithDiscriminatorsTestYAML = `
- when:
    headers:
      X-Example: example
  response:
    status: 501
- when:
    headers:
      X-Example: example
    queryParams:
      anotherExample: yet_another_example
  response:
    status: 400
    body: OK
`

func TestResourceUnmarshal__WithDiscriminators(t *testing.T) {
	r := Resource{}
	err := yaml.Unmarshal([]byte(resourceWithDiscriminatorsTestYAML), &r)

	assert.NoError(t, err)
	assert.NotNil(t, r.Discriminators)
	assert.Len(t, r.Discriminators, 2)

	assert.NotNil(t, r.Discriminators[0].When)
	assert.Equal(t, "example", r.Discriminators[0].When.Headers.Get("x-example"))

	assert.NotNil(t, r.Discriminators[0].Response)
	assert.Equal(t, 501, r.Discriminators[0].Response.Status)

	assert.NotNil(t, r.Discriminators[1].When)
	assert.Equal(t, "example", r.Discriminators[1].When.Headers.Get("x-example"))
	assert.Equal(t, "yet_another_example", r.Discriminators[1].When.QueryParams.Get("anotherExample"))

	assert.NotNil(t, r.Discriminators[1].Response)
	assert.Equal(t, 400, r.Discriminators[1].Response.Status)
	assert.Equal(t, "OK", r.Discriminators[1].Response.Body)
}

const resourceWithTemplatedValuesTestYAML = `
- when:
    headers:
      exists: ${exists}
      missing: ${missing}
  response:
    status: 501
`

func TestResourceUnmarshal__WithTemplatedValues(t *testing.T) {
	r := Resource{}
	err := yaml.Unmarshal([]byte(resourceWithTemplatedValuesTestYAML), &r)

	assert.NoError(t, err)
	assert.NotNil(t, r.Discriminators)
	assert.Len(t, r.Discriminators, 1)

	assert.NotNil(t, r.Discriminators[0].When)
	assert.Equal(t, "", r.Discriminators[0].When.Headers.Get("exists"))

	exists, ok := r.Discriminators[0].When.Headers.TemplatedValues["Exists"]
	assert.True(t, ok)
	assert.True(t, exists("I exist"))

	missing, ok := r.Discriminators[0].When.Headers.TemplatedValues["Missing"]
	assert.True(t, ok)
	assert.True(t, missing(""))
}

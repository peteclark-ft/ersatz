package v2

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiscriminator__Headers__TemplatedExists(t *testing.T) {
	d := NewRequestDiscriminator()
	d.Headers.TemplatedValues["x-testing"] = Exists

	r := httptest.NewRequest("GET", "/url", nil)
	r.Header.Add("x-testing", "I do exist")

	ok := d.SatisfiesDiscriminator(r)
	assert.True(t, ok)
}

func TestDiscriminator__Headers__TemplatedExistsFails(t *testing.T) {
	d := NewRequestDiscriminator()
	d.Headers.TemplatedValues["x-testing"] = Exists

	r := httptest.NewRequest("GET", "/url", nil)

	ok := d.SatisfiesDiscriminator(r)
	assert.False(t, ok)
}

func TestDiscriminator__Headers__TemplatedMissing(t *testing.T) {
	d := NewRequestDiscriminator()
	d.Headers.TemplatedValues["x-testing"] = Missing

	r := httptest.NewRequest("GET", "/url", nil)

	ok := d.SatisfiesDiscriminator(r)
	assert.True(t, ok)
}

func TestDiscriminator__Headers__TemplatedMissingFails(t *testing.T) {
	d := NewRequestDiscriminator()
	d.Headers.TemplatedValues["x-testing"] = Missing

	r := httptest.NewRequest("GET", "/url", nil)
	r.Header.Add("x-testing", "I shouldn't be here")

	ok := d.SatisfiesDiscriminator(r)
	assert.False(t, ok)
}

func TestDiscriminator__Headers__EqualValue(t *testing.T) {
	d := NewRequestDiscriminator()
	d.Headers.Add("x-testing", "specificValue")

	r := httptest.NewRequest("GET", "/url", nil)
	r.Header.Add("x-testing", "specificValue")

	ok := d.SatisfiesDiscriminator(r)
	assert.True(t, ok)
}

func TestDiscriminator__Headers__EqualValueFails(t *testing.T) {
	d := NewRequestDiscriminator()
	d.Headers.Add("x-testing", "specificValue")

	r := httptest.NewRequest("GET", "/url", nil)
	r.Header.Add("x-testing", "wrongValue")

	ok := d.SatisfiesDiscriminator(r)
	assert.False(t, ok)
}

func TestDiscriminator__QueryParams__TemplatedExists(t *testing.T) {
	d := NewRequestDiscriminator()
	d.QueryParams.TemplatedValues["testing"] = Exists

	r := httptest.NewRequest("GET", "/url?testing=exisssst", nil)

	ok := d.SatisfiesDiscriminator(r)
	assert.True(t, ok)
}

func TestDiscriminator__QueryParams__TemplatedExistsFails(t *testing.T) {
	d := NewRequestDiscriminator()
	d.QueryParams.TemplatedValues["testing"] = Exists

	r := httptest.NewRequest("GET", "/url", nil)

	ok := d.SatisfiesDiscriminator(r)
	assert.False(t, ok)
}

func TestDiscriminator__QueryParams__TemplatedMissing(t *testing.T) {
	d := NewRequestDiscriminator()
	d.QueryParams.TemplatedValues["testing"] = Missing

	r := httptest.NewRequest("GET", "/url", nil)

	ok := d.SatisfiesDiscriminator(r)
	assert.True(t, ok)
}

func TestDiscriminator__QueryParams__TemplatedMissingFails(t *testing.T) {
	d := NewRequestDiscriminator()
	d.QueryParams.TemplatedValues["testing"] = Missing

	r := httptest.NewRequest("GET", "/url?testing=missssssssing", nil)

	ok := d.SatisfiesDiscriminator(r)
	assert.False(t, ok)
}

func TestDiscriminator__QueryParams__EqualValue(t *testing.T) {
	d := NewRequestDiscriminator()
	d.QueryParams.Add("testing", "specificValue")

	r := httptest.NewRequest("GET", "/url?testing=specificValue", nil)

	ok := d.SatisfiesDiscriminator(r)
	assert.True(t, ok)
}

func TestDiscriminator__QueryParams__EqualValueFails(t *testing.T) {
	d := NewRequestDiscriminator()
	d.QueryParams.Add("testing", "specificValue")

	r := httptest.NewRequest("GET", "/url?testing=wrongValue", nil)

	ok := d.SatisfiesDiscriminator(r)
	assert.False(t, ok)
}

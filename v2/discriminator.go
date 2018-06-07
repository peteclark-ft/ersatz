package v2

import (
	"net/http"
	"net/textproto"
	"net/url"
)

func NewRequestDiscriminator() RequestDiscriminator {
	return RequestDiscriminator{
		Headers: Headers{
			MIMEHeader:      make(textproto.MIMEHeader),
			TemplatedValues: make(TemplatedValues),
		},
		QueryParams: QueryParams{
			Values:          make(url.Values),
			TemplatedValues: make(TemplatedValues),
		},
	}
}

func (arr Discriminators) AtLeastOneDiscriminatorIsSatisfied(req *http.Request) bool {
	for _, d := range arr {
		if d.When.SatisfiesDiscriminator(req) {
			return true
		}
	}
	return false
}

func (r RequestDiscriminator) SatisfiesDiscriminator(req *http.Request) bool {
	return r.Headers.Validate(req.Header) && r.QueryParams.Validate(req.URL.Query())
}

func (q QueryParams) Validate(actual url.Values) bool {
	for k, template := range q.TemplatedValues {
		v := actual.Get(k)
		if !template(v) {
			return false
		}
	}

	for name := range q.Values {
		expected := q.Get(name)
		value := actual.Get(name)

		ok := Contains(expected, value)
		if !ok {
			return false
		}
	}
	return true
}

// Validate validates the expected headers against the received headers
func (h Headers) Validate(actual http.Header) bool {
	for k, template := range h.TemplatedValues {
		v := actual.Get(k)
		if !template(v) {
			return false
		}
	}

	for name := range h.MIMEHeader {
		expected := h.Get(name)
		value := actual.Get(name)

		ok := Contains(expected, value)
		if !ok {
			return false
		}
	}
	return true
}

// Contains compares the expected values to the actual
func Contains(expected string, actual ...string) bool {
	for _, v := range actual {
		if v == expected {
			return true
		}
	}
	return false
}

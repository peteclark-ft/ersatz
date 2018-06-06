package v2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	expectationActionPrefix = "$#"
	expectationActionSuffix = "#$"

	expectationActionExists  = "exists"
	expectationActionMissing = "miss"
)

type Expectations []Expectation

// AtLeastOneExpectationPasses verifies that for multiple expectations at least one passes successfully, allowing the simulation to proceed
func (e Expectations) AtLeastOneExpectationPasses(r *http.Request) bool {
	for _, expectation := range e {
		if expectation.Validate(r) {
			return true
		}
	}
	return false
}

// AllExpectationsPass checks if all expectations pass
func (e Expectations) AllExpectationsPass(r *http.Request) bool {
	for _, expectation := range e {
		if !expectation.Validate(r) {
			return false
		}
	}
	return true
}

// UnmarshalJSON allows expectations to be declared either as an array or as a singular expectation
func (e *Expectations) UnmarshalJSON(d []byte) error {
	arr := make([]Expectation, 0)
	err := json.Unmarshal(d, &arr)
	if err == nil {
		*e = append(*e, arr...)
		return nil
	}

	single := Expectation{}
	err = json.Unmarshal(d, &single)
	if err != nil {
		return err
	}

	*e = append(*e, single)
	return nil
}

// Expectation contains expectations for the endpoint
type Expectation struct {
	Headers     ExpectedHeaders     `json:"headers"`
	QueryParams ExpectedQueryParams `json:"queryParams"`
}

// NewExpectation returns a setup Expectation used for test cases.
func NewExpectation() Expectation {
	return Expectation{Headers: ExpectedHeaders{MIMEHeader: make(textproto.MIMEHeader)}, QueryParams: ExpectedQueryParams{Values: make(url.Values)}}
}

// Validate validates the query params and headers if they exist
func (e *Expectation) Validate(r *http.Request) bool {
	if err := e.QueryParams.Validate(r.URL.Query()); err != nil {
		log.WithError(err).Error("Failed to validate request query params")
		return false
	}

	if err := e.Headers.Validate(r.Header); err != nil {
		log.WithError(err).Error("Failed to validate request headers")
		return false
	}

	return true
}

// ExpectedHeaders does what it says on the tin
type ExpectedHeaders struct {
	textproto.MIMEHeader
}

// UnmarshalJSON supports non-array declaration of headers
func (e *ExpectedHeaders) UnmarshalJSON(d []byte) error {
	headers := make(map[string]string)
	err := json.Unmarshal(d, &headers)
	if err != nil {
		return err
	}

	e.MIMEHeader = textproto.MIMEHeader{}
	for k, v := range headers {
		e.Add(k, v)
	}
	return nil
}

// Validate validates the expected headers against the received headers
func (e ExpectedHeaders) Validate(actual http.Header) error {
	for name := range e.MIMEHeader {
		expected := e.Get(name)
		value := actual.Get(name)
		if strings.HasPrefix(expected, expectationActionPrefix) && strings.HasSuffix(expected, expectationActionSuffix) {
			expectationAction := strings.Replace(expected, expectationActionPrefix, "", -1)
			expectationAction = strings.Replace(expectationAction, expectationActionSuffix, "", -1)
			if expectationAction != expectationActionExists && expectationAction != expectationActionMissing {
				return fmt.Errorf("invalid given action '%s'", expectationAction)
			}
			if expectationAction == expectationActionExists && len(value) == 0 {
				return fmt.Errorf("expected to find any value for header '%v'", name)
			}
			if expectationAction == expectationActionMissing && len(value) > 0 {
				return fmt.Errorf("expected missing '%v' header", name)
			}
			continue
		}

		ok := Contains(expected, value)
		if !ok {
			return fmt.Errorf(`expected to find header '%v' with value '%v' in actual request`, name, expected)
		}
	}
	return nil
}

// ExpectedQueryParams does what it says on the tin
type ExpectedQueryParams struct {
	url.Values
}

// UnmarshalJSON supports non-array declaration of headers
func (e *ExpectedQueryParams) UnmarshalJSON(d []byte) error {
	query := make(map[string]string)
	err := json.Unmarshal(d, &query)
	if err != nil {
		return err
	}

	e.Values = url.Values{}

	for k, v := range query {
		e.Add(k, v)
	}
	return nil
}

// Validate validates the expected query params vs the actual received params for the request
func (e ExpectedQueryParams) Validate(actual url.Values) error {
	for name := range e.Values {
		expected := e.Get(name)
		value := actual.Get(name)
		if strings.HasPrefix(expected, expectationActionPrefix) && strings.HasSuffix(expected, expectationActionSuffix) {
			expectationAction := strings.Replace(expected, expectationActionPrefix, "", -1)
			expectationAction = strings.Replace(expectationAction, expectationActionSuffix, "", -1)
			if expectationAction != expectationActionExists && expectationAction != expectationActionMissing {
				return fmt.Errorf("invalid given action '%s'", expectationAction)
			}
			if expectationAction == expectationActionExists && len(value) == 0 {
				return fmt.Errorf("expected to find any value for query param '%v'", name)
			}
			if expectationAction == expectationActionMissing && len(value) > 0 {
				return fmt.Errorf("expected missing '%v' query param", name)
			}
			continue
		}

		ok := Contains(expected, value)
		if !ok {
			return fmt.Errorf(`expected to find query param '%v' with value '%v' in actual request '%v'`, name, expected, actual.Encode())
		}
	}
	return nil
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

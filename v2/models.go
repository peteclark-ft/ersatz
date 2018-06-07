package v2

import (
	"net/http"

	"github.com/husobee/vestigo"
)

// Fixtures is the top level object with which endpoints are configured
type Fixtures map[string]Path

// Version returns the fixtures version number represented by this package
func (v Fixtures) Version() int {
	return 2
}

// Path maps absolute paths to http resources
type Path map[string][]Resource

// Resource mocks a particular http method for a given path
type Resource struct {
	Status               int               `json:"status"`
	Headers              map[string]string `json:"headers"`
	Body                 interface{}       `json:"body"`
	Expectations         Expectations      `json:"expectations"`
	AllExpectationsCheck bool              `json:"all_expectations_check"`
}

// Router allows us to test that paths are configured properly
type Router interface {
	Get(path string, handler http.HandlerFunc, middleware ...vestigo.Middleware)
	Put(path string, handler http.HandlerFunc, middleware ...vestigo.Middleware)
	Post(path string, handler http.HandlerFunc, middleware ...vestigo.Middleware)
	Delete(path string, handler http.HandlerFunc, middleware ...vestigo.Middleware)
}

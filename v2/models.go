package v2

import (
	"net/http"
	"net/textproto"
	"net/url"

	"github.com/husobee/vestigo"
)

// Fixtures is the top level object with which endpoints are configured
type Fixtures map[string]Path

// Version returns the fixtures version number represented by this package
func (v Fixtures) Version() int {
	return 2
}

// Path maps absolute paths to http resources
type Path map[string]Resource

type Resource struct {
	Discriminators Discriminators
	Response       Response
}

type Discriminators []Discriminator

type Discriminator struct {
	When     RequestDiscriminator `json:"when"`
	Response Response             `json:"response"`
}

type RequestDiscriminator struct {
	Headers     Headers     `json:"headers"`
	QueryParams QueryParams `json:"queryParams"`
}

// Headers does what it says on the tin
type Headers struct {
	textproto.MIMEHeader
	TemplatedValues TemplatedValues
}

// QueryParams does what it says on the tin
type QueryParams struct {
	url.Values
	TemplatedValues TemplatedValues
}

// TemplatedValues holds "special" values which can be used for fuzzy discriminators - i.e. ${exists} checks for the existence of the header
type TemplatedValues map[string]TemplatedFunction

type TemplatedFunction func(string) bool

// Response mocks a particular http method for a given path
type Response struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    interface{}       `json:"body"`
}

// Router allows us to test that paths are configured properly
type Router interface {
	Get(path string, handler http.HandlerFunc, middleware ...vestigo.Middleware)
	Put(path string, handler http.HandlerFunc, middleware ...vestigo.Middleware)
	Post(path string, handler http.HandlerFunc, middleware ...vestigo.Middleware)
	Delete(path string, handler http.HandlerFunc, middleware ...vestigo.Middleware)
}

package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/husobee/vestigo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockPaths(t *testing.T) {
	f := make(Fixtures)
	p := make(Path)
	r := Resource{Status: http.StatusOK}

	f["/example"] = p
	p["get"] = r
	p["post"] = r
	p["put"] = r
	p["delete"] = r

	mockRouter := new(MockRouter)
	mockRouter.On("Get", "/example", mock.Anything)
	mockRouter.On("Post", "/example", mock.Anything)
	mockRouter.On("Put", "/example", mock.Anything)
	mockRouter.On("Delete", "/example", mock.Anything)

	MockPaths(mockRouter, &f)
	mockRouter.AssertExpectations(t)
}

func TestMockResourcePlaintextResponse(t *testing.T) {
	res := Resource{Status: http.StatusTeapot, Body: "OK", Headers: map[string]string{"content-type": "text/plain"}}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	mockResource(res)(w, r)
	assert.Equal(t, "OK", w.Body.String())
	assert.Equal(t, http.StatusTeapot, w.Code)
}

func TestMockResourceJSONResponseIsDefault(t *testing.T) {
	res := Resource{Status: http.StatusTeapot, Body: "OK"}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	mockResource(res)(w, r)
	assert.Equal(t, `"OK"`, w.Body.String())
	assert.Equal(t, http.StatusTeapot, w.Code)
}

func TestMockResourceAddsHeaders(t *testing.T) {
	res := Resource{Status: http.StatusTeapot, Body: "OK", Headers: map[string]string{"x-request-id": "tid_1234"}}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	mockResource(res)(w, r)
	assert.Equal(t, `"OK"`, w.Body.String())
	assert.Equal(t, http.StatusTeapot, w.Code)
	assert.Equal(t, "tid_1234", w.Header().Get("x-request-id"))
}

func TestMockResourceYAMLResponse(t *testing.T) {
	res := Resource{
		Status: http.StatusTeapot,
		Body: struct {
			Greeting string `json:"greeting"`
		}{"hi"},
		Headers: map[string]string{"content-type": "application/x-yaml"},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	mockResource(res)(w, r)
	assert.Equal(t, "greeting: hi\n", w.Body.String())
	assert.Equal(t, http.StatusTeapot, w.Code)
}

func TestMockResourceNoBody(t *testing.T) {
	res := Resource{
		Status: http.StatusAccepted,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	mockResource(res)(w, r)
	assert.Equal(t, http.StatusAccepted, w.Code)
}

func TestMockResourceExpectationsFail(t *testing.T) {
	e := NewExpectation()
	e.QueryParams.Add("not-there", "v")

	res := Resource{
		Status:       http.StatusAccepted,
		Expectations: Expectations{e},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	mockResource(res)(w, r)
	assert.Equal(t, http.StatusNotImplemented, w.Code)
}

func TestMockResourceAtLeastOneExpectationPasses(t *testing.T) {
	e1 := NewExpectation()
	e1.QueryParams.Add("not-there", "v")

	e2 := NewExpectation()
	e1.QueryParams.Add("but-im-there", "v")

	res := Resource{
		Status:       http.StatusAccepted,
		Expectations: Expectations{e1, e2},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/?but-im-there=v", nil)

	mockResource(res)(w, r)
	assert.Equal(t, http.StatusAccepted, w.Code)
}

type MockRouter struct {
	mock.Mock
}

func (m *MockRouter) Get(path string, handler http.HandlerFunc, middleware ...vestigo.Middleware) {
	m.Called(path, handler)
}
func (m *MockRouter) Put(path string, handler http.HandlerFunc, middleware ...vestigo.Middleware) {
	m.Called(path, handler)
}
func (m *MockRouter) Post(path string, handler http.HandlerFunc, middleware ...vestigo.Middleware) {
	m.Called(path, handler)
}
func (m *MockRouter) Delete(path string, handler http.HandlerFunc, middleware ...vestigo.Middleware) {
	m.Called(path, handler)
}

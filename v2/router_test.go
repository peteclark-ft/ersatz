package v2

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
	p["get"] = []Resource{r}
	p["post"] = []Resource{r}
	p["put"] = []Resource{r}
	p["delete"] = []Resource{r}

	mockRouter := new(MockRouter)
	mockRouter.On("Get", "/example", mock.Anything)
	mockRouter.On("Post", "/example", mock.Anything)
	mockRouter.On("Put", "/example", mock.Anything)
	mockRouter.On("Delete", "/example", mock.Anything)

	MockPaths(mockRouter, &f)
	mockRouter.AssertExpectations(t)
}

func TestMockResourcePlaintextResponse(t *testing.T) {
	res := []Resource{Resource{Status: http.StatusTeapot, Body: "OK", Headers: map[string]string{"content-type": "text/plain"}}}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	mockResource(res)(w, r)
	assert.Equal(t, "OK", w.Body.String())
	assert.Equal(t, http.StatusTeapot, w.Code)
}

func TestMockResourceJSONResponseIsDefault(t *testing.T) {
	res := []Resource{Resource{Status: http.StatusTeapot, Body: "OK"}}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	mockResource(res)(w, r)
	assert.Equal(t, `"OK"`, w.Body.String())
	assert.Equal(t, http.StatusTeapot, w.Code)
}

func TestMockResourceAddsHeaders(t *testing.T) {
	res := []Resource{Resource{Status: http.StatusTeapot, Body: "OK", Headers: map[string]string{"x-request-id": "tid_1234"}}}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	mockResource(res)(w, r)
	assert.Equal(t, `"OK"`, w.Body.String())
	assert.Equal(t, http.StatusTeapot, w.Code)
	assert.Equal(t, "tid_1234", w.Header().Get("x-request-id"))
}

func TestMockResourceYAMLResponse(t *testing.T) {
	res := []Resource{
		Resource{
			Status: http.StatusTeapot,
			Body: struct {
				Greeting string `json:"greeting"`
			}{"hi"},
			Headers: map[string]string{"content-type": "application/x-yaml"},
		},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	mockResource(res)(w, r)
	assert.Equal(t, "greeting: hi\n", w.Body.String())
	assert.Equal(t, http.StatusTeapot, w.Code)
}

func TestMockResourceNoBody(t *testing.T) {
	res := []Resource{
		Resource{
			Status: http.StatusAccepted,
		},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	mockResource(res)(w, r)
	assert.Equal(t, http.StatusAccepted, w.Code)
}

func TestMockResourceExpectationsFail(t *testing.T) {
	e := NewExpectation()
	e.QueryParams.Add("not-there", "v")

	res := []Resource{
		Resource{
			Status:       http.StatusAccepted,
			Expectations: Expectations{e},
		},
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

	res := []Resource{
		Resource{
			Status:       http.StatusAccepted,
			Expectations: Expectations{e1, e2},
		},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/?but-im-there=v", nil)

	mockResource(res)(w, r)
	assert.Equal(t, http.StatusAccepted, w.Code)
}

func TestMockAllExpectationsMissFail(t *testing.T) {
	e1 := NewExpectation()
	e1.QueryParams.Add("not-there", "$#miss#$")
	e1.QueryParams.Add("but-im-there", "$#miss#$")

	res := []Resource{
		Resource{
			Status:       http.StatusAccepted,
			Expectations: Expectations{e1},
		},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/?but-im-there=v", nil)

	mockResource(res)(w, r)
	assert.Equal(t, http.StatusNotImplemented, w.Code)
}

func TestMockAllExpectationsSuccess(t *testing.T) {
	e1 := NewExpectation()
	e1.QueryParams.Add("not-there", "$#miss#$")

	e2 := NewExpectation()
	e2.QueryParams.Add("but-im-there", "$#exists#$")

	e3 := NewExpectation()
	e3.Headers.Add("Accept", "text/plain")

	res := []Resource{
		Resource{
			Status:       http.StatusAccepted,
			Expectations: Expectations{e1, e2, e3},
		},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/?but-im-there=v", nil)
	r.Header.Add("Accept", "text/plain")

	mockResource(res)(w, r)
	assert.Equal(t, http.StatusAccepted, w.Code)
}

func TestMockMultipleResources(t *testing.T) {
	expect1 := NewExpectation()
	expect1.QueryParams.Add("not-there", "${miss}")
	expect1.QueryParams.Add("but-im-there", "${miss}")

	res1 := Resource{
		Status:       http.StatusAccepted,
		Expectations: Expectations{expect1},
	}

	expect3 := NewExpectation()
	expect3.Headers.Add("Accept", "text/plain")
	res2 := Resource{
		Status:       http.StatusContinue,
		Expectations: Expectations{expect3},
	}

	res := []Resource{res1, res2}
	mockedResource := mockResource(res)

	// matching res2
	recorder1 := httptest.NewRecorder()
	req1 := httptest.NewRequest("GET", "/?but-im-there=v", nil)
	req1.Header.Add("Accept", "text/plain")
	mockedResource(recorder1, req1)
	assert.Equal(t, http.StatusContinue, recorder1.Code)

	// matching res1
	recorder2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/", nil)
	mockedResource(recorder2, req2)
	assert.Equal(t, http.StatusAccepted, recorder2.Code)

	// no definition for this situation
	recorder3 := httptest.NewRecorder()
	req3 := httptest.NewRequest("GET", "/?but-im-there=v", nil)
	mockedResource(recorder3, req3)
	assert.Equal(t, http.StatusNotImplemented, recorder3.Code)
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

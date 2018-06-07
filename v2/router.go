package v2

import (
	"encoding/json"
	"mime"
	"net/http"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// MockPaths adds endpoints to the provided router as per the ersatz-fixtures.yml
func MockPaths(r Router, paths *Fixtures) {
	for p, path := range *paths {
		for method, resource := range path {
			switch method {
			case "get":
				r.Get(p, mockResource(resource))
			case "post":
				r.Post(p, mockResource(resource))
			case "put":
				r.Put(p, mockResource(resource))
			case "delete":
				r.Delete(p, mockResource(resource))
			}
		}
	}
}

func mockResource(res Resource) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if res.Discriminators == nil {
			writeMockResponse(res.Response, w, r)
			return
		}

		if res.Discriminators != nil && !res.Discriminators.AtLeastOneDiscriminatorIsSatisfied(r) {
			w.WriteHeader(http.StatusNotImplemented)
			return
		}

		for _, d := range res.Discriminators {
			if d.When.SatisfiesDiscriminator(r) {
				writeMockResponse(d.Response, w, r)
				return
			}
		}
	}
}

func writeMockResponse(res Response, w http.ResponseWriter, r *http.Request) {
	for k, v := range res.Headers {
		w.Header().Add(k, v)
	}

	w.WriteHeader(res.Status)
	if res.Body == nil {
		return
	}

	contentType, ok := res.Headers["content-type"]
	if !ok { // assume json
		contentType = "application/json"
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		log.WithError(err).Error("Failed to parse media type")
		return
	}

	var output []byte

	switch mediaType {
	case "application/json":
		output, err = json.Marshal(res.Body)
	case "application/x-yaml":
		output, err = yaml.Marshal(res.Body)
	case "text/plain":
		output = []byte(res.Body.(string))
	}

	if err != nil {
		log.WithError(err).Error("Failed to marshal body")
		return
	}

	w.Write(output)
}

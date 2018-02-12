package v1

import (
	"encoding/json"
	"mime"
	"net/http"

	"github.com/Financial-Times/http-handlers-go/httphandlers"
	"github.com/ghodss/yaml"
	"github.com/husobee/vestigo"
	metrics "github.com/rcrowley/go-metrics"
	log "github.com/sirupsen/logrus"
)

type Fixtures map[string]Path

func (v Fixtures) Version() int {
	return 1
}

type Path map[string]Resource

type Resource struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    interface{}       `json:"body"`
}

func MockPaths(port string, paths *Fixtures) {
	unmonitoredRouter := vestigo.NewRouter()
	var r http.Handler = unmonitoredRouter
	r = httphandlers.TransactionAwareRequestLoggingHandler(log.StandardLogger(), r)
	r = httphandlers.HTTPMetricsHandler(metrics.DefaultRegistry, r)

	for p, path := range *paths {
		for method, resource := range path {
			switch method {
			case "get":
				unmonitoredRouter.Get(p, mock(resource))
			case "post":
				unmonitoredRouter.Post(p, mock(resource))
			case "put":
				unmonitoredRouter.Put(p, mock(resource))
			case "delete":
				unmonitoredRouter.Delete(p, mock(resource))
			}
		}
	}

	log.Info("Ready to simulate requests!")
	http.Handle("/", r)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Unable to start: %v", err)
	}
}

func mock(res Resource) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
			log.WithError(err).Error("Failed to parse media type...")
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
			log.WithError(err).Error("Failed to marshal body...")
			return
		}

		w.Write(output)
	}
}

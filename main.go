package main

import (
	"encoding/json"
	"io/ioutil"
	"mime"
	"net/http"
	"os"

	"github.com/Financial-Times/http-handlers-go/httphandlers"
	yaml "github.com/ghodss/yaml"
	"github.com/husobee/vestigo"
	"github.com/jawher/mow.cli"
	metrics "github.com/rcrowley/go-metrics"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

	app := cli.App("ersatz", "Mocks shit")
	port := app.String(cli.StringOpt{
		Name:   "port",
		Value:  "8081",
		Desc:   "Port to run ersatz on",
		EnvVar: "PORT",
	})

	configFile := app.StringArg("FILE", "", "The fixtures config file you wish to run")

	app.Action = func() {
		f, err := ioutil.ReadFile(*configFile)
		if err != nil {
			log.WithError(err).Error("Failed to read yml")
			return
		}

		paths := make(map[string]path)
		err = yaml.Unmarshal(f, &paths)

		if err != nil {
			log.WithError(err).Error("Failed to marshal yml")
			return
		}

		mockPaths(*port, paths)
	}

	app.Run(os.Args)
}

func mock(res resource) func(w http.ResponseWriter, r *http.Request) {
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

func mockPaths(port string, paths map[string]path) {
	unmonitoredRouter := vestigo.NewRouter()
	var r http.Handler = unmonitoredRouter
	r = httphandlers.TransactionAwareRequestLoggingHandler(log.StandardLogger(), r)
	r = httphandlers.HTTPMetricsHandler(metrics.DefaultRegistry, r)

	for p, path := range paths {
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

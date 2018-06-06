package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/Financial-Times/http-handlers-go/httphandlers"
	"github.com/ghodss/yaml"
	"github.com/husobee/vestigo"
	"github.com/jawher/mow.cli"
	"github.com/peteclark-ft/ersatz/v1"
	"github.com/peteclark-ft/ersatz/v2"
	log "github.com/sirupsen/logrus"
)

var configureOnce = &sync.Once{}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

	app := cli.App("ersatz", "Super basic stubbing tool for sandboxed component testing.")
	port := app.String(cli.StringOpt{
		Name:   "port p",
		Value:  "9000",
		Desc:   "Port to run ersatz on",
		EnvVar: "PORT",
	})

	fixtures := app.String(cli.StringOpt{
		Name:   "fixtures f",
		Value:  "./_ft/ersatz-fixtures.yml",
		Desc:   "Fixtures file to use to simulate requests",
		EnvVar: "FIXTURES",
	})

	app.Action = func() {
		yml, err := ioutil.ReadFile(*fixtures)
		if err != nil {
			runServer(*port, nil)
			return
		}

		ersatz := ersatz{}
		err = yaml.Unmarshal(yml, &ersatz)
		if err != nil {
			log.WithError(err).Fatal("Failed to unmarshal yaml in provided fixtures file")
		}

		runServer(*port, &ersatz)
	}

	app.Run(os.Args)
}

func acceptFixtures(w http.ResponseWriter, req *http.Request) {
	yml, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.WithError(err).Error("Failed to read request body")
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	ersatz := ersatz{}
	err = yaml.Unmarshal(yml, &ersatz)
	if err != nil {
		log.WithError(err).Error("Failed to unmarshal yaml")
		http.Error(w, "Failed to unmarshal yaml", http.StatusBadRequest)
		return
	}

	configureOnce.Do(func() {
		configureErsatz(ersatz)
	})

	log.Info("Configured fixtures via /__configure endpoint, further requests to this endpoint will not reconfigure ersatz.")
	w.WriteHeader(http.StatusOK)
}

func runServer(port string, ersatz *ersatz) {
	if ersatz != nil {
		configureErsatz(*ersatz)
	} else {
		log.Info("No fixtures file found, ready to accept fixtures data on POST /__configure")
		http.HandleFunc("/__configure", acceptFixtures)
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Unable to start: %v", err)
	}
}

func configureErsatz(ersatz ersatz) {
	unmonitoredRouter := vestigo.NewRouter()
	var r http.Handler = unmonitoredRouter
	r = httphandlers.TransactionAwareRequestLoggingHandler(log.StandardLogger(), r)

	switch ersatz.Version {
	case "1.0.0-rc1":
	case "1.0.0":
		v1.MockPaths(unmonitoredRouter, ersatz.Fixtures.(*v1.Fixtures))
	case "2.0.0":
		v2.MockPaths(unmonitoredRouter, ersatz.Fixtures.(*v2.Fixtures))
	default:
		log.Fatal(ErrUnsupportedVersion.Error())
	}

	log.Info("Ready to simulate requests!")
	http.Handle("/", r)
}

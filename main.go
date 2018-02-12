package main

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Financial-Times/http-handlers-go/httphandlers"
	"github.com/ghodss/yaml"
	"github.com/husobee/vestigo"
	"github.com/jawher/mow.cli"
	"github.com/peteclark-ft/ersatz/v1"
	metrics "github.com/rcrowley/go-metrics"
	log "github.com/sirupsen/logrus"
)

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
		Value:  "./.ft/fixtures.yml",
		Desc:   "Fixtures file to use to simulate requests",
		EnvVar: "FIXTURES",
	})

	app.Action = func() {
		f, err := ioutil.ReadFile(*fixtures)
		if err != nil {
			log.WithError(err).Fatal("Failed to read fixtures file")
		}

		ersatz := ersatz{}
		err = yaml.Unmarshal(f, &ersatz)

		if err != nil {
			log.WithError(err).Fatal("Failed to unmarshal yml in fixtures file")
		}

		runServer(*port, ersatz)
	}

	app.Run(os.Args)
}

func runServer(port string, ersatz ersatz) {
	unmonitoredRouter := vestigo.NewRouter()
	var r http.Handler = unmonitoredRouter
	r = httphandlers.TransactionAwareRequestLoggingHandler(log.StandardLogger(), r)
	r = httphandlers.HTTPMetricsHandler(metrics.DefaultRegistry, r)

	switch ersatz.Version {
	case 1:
		v1.MockPaths(unmonitoredRouter, ersatz.Fixtures.(*v1.Fixtures))
	default:
		log.Fatal(ErrUnsupportedVersion.Error())
	}

	log.Info("Ready to simulate requests!")
	http.Handle("/", r)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Unable to start: %v", err)
	}
}

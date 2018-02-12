package main

import (
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/jawher/mow.cli"
	"github.com/peteclark-ft/ersatz/v1"
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

		ersatz := ersatz{}
		err = yaml.Unmarshal(f, &ersatz)

		if err != nil {
			log.WithError(err).Error("Failed to marshal yml")
			return
		}

		switch ersatz.Version {
		case 1:
			v1.MockPaths(*port, ersatz.Fixtures.(*v1.Fixtures))
		}
	}

	app.Run(os.Args)
}

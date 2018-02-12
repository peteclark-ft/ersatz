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

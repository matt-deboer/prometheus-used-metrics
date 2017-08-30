package main

import (
	"os"
	"strings"

	"github.com/matt-deboer/mpp/pkg/version"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {

	envBase := strings.ToUpper(version.Name) + "_"

	app := cli.NewApp()
	app.Name = version.Name
	app.Usage = `
		Queries prometheus and Grafana to produce a list of all used metrics
		`
	app.Version = version.Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "prometheus-api-endpoint",
			Usage:  "The prometheus API endpoint to contact",
			EnvVar: envBase + "PROMETHEUS_API_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "grafana-api-endpoint",
			Usage:  "The grafana API endpoint to contact",
			EnvVar: envBase + "GRAFANA_API_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "grafana-credentials",
			Usage:  "The credentials used to authenticate to the grafana API",
			EnvVar: envBase + "GRAFANA_CREDENTIALS",
		},
		cli.BoolFlag{
			Name:   "trace-requests, T",
			Usage:  "Log information about all requests",
			EnvVar: envBase + "TRACE_REQUESTS",
		},
		cli.BoolFlag{
			Name:   "verbose, V",
			Usage:  "Log extra information about steps taken",
			EnvVar: envBase + "VERBOSE",
		},
	}
	app.Action = func(c *cli.Context) {

		if c.Bool("verbose") {
			log.SetLevel(log.DebugLevel)
		}

		prometheusAPI := c.String("prometheus-api-endpoint")
		grafanaAPI := c.String("grafana-api-endpoint")
		grafanaCreds := c.String("grafana-credentials")

		if log.GetLevel() >= log.DebugLevel {
			log.Debugf("Querying for metrics usage: prometheus-api: %s, grafana-api: %s", prometheusAPI, grafanaAPI)
		}

	}
	app.Run(os.Args)
}

package main

import (
	"bufio"
	"io"
	"os"
	"strings"

	"encoding/json"

	"github.com/matt-deboer/prometheus-used-metrics/pkg/query"
	"github.com/matt-deboer/prometheus-used-metrics/pkg/version"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	run(os.Args, os.Stdout)
}

func run(args []string, stdout io.Writer) {
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
		cli.StringFlag{
			Name:   "output-format, o",
			Usage:  "The output format; one of 'json', 'yaml'",
			Value:  "json",
			EnvVar: envBase + "OUTPUT_FORMAT",
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
		if len(prometheusAPI) == 0 {
			argError(c, "'prometheus-api-endpoint' is required")
		}
		grafanaAPI := c.String("grafana-api-endpoint")
		grafanaCreds := c.String("grafana-credentials")

		prometheusResolver, err := query.NewPrometheusMetricsResolver(prometheusAPI)
		if err != nil {
			log.Fatalf("Failed to create prometheus resolver for api '%s'; %v", prometheusAPI, err)
		}
		promAlertUsage, err := prometheusResolver.GetAlertMetricUsage()
		if err != nil {
			log.Fatalf("Failed to resolve metrics usage from prometheus alerts '%s'; %v", prometheusAPI, err)
		}

		usage := make(map[string]map[string][]string)
		for k, v := range promAlertUsage {
			u, found := usage[k]
			if !found {
				u = make(map[string][]string)
				usage[k] = u
			}
			alerts := []string{}
			for vk := range v {
				alerts = append(alerts, vk)
			}
			u["prometheus_alerts"] = alerts
		}

		if len(grafanaAPI) > 0 {
			grafanaResolver := query.NewGrafanaDashboardsResolver(grafanaAPI, grafanaCreds)
			grafanaUsage, err := grafanaResolver.GetMetricUsage()
			if err != nil {
				log.Fatalf("Failed to resolve metrics usage from grafana '%s'; %v", grafanaAPI, err)
			}

			for k, v := range grafanaUsage {
				u, found := usage[k]
				if !found {
					u = make(map[string][]string)
					usage[k] = u
				}
				graphs := []string{}
				for vk := range v {
					graphs = append(graphs, vk)
				}
				u["grafana_graphs"] = graphs
			}
		}

		serialize(usage, stdout)

	}
	app.Run(args)
}

func argError(c *cli.Context, msg string, args ...interface{}) {
	log.Errorf(msg+"\n", args...)
	cli.ShowAppHelp(c)
	os.Exit(1)
}

func serialize(usage interface{}, stdout io.Writer) {
	data, err := json.Marshal(usage)
	if err != nil {
		log.Fatalf("Failed to serialize metrics usage; %v", err)
	}
	w := bufio.NewWriter(stdout)
	_, err = w.Write(data)
	if err != nil {
		log.Fatalf("Failed to write metrics json output; %v", err)
	}
}

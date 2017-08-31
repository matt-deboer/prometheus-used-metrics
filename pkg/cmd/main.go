package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"encoding/json"

	"github.com/ghodss/yaml"
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
			Queries prometheus and Grafana to produce a list of all used metrics;
			
			returns {
				"a_metric_name" : {
					"grafana_graphs" : [ {{ an array of 'slug' values for the graphs using this metric }}],
					"prometheus_alerts" : [ {{ an array of alert names for the alerts using this metric }}],
					"prometheus_jobs" : [ {{ an array of job names for the jobs that provide this metric }}]
				},
				...
			}
			`
	app.Version = version.Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "prometheus-api-endpoint",
			Usage:  "The Prometheus API endpoint to contact",
			EnvVar: envBase + "PROMETHEUS_API_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "grafana-api-endpoint",
			Usage:  "The Grafana API endpoint to contact",
			EnvVar: envBase + "GRAFANA_API_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "grafana-credentials",
			Usage:  "The credentials used to authenticate to the grafana API",
			EnvVar: envBase + "GRAFANA_CREDENTIALS",
		},
		cli.StringFlag{
			Name:   "output-format, o",
			Usage:  "The output format; one of 'json', 'yaml', or 'whitelist'(produces a Prometheus config snippet with metric_relabel_configs, keep entries)",
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

		err = prometheusResolver.ApplyMetricDetails(&usage)
		if err != nil {
			log.Fatalf("Failed to apply metric details using prometheus endpoint '%s'; %v", prometheusAPI, err)
		}

		serialize(usage, stdout, c.String("output-format"))

	}
	app.Run(args)
}

func argError(c *cli.Context, msg string, args ...interface{}) {
	log.Errorf(msg+"\n", args...)
	cli.ShowAppHelp(c)
	os.Exit(1)
}

const metricRelabelTemplate = `
  - job_name: %s	
    metric_relabel_configs:
`
const metricKeepTemplate = `
    - source_labels: [__name__]
      regex: '%s'
      action: keep
`

func serialize(usage map[string]map[string][]string, stdout io.Writer, format string) {
	var data []byte
	var err error
	if format == "whitelist" {
		// Creates a prometheus config yaml snippet for each job that includes the
		// 'metric_relabel_configs' section with a 'keep' entry for each used metric
		metricsByJob := make(map[string][]string)
		for metricName, details := range usage {
			if jobs, found := details["prometheus_jobs"]; found {
				for _, job := range jobs {
					if _, exists := metricsByJob[job]; !exists {
						metricsByJob[job] = []string{}
					}
					metricsByJob[job] = append(metricsByJob[job], metricName)
				}
			}
		}
		buff := bytes.NewBufferString("")
		for job, metrics := range metricsByJob {
			buff.WriteString(fmt.Sprintf(metricRelabelTemplate, job))
			for _, metric := range metrics {
				buff.WriteString(fmt.Sprintf(metricKeepTemplate, metric))
			}
		}
		data = buff.Bytes()
	} else if format == "yaml" {
		data, err = yaml.Marshal(usage)
		if err != nil {
			log.Fatalf("Failed to serialize metrics usage as yaml; %v", err)
		}
	} else {
		data, err = json.Marshal(usage)
		if err != nil {
			log.Fatalf("Failed to serialize metrics usage as json; %v", err)
		}
	}
	w := bufio.NewWriter(stdout)
	_, err = w.Write(data)
	if err != nil {
		log.Fatalf("Failed to write metrics output; %v", err)
	}
}

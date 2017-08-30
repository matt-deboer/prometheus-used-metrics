prometheus-used-metrics
===

Produces a manifest of metrics used across one or more sources (of which Grafana and Prometheus Rules are currently supported).

Motivation
---

The main drive for this tool is to allow creation of a whitelist of metrics to keep for a given Prometheus installation. Often,
Prometheus is sraping and storing so much more information than is actually needed--this allows a way to programmatically generate
a whitelist that can be used to limit Prometheus' retained metrics to the bare minimum needed to support current usage.

Usage
---

```
NAME:
   prometheus-used-metrics -
      Queries prometheus and Grafana to produce a list of all used metrics


USAGE:
   prometheus-used-metrics [global options] command [command options] [arguments...]

VERSION:
   b9dcd11+local_changes

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --prometheus-api-endpoint value  The prometheus API endpoint to contact [$PROMETHEUS-USED-METRICS_PROMETHEUS_API_ENDPOINT]
   --grafana-api-endpoint value     The grafana API endpoint to contact [$PROMETHEUS-USED-METRICS_GRAFANA_API_ENDPOINT]
   --grafana-credentials value      The credentials used to authenticate to the grafana API [$PROMETHEUS-USED-METRICS_GRAFANA_CREDENTIALS]
   --output-format value, -o value  The output format; one of 'json', 'yaml' (default: "json") [$PROMETHEUS-USED-METRICS_OUTPUT_FORMAT]
   --trace-requests, -T             Log information about all requests [$PROMETHEUS-USED-METRICS_TRACE_REQUESTS]
   --verbose, -V                    Log extra information about steps taken [$PROMETHEUS-USED-METRICS_VERBOSE]
   --help, -h                       show help
   --version, -v                    print the version
```


Example usage 

```sh
prometheus-used-metrics --prometheus-api-endpoint http://prometheus.example.org --grafana-api-endpoint http://grafana.example.org --grafana-credentials "reallylongapitokenhere"
```

Output looks like:

```json
{
  "metric_name_1": {
    "grafana_graphs": ["Graph1","Graph2"],
    "prometheus_alerts": ["alert1","another_alert"],
  }
}
```
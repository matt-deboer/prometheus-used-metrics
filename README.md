prometheus-used-metrics
===
[![Build Status](https://travis-ci.org/matt-deboer/prometheus-used-metrics.svg?branch=master)](https://travis-ci.org/matt-deboer/prometheus-used-metrics)
[![Docker Pulls](https://img.shields.io/docker/pulls/mattdeboer/prometheus-used-metrics.svg)](https://hub.docker.com/r/mattdeboer/prometheus-used-metrics/)
[![Coverage Status](https://coveralls.io/repos/github/matt-deboer/prometheus-used-metrics/badge.svg?branch=master)](https://coveralls.io/github/matt-deboer/prometheus-used-metrics?branch=master) 



Produces a manifest of metrics used across one or more sources (of which Grafana and Prometheus Rules are currently supported).

Motivation
---

Often, Prometheus is sraping and storing so much more information than is actually needed--this tool allows you to programmatically generate
a whitelist that can be used to limit Prometheus' retained metrics to the bare minimum needed to support your current usage.

Usage
---

```
NAME:
   prometheus-used-metrics -
      Queries prometheus and Grafana to produce a list of all used metrics;

      returns {
        "a_metric_name" : {
          "grafana_graphs" : [ {{ an array of 'slug' values for the graphs using this metric }}],
          "prometheus_alerts" : [ {{ an array of alert names for the alerts using this metric }}],
          "prometheus_jobs" : [ {{ an array of job names for the jobs that provide this metric }}]
        },
        ...
      }


USAGE:
   prometheus-used-metrics [global options] command [command options] [arguments...]

VERSION:
   01ca8b6+local_changes

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --prometheus-api-endpoint value  The Prometheus API endpoint to contact [$PROMETHEUS-USED-METRICS_PROMETHEUS_API_ENDPOINT]
   --grafana-api-endpoint value     The Grafana API endpoint to contact [$PROMETHEUS-USED-METRICS_GRAFANA_API_ENDPOINT]
   --grafana-credentials value      The credentials used to authenticate to the grafana API [$PROMETHEUS-USED-METRICS_GRAFANA_CREDENTIALS]
   --output-format value, -o value  The output format; one of 'json', 'yaml', or 'whitelist'(produces a Prometheus config snippet with metric_relabel_configs, keep entries) (default: "json") [$PROMETHEUS-USED-METRICS_OUTPUT_FORMAT]
   --trace-requests, -T             Log information about all requests [$PROMETHEUS-USED-METRICS_TRACE_REQUESTS]
   --verbose, -V                    Log extra information about steps taken [$PROMETHEUS-USED-METRICS_VERBOSE]
   --help, -h                       show help
   --version, -v                    print the version
```


## Example usage 

```sh
prometheus-used-metrics --prometheus-api-endpoint http://prometheus.example.org --grafana-api-endpoint http://grafana.example.org --grafana-credentials "reallylongapitokenhere"
```

Output looks like:

```json
{
  "metric_name_1": {
    "grafana_graphs": ["Graph1","Graph2"],
    "prometheus_alerts": ["alert1","another_alert"],
    "prometheus_jobs": ["job1","job2"]
  },
  ...
}
```

## Output Format Examples

- `json`

  ```json
  {
    "mesos_slave_cpus": {
      "grafana_graphs": [
        "db/resource-availability"
      ],
      "prometheus_alerts": [
        "system_low_cpu_in_rack",
        "system_low_cpu_in_az"
      ],
      "prometheus_jobs": [
        "marathon_exporter",
        "mesos_exporter_agents"
      ]
    },
    "mesos_slave_mem": {
      "grafana_graphs": [
        "db/resource-availability"
      ],
      "prometheus_alerts": [
        "system_low_mem_in_rack",
        "system_low_mem_in_az"
      ],
      "prometheus_jobs": [
        "mesos_exporter_agents"
      ]
    },
    ...
  }
  ```

- `yaml`

  ```yaml
  mesos_slave_cpus:
    grafana_graphs:
    - db/resource-availability
    prometheus_alerts:
    - system_low_cpu_in_rack
    - system_low_cpu_in_az
    prometheus_jobs:
    - marathon_exporter
    - mesos_exporter_agents
  mesos_slave_mem:
    grafana_graphs:
    - db/resource-availability
    prometheus_alerts:
    - system_low_mem_in_az
    - system_low_mem_in_rack
    prometheus_jobs:
    - mesos_exporter_agents
  ```

- `whitelist`

  ```yaml
  - job_name: node_exporter_agents
    metric_relabel_configs:

    - source_labels: [__name__]
      regex: '(node_filesystem_size|node_filesystem_free|node_network_transmit_bytes|node_network_receive_bytes|node_cpu)'
      action: keep

  - job_name: node_exporter_masters
    metric_relabel_configs:

    - source_labels: [__name__]
      regex: '(node_filesystem_size|node_filesystem_free|node_filesystem_sizenode_network_transmit_bytes|node_network_receive_bytes)'
      action: keep
  ```

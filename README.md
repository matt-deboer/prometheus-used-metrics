prometheus-used-metrics
===

Produces a manifest of metrics used across one or more sources (of which Grafana and Prometheus Rules are currently supported).

Motivation
---

The main drive for this tool is to allow creation of a whitelist of metrics to keep for a given Prometheus installation. Often,
Prometheus is sraping and storing so much more information than is actually needed--this allows a way to programmatically generate
a whitelist that can be used to limit Prometheus' retained metrics to the bare minimum needed to support current usage.

Example usage should be:

```sh
prometheus-used-metrics --prometheus-api-endpoint http://somewhere.com --grafana-api-endpoint http://grafana.something --grafana-credentials "Blah3asdfashdasdfnafsdfasdf"
```

Output should be:

```json
{
  "metric_name_1": {
    "grafana_graphs": ["Graph1","Graph2"],
    "prometheus_alerts": ["alert1","another_alert"],
  }
}
```
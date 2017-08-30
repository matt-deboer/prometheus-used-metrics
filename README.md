prometheus-used-metrics
===

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
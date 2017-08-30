changelog
===

v0.2.2 [2017-07-06]
---

**features:**

- added viability check to endpoints;  in addition to the scraped comparision metric value,
  endpoints must be able to respond to a query for the `up` metric in < 1 sec.

**fixes:**

- updated `random` selector strategy to use direct scrape of `prometheus_build_info`--to
  be consistent with `single_most_data` strategy, and to allow usage in scenarios where
  prometheus does not scrape its own metrics.

v0.2.0 [2017-04-28]
---

**features:**

- added a `-alpine` release of the docker image

**fixes:**

- fixed inconsistencies in labels and env. variables
- fixed prometheus-specific metrics lookup to use prometheus endpoints, rather than query--since
  prometheus self-scraping may not be configured (at least initially)

v0.2.0-a2 [2017-04-27]
---

**fixes:**

- added `--kube-namespace` requirement for kubernetes deployments to fix bug in k8s locator

v0.2.0-a1 [2017-04-27]
---

**features:**

- initial kubernetes support
- added metrics for
  - build info
  - affinity hits by type
  - selection events
  - requests by backend
  - repsonse time by backend
- added selector strategies:
  - `minimum-history`
  - `random`
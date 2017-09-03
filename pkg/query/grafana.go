package query

import (
	"fmt"
	"regexp"

	"github.com/matt-deboer/sdk"
	log "github.com/sirupsen/logrus"

	"github.com/prometheus/prometheus/promql"
)

// GrafanaDashboardsResolver resolves dashboards
type GrafanaDashboardsResolver struct {
	client *sdk.Client
}

// NewGrafanaDashboardsResolver creates a new GrafanaDashboardsResolver
// for the provided endpoint, authenticating using the provided credentials
func NewGrafanaDashboardsResolver(grafanaAPIEndpoint, grafanaCredentials string) *GrafanaDashboardsResolver {
	return &GrafanaDashboardsResolver{client: sdk.NewClient(grafanaAPIEndpoint, grafanaCredentials, sdk.DefaultHTTPClient)}
}

// GetMetricUsage returns a map of metrics to arrays of dashboard names in which they are used
func (g *GrafanaDashboardsResolver) GetMetricUsage() (map[string]map[string]string, error) {
	boards, err := g.client.SearchDashboards("", false)
	if err != nil {
		return nil, fmt.Errorf("Failed to resolve dashboards; %v", err)
	}

	metricsUsage := make(map[string]map[string]string)
	for _, b := range boards {
		board, _, err := g.client.GetDashboard(b.URI)
		if err != nil {
			return nil, fmt.Errorf("Failed to resolve dashboard '%s'; %v", b.URI, err)
		}
		if log.GetLevel() >= log.DebugLevel {
			log.Debugf("Parsing dashboard '%s'...", board.Title)
		}
		metricsUsed := []string{}
		// parse strings in the board looking for prometheus metrics
		for _, row := range board.Rows {
			for _, panel := range row.Panels {
				targets := panel.GetTargets()
				if targets != nil {
					for _, target := range *targets {
						// This is our check for prometheus => this value is only filled in for prometheus targets
						if len(target.Expr) > 0 {
							if log.GetLevel() >= log.DebugLevel {
								log.Debugf("Dashboard '%s' has Prometheus expression '%s'; parsing used metrics...", board.Title, target.Expr)
							}
							metrics, err := parseUsedMetrics(target.Expr)
							if err != nil {
								return nil, fmt.Errorf("Failed to parse expression for dashboard '%s'; expression '%s'; %v", board.Title, target.Expr, err)
							}
							metricsUsed = append(metricsUsed, metrics...)
						}
					}
				}
			}
		}

		for _, templateVar := range board.Templating.List {
			if templateVar.Datasource != nil {
				if log.GetLevel() >= log.DebugLevel {
					log.Debugf("Dashboard '%s' has template var '%s' with query '%s'; parsing used metrics...", board.Title, templateVar.Name, templateVar.Query)
				}
				metrics, err := parseUsedMetrics(templateVar.Query)
				if err != nil {
					log.Warnf("Failed to parse expression for template var '%s'; expression '%s'; %v", templateVar.Name, templateVar.Query, err)
				} else {
					metricsUsed = append(metricsUsed, metrics...)
				}
			}
		}

		if log.GetLevel() >= log.DebugLevel {
			log.Debugf("Dashboard '%s' requires metrics %#v", board.Title, metricsUsed)
		}

		for _, metric := range metricsUsed {
			usage, found := metricsUsage[metric]
			if !found {
				usage = make(map[string]string)
				metricsUsage[metric] = usage
			}
			usage[b.URI] = b.Title
		}
	}
	return metricsUsage, nil
}

func parseUsedMetrics(expression string) ([]string, error) {

	// replace any variable references with plain names
	rexp := regexp.MustCompile(`\$([A-Za-z_]+)`)
	sanitized := rexp.ReplaceAllString(expression, "$1")

	expr, err := promql.ParseExpr(sanitized)
	if err != nil {
		return nil, err
	}
	v := &visitor{}
	promql.Walk(v, expr)
	return v.metrics, err
}

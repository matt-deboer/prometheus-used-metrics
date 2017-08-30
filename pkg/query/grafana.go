package query

import (
	"fmt"

	"github.com/grafana-tools/sdk"
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

// GetDashboards returns a map
func (g *GrafanaDashboardsResolver) GetDashboards() (map[string][]string, error) {
	boards, err := g.client.SearchDashboards("*", false)
	if err != nil {
		return nil, fmt.Errorf("Failed to resolve dashboards; %v", err)
	}

	metricsUsage := make(map[string][]string)
	for _, b := range boards {
		board, props, err := g.client.GetDashboard(b.Title)
		if err != nil {
			return nil, fmt.Errorf("Failed to resolve dashboard '%s'; %v", b.Title, err)
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
							metricsUsed = append(metricsUsed, parseUsedMetrics(target.Expr)...)
						}
					}
				}
			}
		}
		for _, metric := range metricsUsed {
			metricUsage, found := metricsUsage[metric]
			if !found {
				metricUsage = []string{}
				metricsUsage[metric] = metricUsage
			}
			metricUsage = append(metricUsage, b.Title)
		}
	}
	return metricsUsage, nil
}

func parseUsedMetrics(expression string) []string {

}

package query

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/api/prometheus"
	"github.com/prometheus/prometheus/promql"
	log "github.com/sirupsen/logrus"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var alertSplitter = regexp.MustCompile(` ALERT `)

type PrometheusMetricsResolver struct {
	queryAPI prometheus.QueryAPI
	endpoint string
}

func NewPrometheusMetricsResolver(prometheusAPIEndpoint string) (*PrometheusMetricsResolver, error) {
	client, err := prometheus.New(prometheus.Config{
		Address: prometheusAPIEndpoint,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create prometheus client for endpoint '%s'; %v", prometheusAPIEndpoint, err)
	}
	return &PrometheusMetricsResolver{
		queryAPI: prometheus.NewQueryAPI(client),
		endpoint: prometheusAPIEndpoint,
	}, nil
}

// GetAlertMetricUsage returns a map of metric names to alerts where they are used
func (p *PrometheusMetricsResolver) GetAlertMetricUsage() (map[string]map[string]bool, error) {
	resp, err := http.Get(fmt.Sprintf("%s/rules", p.endpoint))
	if err != nil {
		panic(err)
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	metricsUsage := make(map[string]map[string]bool)
	// define a matcher
	matcher := func(n *html.Node) bool {
		// must check for nil values
		return n.DataAtom == atom.Pre
	}
	// grab all articles and print them
	for _, rules := range scrape.FindAll(root, matcher) {
		text := scrape.Text(rules)
		rawAlerts := alertSplitter.Split(text, -1)
		for _, rawAlert := range rawAlerts {
			if !strings.HasPrefix(rawAlert, "ALERT ") {
				rawAlert = "ALERT " + rawAlert
			}
			if log.GetLevel() >= log.DebugLevel {
				log.Debugf("Parsing alert:\n%s\n", rawAlert)
			}
			statements, err := promql.ParseStmts(rawAlert)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse alert '%s'; %v", rawAlert, err)
			}
			if len(statements) != 1 {
				return nil, fmt.Errorf("Expected a single statement for alert '%s'; got %d", rawAlert, len(statements))
			}
			v := &visitor{}
			promql.Walk(v, statements[0])
			for _, metric := range v.metrics {
				usage, found := metricsUsage[metric]
				if !found {
					usage = make(map[string]bool)
					metricsUsage[metric] = usage
				}
				usage[v.alertName] = true
			}
		}
	}
	return metricsUsage, nil
}

func (p *PrometheusMetricsResolver) GetMetricDetails(metricName string) (map[string]string, error) {

}

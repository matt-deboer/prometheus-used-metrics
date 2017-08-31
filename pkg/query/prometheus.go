package query

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/prometheus/client_golang/api/prometheus"
	"github.com/prometheus/common/model"
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

const nameLabel = model.LabelName("__name__")
const jobLabel = model.LabelName("job")

func (p *PrometheusMetricsResolver) ApplyMetricDetails(metrics *map[string]map[string][]string) error {

	keys := make([]string, len(*metrics))
	i := 0
	for k := range *metrics {
		keys[i] = k
		i++
	}
	names := strings.Join(keys, "|")
	result, err := p.queryAPI.Query(context.Background(), fmt.Sprintf(`count({__name__=~"(%s)"}) by (__name__,job)`, names), time.Now())
	if err != nil {
		return fmt.Errorf("Failed to retrieve jobs for metrics '%s'", keys)
	}

	switch result.Type() {
	case model.ValVector:
		for _, sample := range result.(model.Vector) {
			if log.GetLevel() >= log.DebugLevel {
				log.Debugf("Parsing Vector sample: %s...", sample.String)
			}
			labels := model.LabelSet(sample.Metric)
			name := string(labels[nameLabel])
			job := string(labels[jobLabel])
			metric := (*metrics)[name]
			if _, found := metric["prometheus_jobs"]; !found {
				metric["prometheus_jobs"] = []string{}
			}
			metric["prometheus_jobs"] = append(metric["prometheus_jobs"], job)
		}
		break
	}
	return nil
}

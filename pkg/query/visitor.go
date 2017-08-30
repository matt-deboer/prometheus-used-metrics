package query

import (
	"github.com/prometheus/prometheus/promql"
)

type visitor struct {
	metrics   []string
	alertName string
}

func (v *visitor) Visit(node promql.Node) promql.Visitor {
	switch node.(type) {
	case *promql.MatrixSelector:
		v.metrics = append(v.metrics, node.(*promql.MatrixSelector).Name)
		break
	case *promql.VectorSelector:
		v.metrics = append(v.metrics, node.(*promql.VectorSelector).Name)
		break
	case *promql.AlertStmt:
		v.alertName = node.(*promql.AlertStmt).Name
		break
	}
	return v
}

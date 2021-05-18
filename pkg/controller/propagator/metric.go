package propagator

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	roothandlerMeasure = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "func_handle_root_policy_elapsed_seconds",
		Help: "how long do the handleRootPolicy take to complete",
	})
)

func init() {
	metrics.Registry.MustRegister(roothandlerMeasure)
}

package metrics

import (
	"context"
	"fmt"
	"sync"
	"time"

	gokitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const metricsContextKey = "metrics_context"

// The MetricsContext has the metrics publishers to prometheus.
type MetricsContext struct {
	Counters   map[string]metrics.Counter
	Histograms map[string]metrics.Histogram
	Namespace  string
	Logger     gokitlog.Logger
	mu         sync.Mutex
}

/** The New returns a MetricsContext that implements the MetricsContext interface.
 * When no logger is provided, it defaults to go-kit's NewNopLogger
 */
func New(logger gokitlog.Logger, namespace string) *MetricsContext {
	if logger == nil {
		logger = gokitlog.NewNopLogger()
	}

	return &MetricsContext{
		Namespace:  namespace,
		Counters:   make(map[string]metrics.Counter),
		Histograms: make(map[string]metrics.Histogram),
		Logger:     logger,
	}
}

/* The IntoContext injects a MetricsContext into the context.Context which can be extract
 * with metrics.FromContext
 */
func IntoContext(ctx context.Context, m *MetricsContext) context.Context {
	return context.WithValue(ctx, metricsContextKey, m)
}

/* The FromContext extract a MetricsContext singleton from the context.Context if it was previously
 * inject using metrics.IntoContext
 */
func FromContext(ctx context.Context) (*MetricsContext, error) {
	value := ctx.Value(metricsContextKey)
	if value == nil {
		return nil, fmt.Errorf("no MetricsContext configured or present in the context")
	}
	metricsCtx, ok := value.(*MetricsContext)
	if !ok {
		return nil, fmt.Errorf("MetricsContext is misconfigured and cannot be extracted")
	}
	return metricsCtx, nil
}

/** The Record function emits a metric the provided name and list of labels and values. Even positions
 * on the list represent the labels and the remains the values. Any error that occurs is recovered and
 * logged to avoid interrupting the main computation
 */
func Record(ctx context.Context, name string, count float64, labelValues ...string) error {
	metrics, err := FromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to extract the MetricsContext from context.Context")
	}

	return metrics.Record(name, count, labelValues...)
}

func (c *MetricsContext) Record(name string, count float64, labelValues ...string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	defer c.recoverFromPanic(name)

	metric, ok := c.Counters[name]
	if !ok {
		ops := stdprometheus.CounterOpts{
			Namespace: c.Namespace,
			Name:      name,
			Help:      name + " is a counter",
		}

		metric = prometheus.NewCounterFrom(ops, getLabels(labelValues))
		c.Counters[name] = metric
	}

	metric.With(labelValues...).Add(count)
	return nil
}

/** The RecordTiming function emits a metric using the provided name, duration value, and list of labels and
 * values. Even positions on the list represent the labels and the remains the values. Any error that occurs
 * is recovered and logged to avoid interrupting the main computation
 */
func RecordTiming(ctx context.Context, name string, duration time.Duration, labelValues ...string) error {
	metrics, err := FromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to extract the MetricsContext from context.Context")
	}

	return metrics.RecordTiming(name, duration, labelValues...)
}

func (c *MetricsContext) RecordTiming(name string, duration time.Duration, labelValues ...string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	defer c.recoverFromPanic(name)

	metric, ok := c.Histograms[name]
	if !ok {
		ops := stdprometheus.HistogramOpts{
			Namespace: c.Namespace,
			Name:      name,
			Help:      name + " is a histogram",
		}

		metric = prometheus.NewHistogramFrom(ops, getLabels(labelValues))
		c.Histograms[name] = metric
	}

	metric.With(labelValues...).Observe(duration.Seconds())
	return nil
}

type void struct{}

func (c *MetricsContext) CounterNames() map[string]void {
	c.mu.Lock()
	defer c.mu.Unlock()

	names := make(map[string]void, len(c.Counters))
	for k := range c.Counters {
		names[k] = void{}
	}
	return names
}

func (c *MetricsContext) TimingNames() map[string]void {
	c.mu.Lock()
	defer c.mu.Unlock()

	names := make(map[string]void, len(c.Histograms))
	for k := range c.Histograms {
		names[k] = void{}
	}
	return names
}

func getLabels(lvs []string) []string {
	length := len(lvs)
	if length < 1 {
		return []string{}
	}

	keyCounter := 0
	totalKeys := len(lvs) / 2
	keys := make([]string, totalKeys)
	for idx, label := range lvs {
		if idx%2 == 0 && keyCounter < totalKeys {
			keys[keyCounter] = label
			keyCounter += 1
		}
	}

	return keys
}

func (c *MetricsContext) recoverFromPanic(name string) {
	if r := recover(); r != nil {
		c.Logger.Log("metrics", name, "kind", "duration", "err", r)
	}
}

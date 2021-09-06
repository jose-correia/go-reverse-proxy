// +build unit

package metrics_test

import (
	"crypto/rand"
	"go-reverse-proxy/app/common/log"
	"go-reverse-proxy/app/common/metrics"
	"testing"
	"time"

	gokitmetrics "github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestFromContextOnNonInitializedMetricContext(t *testing.T) {
	baseCtx := context.Background()
	metricsCtx, err := metrics.FromContext(baseCtx)
	assert.Nil(t, metricsCtx)
	assert.Error(t, err, "")
}

func TestFromContextOnInitializedMetricContext(t *testing.T) {
	mtr := &metrics.MetricsContext{
		Namespace:  "reverseproxy",
		Counters:   make(map[string]gokitmetrics.Counter),
		Histograms: make(map[string]gokitmetrics.Histogram),
		Logger:     log.NewLogger(),
	}
	ctxWithMetrics := metrics.IntoContext(context.Background(), mtr)
	metricsCtx, err := metrics.FromContext(ctxWithMetrics)
	assert.Nil(t, err)
	assert.Equal(t, metricsCtx, mtr)
}

func TestNestedContextWrapping(t *testing.T) {
	mtr := &metrics.MetricsContext{
		Namespace:  "reverseproxy",
		Counters:   make(map[string]gokitmetrics.Counter),
		Histograms: make(map[string]gokitmetrics.Histogram),
		Logger:     log.NewLogger(),
	}

	ctxWithMetrics := metrics.IntoContext(context.Background(), mtr)
	wrappedCtx := metrics.IntoContext(ctxWithMetrics, mtr)
	metricsCtx, err := metrics.FromContext(wrappedCtx)
	assert.Nil(t, err)
	assert.NotNil(t, metricsCtx)
}

func TestRecordWhenMetricDoesntExist(t *testing.T) {
	metricName := RandStringBytes(t, 10)
	lvs := []string{"label_1", "1", "label_2", "2"}
	counters := make(map[string]gokitmetrics.Counter)

	metricsCtx := &metrics.MetricsContext{
		Namespace:  "reverseproxy",
		Counters:   counters,
		Histograms: make(map[string]gokitmetrics.Histogram),
	}

	assert.Nil(t, metricsCtx.Counters[metricName])
	ctx := metrics.IntoContext(context.Background(), metricsCtx)
	err := metrics.Record(ctx, metricName, 1, lvs...)
	assert.Nil(t, err)
	assert.NotNil(t, metricsCtx.Counters[metricName])
}

func TestRecordWhenMetricExists(t *testing.T) {
	metricName := RandStringBytes(t, 10)
	keys := []string{"label_1", "label_2"}
	lvs := []string{"label_1", "1", "label_2", "2"}
	counters := make(map[string]gokitmetrics.Counter)
	ops := stdprometheus.CounterOpts{
		Namespace: "reverseproxy",
		Name:      metricName,
		Help:      metricName,
	}
	counters[metricName] = prometheus.NewCounterFrom(ops, keys)

	metricsCtx := &metrics.MetricsContext{
		Namespace:  "reverseproxy",
		Counters:   counters,
		Histograms: make(map[string]gokitmetrics.Histogram),
	}
	assert.NotNil(t, metricsCtx.Counters[metricName])
	ctx := metrics.IntoContext(context.Background(), metricsCtx)
	err := metrics.Record(ctx, metricName, 1, lvs...)
	assert.Nil(t, err)
	assert.NotNil(t, metricsCtx.Counters[metricName])
}

func TestRecordWhenErrorOccursDoNotCrash(t *testing.T) {
	metricName := RandStringBytes(t, 10)
	keys := []string{"label_1", "label_2"}
	lvs := []string{"label_1", "1", "label_2", "2", "label_3"}
	counters := make(map[string]gokitmetrics.Counter)
	ops := stdprometheus.CounterOpts{
		Namespace: "reverseproxy",
		Name:      metricName,
		Help:      metricName,
	}
	counters[metricName] = prometheus.NewCounterFrom(ops, keys)

	metricsCtx := &metrics.MetricsContext{
		Namespace:  "reverseproxy",
		Counters:   counters,
		Histograms: make(map[string]gokitmetrics.Histogram),
		Logger:     log.NewLogger(),
	}
	assert.NotNil(t, metricsCtx.Counters[metricName])
	ctx := metrics.IntoContext(context.Background(), metricsCtx)
	err := metrics.Record(ctx, metricName, 1, lvs...)
	assert.Nil(t, err)
	assert.NotNil(t, metricsCtx.Counters[metricName])
}

func RecordTimingWhenMetricDoesntExist(t *testing.T) {
	metricName := RandStringBytes(t, 10)
	lvs := []string{"label_1", "1", "label_2", "2"}
	histograms := make(map[string]gokitmetrics.Histogram)
	metricsCtx := &metrics.MetricsContext{
		Namespace:  "reverseproxy",
		Counters:   make(map[string]gokitmetrics.Counter),
		Histograms: histograms,
	}

	assert.Nil(t, metricsCtx.Histograms[metricName])
	ctx := metrics.IntoContext(context.Background(), metricsCtx)
	err := metrics.RecordTiming(ctx, metricName, 1*time.Second, lvs...)
	assert.Nil(t, err)
	assert.NotNil(t, metricsCtx.Histograms[metricName])
}

func RecordTimingWhenMetricExists(t *testing.T) {
	metricName := RandStringBytes(t, 10)
	keys := []string{"label_1", "label_2"}
	lvs := []string{"label_1", "1", "label_2", "2"}
	histograms := make(map[string]gokitmetrics.Histogram)
	ops := stdprometheus.HistogramOpts{
		Namespace: "reverseproxy",
		Name:      metricName,
		Help:      metricName,
	}
	histograms[metricName] = prometheus.NewHistogramFrom(ops, keys)

	metricsCtx := &metrics.MetricsContext{
		Namespace:  "reverseproxy",
		Counters:   make(map[string]gokitmetrics.Counter),
		Histograms: histograms,
		Logger:     log.NewLogger(),
	}

	assert.NotNil(t, metricsCtx.Histograms[metricName])
	ctx := metrics.IntoContext(context.Background(), metricsCtx)
	err := metrics.RecordTiming(ctx, metricName, 1*time.Second, lvs...)
	assert.Nil(t, err)
	assert.NotNil(t, metricsCtx.Histograms[metricName])
}

func RecordTimingWhenErrorOccursDoNotCrash(t *testing.T) {
	metricName := RandStringBytes(t, 10)
	keys := []string{"label_1", "label_2"}
	lvs := []string{"label_1", "1", "label_2", "2", "label_3"}
	histograms := make(map[string]gokitmetrics.Histogram)
	ops := stdprometheus.HistogramOpts{
		Namespace: "reverseproxy",
		Name:      metricName,
		Help:      metricName,
	}
	histograms[metricName] = prometheus.NewHistogramFrom(ops, keys)

	metricsCtx := &metrics.MetricsContext{
		Namespace:  "reverseproxy",
		Counters:   make(map[string]gokitmetrics.Counter),
		Histograms: histograms,
		Logger:     log.NewLogger(),
	}

	assert.NotNil(t, metricsCtx.Histograms[metricName])
	ctx := metrics.IntoContext(context.Background(), metricsCtx)
	err := metrics.RecordTiming(ctx, metricName, 1*time.Second, lvs...)
	assert.Nil(t, err)
	assert.NotNil(t, metricsCtx.Histograms[metricName])
}

func RandStringBytes(t *testing.T, n int) string {
	value, err := randStringBytes(n)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	return value
}

func randStringBytes(n int) (string, error) {
	const alphanum = "abcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes), nil
}

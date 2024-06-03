package cmetrics

import (
	"github.com/gocopper/copper/cerrors"
	"github.com/gocopper/copper/clogger"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics interface {
	CounterInc(name string, labels map[string]string)
	HistogramObserve(name string, labels map[string]string, value float64)
}

type metrics struct {
	counters   map[string]*prometheus.CounterVec
	histograms map[string]*prometheus.HistogramVec

	logger clogger.Logger
}

func NewMetrics(registry *Registry, logger clogger.Logger) (Metrics, error) {
	countersByName := make(map[string]*prometheus.CounterVec)
	for i := range registry.Counters {
		countersByName[registry.Counters[i].Name] = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: registry.Counters[i].Name,
		}, registry.Counters[i].Labels)

		err := prometheus.DefaultRegisterer.Register(countersByName[registry.Counters[i].Name])
		if err != nil {
			return nil, cerrors.New(err, "failed to register counter metric", map[string]interface{}{
				"name": registry.Counters[i].Name,
			})
		}
	}

	histogramsByName := make(map[string]*prometheus.HistogramVec)
	for i := range registry.Histograms {
		histogramsByName[registry.Histograms[i].Name] = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    registry.Histograms[i].Name,
			Buckets: registry.Histograms[i].Buckets,
		}, registry.Histograms[i].Labels)

		err := prometheus.DefaultRegisterer.Register(histogramsByName[registry.Histograms[i].Name])
		if err != nil {
			return nil, cerrors.New(err, "failed to register histogram metric", map[string]interface{}{
				"name": registry.Histograms[i].Name,
			})

		}
	}

	return &metrics{
		counters:   countersByName,
		histograms: histogramsByName,

		logger: logger,
	}, nil
}

func (m *metrics) HistogramObserve(name string, labels map[string]string, value float64) {
	histogram, ok := m.histograms[name]
	if !ok {
		m.logger.WithTags(map[string]interface{}{
			"name": name,
		}).Warn("Histogram is not registered. Ignoring..", nil)
		return
	}

	metric, err := histogram.GetMetricWith(labels)
	if err != nil {
		m.logger.WithTags(map[string]interface{}{
			"name": name,
		}).Warn("Failed to get histogram metric with labels", err)
		return

	}

	metric.Observe(value)
}

func (m *metrics) CounterInc(name string, labels map[string]string) {
	counter, ok := m.counters[name]
	if !ok {
		m.logger.WithTags(map[string]interface{}{
			"name": name,
		}).Warn("Counter is not registered. Ignoring..", nil)
		return
	}

	metric, err := counter.GetMetricWith(labels)
	if err != nil {
		m.logger.WithTags(map[string]interface{}{
			"name": name,
		}).Warn("Failed to get counter metric with labels", err)
		return
	}

	metric.Inc()
}

package cmetrics

func NewNoopMetrics() Metrics {
	return &noop{}
}

type noop struct{}

func (m *noop) CounterInc(name string, labels map[string]string) {
}

func (m *noop) HistogramObserve(name string, labels map[string]string, value float64) {
}

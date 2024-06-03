package cmetrics

type (
	Registry struct {
		Counters   []Counter
		Histograms []Histogram
	}

	Counter struct {
		Name   string
		Labels []string
	}

	Histogram struct {
		Name    string
		Labels  []string
		Buckets []float64
	}
)

type NewRegistryParams struct {
	Counters   []Counter
	Histograms []Histogram
}

var (
	internalCounters = []Counter{
		{
			Name:   "http_requests_total",
			Labels: []string{"status_code", "path"},
		},
	}

	internalHistograms = []Histogram{
		{
			Name:    "http_request_duration_seconds",
			Labels:  []string{"status_code", "path"},
			Buckets: []float64{0.1, 0.2, 0.5, 1.0, 2.0, 5.0, 10.0},
		},
	}
)

func NewRegistry(p NewRegistryParams) *Registry {
	return &Registry{
		Counters:   append(p.Counters, internalCounters...),
		Histograms: append(p.Histograms, internalHistograms...),
	}
}

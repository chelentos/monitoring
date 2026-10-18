package handler

type MetricsStorage interface {
	GetGauge(name string) (float64, bool)
	GetGauges() map[string]float64
	SetGauge(name string, value float64)
	GetCounter(name string) (int64, bool)
	GetCounters() map[string]int64
	AddCounter(name string, delta int64)
}

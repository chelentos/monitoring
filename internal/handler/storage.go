package handler

type MetricsStorage interface {
	SetGauge(name string, value float64)
	AddCounter(name string, delta int64)
}

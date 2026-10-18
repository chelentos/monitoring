package handler

type MetricsStorage interface {
	GetGauge(name string) (float64, bool)
	SetGauge(name string, value float64)
	GetCounter(name string) (int64, bool)
	AddCounter(name string, delta int64)
}

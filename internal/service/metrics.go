package service

import "errors"

var ErrEmptyName = errors.New("metric name is empty")
var ErrMetricNotFound = errors.New("metric not found")

type MetricsStorage interface {
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
	GetGauges() map[string]float64
	GetCounters() map[string]int64
	SetGauge(name string, value float64)
	AddCounter(name string, delta int64)
}

type MetricsService struct {
	storage MetricsStorage
}

func NewMetricsService(storage MetricsStorage) *MetricsService {
	return &MetricsService{storage: storage}
}

func (s *MetricsService) GetGauge(name string) (float64, error) {
	if name == "" {
		return 0, ErrEmptyName
	}
	value, ok := s.storage.GetGauge(name)
	if !ok {
		return 0, ErrMetricNotFound
	}
	return value, nil
}

func (s *MetricsService) GetCounter(name string) (int64, error) {
	if name == "" {
		return 0, ErrEmptyName
	}
	value, ok := s.storage.GetCounter(name)
	if !ok {
		return 0, ErrMetricNotFound
	}
	return value, nil
}

func (s *MetricsService) GetGauges() map[string]float64 {
	return s.storage.GetGauges()
}

func (s *MetricsService) GetCounters() map[string]int64 {
	return s.storage.GetCounters()
}

func (s *MetricsService) UpdateGauge(name string, value float64) error {
	if name == "" {
		return ErrEmptyName
	}

	s.storage.SetGauge(name, value)

	return nil
}

func (s *MetricsService) UpdateCounter(name string, delta int64) error {
	if name == "" {
		return ErrEmptyName
	}

	s.storage.AddCounter(name, delta)

	return nil
}

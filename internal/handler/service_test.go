package handler

import "github.com/chelentos/monitoring/internal/service"

type serviceStub struct {
	gauges                   map[string]float64
	counters                 map[string]int64
	err                      error
	gaugeCalls, counterCalls int
	name                     string
	value                    float64
	delta                    int64
}

func (s *serviceStub) UpdateGauge(name string, value float64) error {
	s.gaugeCalls++
	s.name, s.value = name, value
	return s.err
}

func (s *serviceStub) UpdateCounter(name string, delta int64) error {
	s.counterCalls++
	s.name, s.delta = name, delta
	return s.err
}

func (s *serviceStub) GetGauge(name string) (float64, error) {
	if s.err != nil {
		return 0, s.err
	}
	value, ok := s.gauges[name]
	if !ok {
		return 0, service.ErrMetricNotFound
	}
	return value, nil
}

func (s *serviceStub) GetCounter(name string) (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	value, ok := s.counters[name]
	if !ok {
		return 0, service.ErrMetricNotFound
	}
	return value, nil
}

func (s *serviceStub) GetGauges() map[string]float64 { return s.gauges }
func (s *serviceStub) GetCounters() map[string]int64 { return s.counters }

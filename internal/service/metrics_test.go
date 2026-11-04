package service

import (
	"errors"
	"testing"
)

type storageStub struct {
	gauges                   map[string]float64
	counters                 map[string]int64
	readCalls                int
	gaugeCalls, counterCalls int
	name                     string
	value                    float64
	delta                    int64
}

func (s *storageStub) GetGauge(name string) (float64, bool) {
	s.readCalls++
	value, ok := s.gauges[name]
	return value, ok
}

func (s *storageStub) GetCounter(name string) (int64, bool) {
	s.readCalls++
	value, ok := s.counters[name]
	return value, ok
}

func (s *storageStub) GetGauges() map[string]float64 { return s.gauges }
func (s *storageStub) GetCounters() map[string]int64 { return s.counters }

func TestGetMetric(t *testing.T) {
	for _, name := range []string{"existing", "zero", "missing", ""} {
		t.Run(name, func(t *testing.T) {
			storage := &storageStub{
				gauges:   map[string]float64{"existing": -12.5, "zero": 0},
				counters: map[string]int64{"existing": 9007199254740993, "zero": 0},
			}
			svc := NewMetricsService(storage)
			gauge, gaugeErr := svc.GetGauge(name)
			counter, counterErr := svc.GetCounter(name)
			var wantErr error
			switch name {
			case "":
				wantErr = ErrEmptyName
			case "missing":
				wantErr = ErrMetricNotFound
			}
			if !errors.Is(gaugeErr, wantErr) || !errors.Is(counterErr, wantErr) {
				t.Fatalf("errors = %v, %v; want %v", gaugeErr, counterErr, wantErr)
			}
			if wantErr == nil && (gauge != storage.gauges[name] || counter != storage.counters[name]) {
				t.Errorf("values = %v, %d", gauge, counter)
			}
			if name == "" && storage.readCalls != 0 {
				t.Errorf("empty name reached storage: %d calls", storage.readCalls)
			}
		})
	}
}

func TestListMetrics(t *testing.T) {
	storage := &storageStub{
		gauges:   map[string]float64{"Alloc": 12.5},
		counters: map[string]int64{"PollCount": 3},
	}
	svc := NewMetricsService(storage)
	gauges, counters := svc.GetGauges(), svc.GetCounters()
	if len(gauges) != 1 || gauges["Alloc"] != 12.5 || len(counters) != 1 || counters["PollCount"] != 3 {
		t.Fatalf("metrics = %v, %v", gauges, counters)
	}
}

func (s *storageStub) SetGauge(name string, value float64) {
	s.gaugeCalls++
	s.name, s.value = name, value
}

func (s *storageStub) AddCounter(name string, delta int64) {
	s.counterCalls++
	s.name, s.delta = name, delta
}

func TestUpdateGauge(t *testing.T) {
	for _, name := range []string{"Alloc", ""} {
		t.Run(name, func(t *testing.T) {
			storage := &storageStub{}
			err := NewMetricsService(storage).UpdateGauge(name, -12.5)
			if name == "" {
				if !errors.Is(err, ErrEmptyName) || storage.gaugeCalls != 0 || storage.counterCalls != 0 {
					t.Fatalf("empty name: err=%v, storage=%+v", err, storage)
				}
				return
			}
			if err != nil || storage.gaugeCalls != 1 || storage.counterCalls != 0 || storage.name != name || storage.value != -12.5 {
				t.Fatalf("update: err=%v, storage=%+v", err, storage)
			}
		})
	}
}

func TestUpdateCounter(t *testing.T) {
	for _, name := range []string{"PollCount", ""} {
		t.Run(name, func(t *testing.T) {
			storage := &storageStub{}
			const delta int64 = 9007199254740993
			err := NewMetricsService(storage).UpdateCounter(name, delta)
			if name == "" {
				if !errors.Is(err, ErrEmptyName) || storage.gaugeCalls != 0 || storage.counterCalls != 0 {
					t.Fatalf("empty name: err=%v, storage=%+v", err, storage)
				}
				return
			}
			if err != nil || storage.counterCalls != 1 || storage.gaugeCalls != 0 || storage.name != name || storage.delta != delta {
				t.Fatalf("update: err=%v, storage=%+v", err, storage)
			}
		})
	}
}

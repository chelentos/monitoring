package repository

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func InitMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (s *MemStorage) SetGauge(name string, value float64) {
	s.gauges[name] = value
}

func (s *MemStorage) AddCounter(name string, value int64) {
	s.counters[name] += value
}

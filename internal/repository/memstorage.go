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

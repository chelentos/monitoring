package repository

import (
	"reflect"
	"testing"
)

func TestInitMemStorage(t *testing.T) {
	t.Run("initialized empty maps", func(t *testing.T) {
		got := InitMemStorage()
		want := &MemStorage{
			gauges: map[string]float64{}, counters: map[string]int64{},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("InitMemStorage() = %#v, want %#v", got, want)
		}
	})
}

func TestMemStorage_SetGauge(t *testing.T) {
	tests := []struct {
		name    string
		initial map[string]float64
		value   float64
		want    map[string]float64
	}{
		{
			name: "new metric", initial: map[string]float64{}, value: 12.5,
			want: map[string]float64{"Alloc": 12.5},
		},
		{
			name:    "replace existing value",
			initial: map[string]float64{"Alloc": 12.5, "Sys": 100}, value: 0,
			want: map[string]float64{"Alloc": 0, "Sys": 100},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := InitMemStorage()
			s.gauges = tt.initial
			s.SetGauge("Alloc", tt.value)
			if !reflect.DeepEqual(s.gauges, tt.want) {
				t.Errorf("gauges = %v, want %v", s.gauges, tt.want)
			}
		})
	}
}

func TestMemStorage_AddCounter(t *testing.T) {
	tests := []struct {
		name    string
		initial map[string]int64
		value   int64
		want    map[string]int64
	}{
		{
			name: "new metric", initial: map[string]int64{}, value: 3,
			want: map[string]int64{"PollCount": 3},
		},
		{
			name:    "add to existing value",
			initial: map[string]int64{"PollCount": 3, "Other": 10}, value: 2,
			want: map[string]int64{"PollCount": 5, "Other": 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := InitMemStorage()
			s.counters = tt.initial
			s.AddCounter("PollCount", tt.value)
			if !reflect.DeepEqual(s.counters, tt.want) {
				t.Errorf("counters = %v, want %v", s.counters, tt.want)
			}
		})
	}
}

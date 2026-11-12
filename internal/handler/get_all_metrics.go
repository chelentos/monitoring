package handler

import (
	"html/template"
	"log"
	"net/http"
)

var metricsTemplate = template.Must(template.New("metrics").Parse(`
<!DOCTYPE html>
<html lang="ru">
<head><meta charset="utf-8"><title>Метрики</title></head>
<body>
	<h2>Метрики</h2>
	<ul>
		{{range $name, $value := .Gauges}}
			<li>{{$name}}: {{$value}}</li>
		{{else}}
			<li>Метрик пока нет</li>
		{{end}}
	</ul>

	<h2>Счетчики</h2>
	<ul>
	{{range $name, $value := .Counters}}
		<li>{{$name}}: {{$value}}</li>
	{{else}}
		<li>Счётчиков пока нет</li>
	{{end}}
	</ul>
</body>
</html>
`))

type MetricsListService interface {
	GetGauges() map[string]float64
	GetCounters() map[string]int64
}

func GetAllMetricHandler(reader MetricsListService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Gauges   map[string]float64
			Counters map[string]int64
		}{
			Gauges:   reader.GetGauges(),
			Counters: reader.GetCounters(),
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := metricsTemplate.Execute(w, data); err != nil {
			log.Printf("render metrics: %v", err)
		}
	}
}

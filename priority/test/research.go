package test

import (
	"sort"
	"time"

	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
)

func FilterByKind(gauges []Gauge, kind GaugeKind) []Gauge {
	out := make([]Gauge, 0, len(gauges))

	for _, gauge := range gauges {
		if gauge.Kind != kind {
			continue
		}

		out = append(out, gauge)
	}

	return out
}

func FilterByPriority(gauges []Gauge, priority uint) []Gauge {
	out := make([]Gauge, 0, len(gauges))

	for _, gauge := range gauges {
		if gauge.Priority != priority {
			continue
		}

		out = append(out, gauge)
	}

	return out
}

func SortByData(gauges []Gauge) {
	less := func(i int, j int) bool {
		return gauges[i].Data < gauges[j].Data
	}

	sort.SliceStable(gauges, less)
}

func SortByDuration(gauges []Gauge) {
	less := func(i int, j int) bool {
		return gauges[i].Duration < gauges[j].Duration
	}

	sort.SliceStable(gauges, less)
}

func SortDurations(durations []time.Duration) {
	less := func(i int, j int) bool {
		return durations[i] < durations[j]
	}

	sort.SliceStable(durations, less)
}

func CalcDataQuantityOverTime(
	gauges []Gauge,
	resolution time.Duration,
	unit time.Duration,
) (map[uint][]chartsopts.LineData, []time.Duration) {
	if len(gauges) == 0 {
		return nil, nil
	}

	SortByDuration(gauges)

	maxDuration := gauges[len(gauges)-1].Duration
	seriesesSize := maxDuration / resolution

	serieses := make(map[uint][]chartsopts.LineData)
	intervals := make([]time.Duration, 0, seriesesSize)

	for _, gauge := range gauges {
		if _, exists := serieses[gauge.Priority]; exists {
			continue
		}

		serieses[gauge.Priority] = make([]chartsopts.LineData, 0, seriesesSize)
	}

	intervalEdge := 0

	for duration := time.Duration(0); duration <= maxDuration; duration += resolution {
		intervals = append(intervals, duration/unit)

		intervalQuantities := make(map[uint]uint)

		for id, gauge := range gauges[intervalEdge:] {
			if gauge.Duration > duration {
				intervalEdge += id
				break
			}

			intervalQuantities[gauge.Priority]++
		}

		for priority, quantity := range intervalQuantities {
			item := chartsopts.LineData{
				Name:  duration.String(),
				Value: quantity,
			}

			serieses[priority] = append(serieses[priority], item)
		}

		for priority := range serieses {
			if _, exists := intervalQuantities[priority]; exists {
				continue
			}

			item := chartsopts.LineData{
				Name:  duration.String(),
				Value: 0,
			}

			serieses[priority] = append(serieses[priority], item)
		}
	}

	return serieses, intervals
}

func CalcWriteToFeedbackLatency(
	gauges []Gauge,
	interval time.Duration,
) (map[uint][]chartsopts.BarData, []int) {
	if len(gauges) == 0 {
		return nil, nil
	}

	SortByData(gauges)

	latencies := make(map[uint][]time.Duration)
	pairs := make(map[uint]Gauge)

	for _, gauge := range gauges {
		switch gauge.Kind {
		case GaugeKindCompleted:
			if _, exists := pairs[gauge.Priority]; !exists {
				pairs[gauge.Priority] = gauge
				continue
			}

			latency := gauge.Duration - pairs[gauge.Priority].Duration

			latencies[gauge.Priority] = append(latencies[gauge.Priority], latency)

			delete(pairs, gauge.Priority)
		case GaugeKindProcessed:
			if _, exists := pairs[gauge.Priority]; !exists {
				pairs[gauge.Priority] = gauge
				continue
			}

			latency := pairs[gauge.Priority].Duration - gauge.Duration

			latencies[gauge.Priority] = append(latencies[gauge.Priority], latency)

			delete(pairs, gauge.Priority)
		}
	}

	return ProcessLatencies(latencies, interval)
}

func ProcessLatencies(
	latencies map[uint][]time.Duration,
	interval time.Duration,
) (map[uint][]chartsopts.BarData, []int) {
	for priority := range latencies {
		SortDurations(latencies[priority])
	}

	maxLatency := time.Duration(0)
	minlatency := time.Duration(0)

	for priority := range latencies {
		if len(latencies[priority]) == 0 {
			continue
		}

		latency := latencies[priority][len(latencies[priority])-1]

		if latency > maxLatency {
			maxLatency = latency
		}

		if latency < minlatency {
			minlatency = latency
		}
	}

	seriesesSize := (maxLatency - minlatency) / interval

	serieses := make(map[uint][]chartsopts.BarData)
	sequences := make([]int, 0, seriesesSize)

	for priority := range latencies {
		if _, exists := serieses[priority]; exists {
			continue
		}

		serieses[priority] = make([]chartsopts.BarData, 0, seriesesSize)
	}

	sequence := 0
	intervalEdge := make(map[uint]int)

	for intervalLatency := minlatency; intervalLatency <= maxLatency; intervalLatency += interval {
		sequences = append(sequences, sequence)
		sequence++

		intervalQuantities := make(map[uint]uint)

		for priority := range latencies {
			for id, latency := range latencies[priority][intervalEdge[priority]:] {
				if latency > intervalLatency {
					intervalEdge[priority] += id
					break
				}

				intervalQuantities[priority]++
			}
		}

		for priority, quantity := range intervalQuantities {
			item := chartsopts.BarData{
				Name:    intervalLatency.String(),
				Tooltip: &chartsopts.Tooltip{Show: true},
				Value:   quantity,
			}

			serieses[priority] = append(serieses[priority], item)
		}

		for priority := range serieses {
			if _, exists := intervalQuantities[priority]; exists {
				continue
			}

			item := chartsopts.BarData{
				Name:    intervalLatency.String(),
				Tooltip: &chartsopts.Tooltip{Show: true},
				Value:   0,
			}

			serieses[priority] = append(serieses[priority], item)
		}
	}

	return serieses, sequences
}

func CalcInProcessingOverTime(
	gauges []Gauge,
	resolution time.Duration,
	unit time.Duration,
) (map[uint][]chartsopts.LineData, []time.Duration) {
	if len(gauges) == 0 {
		return nil, nil
	}

	SortByDuration(gauges)

	maxDuration := gauges[len(gauges)-1].Duration
	seriesesSize := maxDuration / resolution

	serieses := make(map[uint][]chartsopts.LineData)
	intervals := make([]time.Duration, 0, seriesesSize)

	for _, gauge := range gauges {
		if _, exists := serieses[gauge.Priority]; exists {
			continue
		}

		serieses[gauge.Priority] = make([]chartsopts.LineData, 0, seriesesSize)
	}

	intervalEdge := 0

	for duration := time.Duration(0); duration <= maxDuration; duration += resolution {
		intervals = append(intervals, duration/unit)

		receivedQuantities := make(map[uint]map[uint]uint)
		completedQuantities := make(map[uint]map[uint]uint)

		for priority := range serieses {
			receivedQuantities[priority] = make(map[uint]uint)
			completedQuantities[priority] = make(map[uint]uint)
		}

		for id, gauge := range gauges[intervalEdge:] {
			if gauge.Duration > duration {
				intervalEdge += id
				break
			}

			switch gauge.Kind {
			case GaugeKindReceived:
				receivedQuantities[gauge.Priority][gauge.Data]++
			case GaugeKindCompleted:
				completedQuantities[gauge.Priority][gauge.Data]++
			}
		}

		for priority, subset := range receivedQuantities {
			quantity := uint(0)

			for data, amount := range subset {
				if _, exists := completedQuantities[priority][data]; exists {
					continue
				}

				quantity += amount
			}

			item := chartsopts.LineData{
				Name:  duration.String(),
				Value: quantity,
			}

			serieses[priority] = append(serieses[priority], item)
		}

		for priority := range serieses {
			if _, exists := receivedQuantities[priority]; exists {
				continue
			}

			item := chartsopts.LineData{
				Name:  duration.String(),
				Value: 0,
			}

			serieses[priority] = append(serieses[priority], item)
		}
	}

	return serieses, intervals
}

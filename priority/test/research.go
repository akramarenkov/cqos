package test

import (
	"sort"
	"time"

	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
)

type QuantityOverTime struct {
	Quantity     uint
	RelativeTime time.Duration
}

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

func SortByRelativeTime(gauges []Gauge) {
	less := func(i int, j int) bool {
		return gauges[i].RelativeTime < gauges[j].RelativeTime
	}

	sort.SliceStable(gauges, less)
}

func SortDurations(durations []time.Duration) {
	less := func(i int, j int) bool {
		return durations[i] < durations[j]
	}

	sort.SliceStable(durations, less)
}

func CalcDataQuantity(
	gauges []Gauge,
	resolution time.Duration,
) map[uint][]QuantityOverTime {
	if len(gauges) == 0 {
		return nil
	}

	SortByRelativeTime(gauges)

	minRelativeTime := time.Duration(0)
	maxRelativeTime := gauges[len(gauges)-1].RelativeTime

	quantitiesCapacity := (maxRelativeTime - minRelativeTime) / resolution

	quantities := make(map[uint][]QuantityOverTime)

	for _, gauge := range gauges {
		if _, exists := quantities[gauge.Priority]; exists {
			continue
		}

		quantities[gauge.Priority] = make([]QuantityOverTime, 0, quantitiesCapacity)
	}

	gaugesEdge := 0

	for relativeTime := minRelativeTime; relativeTime <= maxRelativeTime; relativeTime += resolution {
		intervalQuantities := make(map[uint]uint)

		for id, gauge := range gauges[gaugesEdge:] {
			if gauge.RelativeTime > relativeTime {
				gaugesEdge += id
				break
			}

			intervalQuantities[gauge.Priority]++
		}

		for priority, quantity := range intervalQuantities {
			item := QuantityOverTime{
				RelativeTime: relativeTime,
				Quantity:     quantity,
			}

			quantities[priority] = append(quantities[priority], item)
		}

		for priority := range quantities {
			if _, exists := intervalQuantities[priority]; exists {
				continue
			}

			item := QuantityOverTime{
				RelativeTime: relativeTime,
				Quantity:     0,
			}

			quantities[priority] = append(quantities[priority], item)
		}
	}

	return quantities
}

func CalcInProcessing(
	gauges []Gauge,
	resolution time.Duration,
) map[uint][]QuantityOverTime {
	if len(gauges) == 0 {
		return nil
	}

	SortByRelativeTime(gauges)

	minRelativeTime := time.Duration(0)
	maxRelativeTime := gauges[len(gauges)-1].RelativeTime

	quantitiesCapacity := (maxRelativeTime - minRelativeTime) / resolution

	quantities := make(map[uint][]QuantityOverTime)

	for _, gauge := range gauges {
		if _, exists := quantities[gauge.Priority]; exists {
			continue
		}

		quantities[gauge.Priority] = make([]QuantityOverTime, 0, quantitiesCapacity)
	}

	gaugesEdge := 0

	for relativeTime := minRelativeTime; relativeTime <= maxRelativeTime; relativeTime += resolution {
		receivedQuantities := make(map[uint]map[uint]uint)
		completedQuantities := make(map[uint]map[uint]uint)

		for priority := range quantities {
			receivedQuantities[priority] = make(map[uint]uint)
			completedQuantities[priority] = make(map[uint]uint)
		}

		for id, gauge := range gauges[gaugesEdge:] {
			if gauge.RelativeTime > relativeTime {
				gaugesEdge += id
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

			item := QuantityOverTime{
				RelativeTime: relativeTime,
				Quantity:     quantity,
			}

			quantities[priority] = append(quantities[priority], item)
		}

		for priority := range quantities {
			if _, exists := receivedQuantities[priority]; exists {
				continue
			}

			item := QuantityOverTime{
				RelativeTime: relativeTime,
				Quantity:     0,
			}

			quantities[priority] = append(quantities[priority], item)
		}
	}

	return quantities
}

func CalcWriteToFeedbackLatency(
	gauges []Gauge,
	interval time.Duration,
) map[uint][]QuantityOverTime {
	if len(gauges) == 0 {
		return nil
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

			latency := gauge.RelativeTime - pairs[gauge.Priority].RelativeTime

			latencies[gauge.Priority] = append(latencies[gauge.Priority], latency)

			delete(pairs, gauge.Priority)
		case GaugeKindProcessed:
			if _, exists := pairs[gauge.Priority]; !exists {
				pairs[gauge.Priority] = gauge
				continue
			}

			latency := pairs[gauge.Priority].RelativeTime - gauge.RelativeTime

			latencies[gauge.Priority] = append(latencies[gauge.Priority], latency)

			delete(pairs, gauge.Priority)
		}
	}

	return ProcessLatencies(latencies, interval)
}

func ProcessLatencies(
	latencies map[uint][]time.Duration,
	interval time.Duration,
) map[uint][]QuantityOverTime {
	for priority := range latencies {
		SortDurations(latencies[priority])
	}

	minLatency := time.Duration(0)
	maxLatency := time.Duration(0)

	for priority := range latencies {
		if len(latencies[priority]) == 0 {
			continue
		}

		latency := latencies[priority][len(latencies[priority])-1]

		if latency > maxLatency {
			maxLatency = latency
		}
	}

	quantitiesCapacity := (maxLatency - minLatency) / interval

	quantities := make(map[uint][]QuantityOverTime)

	for priority := range latencies {
		if _, exists := quantities[priority]; exists {
			continue
		}

		quantities[priority] = make([]QuantityOverTime, 0, quantitiesCapacity)
	}

	latenciesEdge := make(map[uint]int)

	for intervalLatency := minLatency; intervalLatency <= maxLatency; intervalLatency += interval {
		intervalQuantities := make(map[uint]uint)

		for priority := range latencies {
			for id, latency := range latencies[priority][latenciesEdge[priority]:] {
				if latency > intervalLatency {
					latenciesEdge[priority] += id
					break
				}

				intervalQuantities[priority]++
			}
		}

		for priority, quantity := range intervalQuantities {
			item := QuantityOverTime{
				RelativeTime: intervalLatency,
				Quantity:     quantity,
			}

			quantities[priority] = append(quantities[priority], item)
		}

		for priority := range quantities {
			if _, exists := intervalQuantities[priority]; exists {
				continue
			}

			item := QuantityOverTime{
				RelativeTime: intervalLatency,
				Quantity:     0,
			}

			quantities[priority] = append(quantities[priority], item)
		}
	}

	return quantities
}

func ConvertQuantityOverTimeToLineEcharts(
	quantities map[uint][]QuantityOverTime,
	relativeTimeUnit time.Duration,
) (map[uint][]chartsopts.LineData, []uint) {
	serieses := make(map[uint][]chartsopts.LineData)
	xaxis := []uint(nil)

	for priority := range quantities {
		if _, exists := serieses[priority]; !exists {
			serieses[priority] = make([]chartsopts.LineData, 0, len(quantities[priority]))

			if xaxis == nil {
				xaxis = make([]uint, 0, len(quantities[priority]))
			}
		}

		for _, quantity := range quantities[priority] {
			item := chartsopts.LineData{
				Name:  quantity.RelativeTime.String(),
				Value: quantity.Quantity,
			}

			serieses[priority] = append(serieses[priority], item)

			if len(xaxis) < len(quantities[priority]) {
				xaxis = append(xaxis, uint(quantity.RelativeTime/relativeTimeUnit))
			}
		}
	}

	return serieses, xaxis
}

func ConvertQuantityOverTimeToBarEcharts(
	quantities map[uint][]QuantityOverTime,
) (map[uint][]chartsopts.BarData, []uint) {
	serieses := make(map[uint][]chartsopts.BarData)
	xaxis := []uint(nil)

	for priority := range quantities {
		if _, exists := serieses[priority]; !exists {
			serieses[priority] = make([]chartsopts.BarData, 0, len(quantities[priority]))

			if xaxis == nil {
				xaxis = make([]uint, 0, len(quantities[priority]))
			}
		}

		for id, quantity := range quantities[priority] {
			item := chartsopts.BarData{
				Name: quantity.RelativeTime.String(),
				Tooltip: &chartsopts.Tooltip{
					Show: true,
				},
				Value: quantity.Quantity,
			}

			serieses[priority] = append(serieses[priority], item)

			if len(xaxis) < len(quantities[priority]) {
				xaxis = append(xaxis, uint(id))
			}
		}
	}

	return serieses, xaxis
}

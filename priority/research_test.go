package priority

import (
	"sort"
	"time"

	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
)

type quantityOverTime struct {
	Quantity     uint
	RelativeTime time.Duration
}

func filterByKind(gauges []gauge, kind gaugeKind) []gauge {
	out := make([]gauge, 0, len(gauges))

	for _, gauge := range gauges {
		if gauge.Kind != kind {
			continue
		}

		out = append(out, gauge)
	}

	return out
}

func sortByData(gauges []gauge) {
	less := func(i int, j int) bool {
		return gauges[i].Data < gauges[j].Data
	}

	sort.SliceStable(gauges, less)
}

func sortByRelativeTime(gauges []gauge) {
	less := func(i int, j int) bool {
		return gauges[i].RelativeTime < gauges[j].RelativeTime
	}

	sort.SliceStable(gauges, less)
}

func sortDurations(durations []time.Duration) {
	less := func(i int, j int) bool {
		return durations[i] < durations[j]
	}

	sort.SliceStable(durations, less)
}

func calcDataQuantity(
	gauges []gauge,
	resolution time.Duration,
) map[uint][]quantityOverTime {
	if len(gauges) == 0 {
		return nil
	}

	sortByRelativeTime(gauges)

	minRelativeTime := time.Duration(0)
	maxRelativeTime := gauges[len(gauges)-1].RelativeTime

	quantitiesCapacity := (maxRelativeTime - minRelativeTime) / resolution

	quantities := make(map[uint][]quantityOverTime)

	for _, gauge := range gauges {
		if _, exists := quantities[gauge.Priority]; exists {
			continue
		}

		quantities[gauge.Priority] = make([]quantityOverTime, 0, quantitiesCapacity)
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
			item := quantityOverTime{
				RelativeTime: relativeTime,
				Quantity:     quantity,
			}

			quantities[priority] = append(quantities[priority], item)
		}

		for priority := range quantities {
			if _, exists := intervalQuantities[priority]; exists {
				continue
			}

			item := quantityOverTime{
				RelativeTime: relativeTime,
				Quantity:     0,
			}

			quantities[priority] = append(quantities[priority], item)
		}
	}

	return quantities
}

func calcInProcessing(
	gauges []gauge,
	resolution time.Duration,
) map[uint][]quantityOverTime {
	if len(gauges) == 0 {
		return nil
	}

	sortByRelativeTime(gauges)

	minRelativeTime := time.Duration(0)
	maxRelativeTime := gauges[len(gauges)-1].RelativeTime

	quantitiesCapacity := (maxRelativeTime - minRelativeTime) / resolution

	quantities := make(map[uint][]quantityOverTime)

	for _, gauge := range gauges {
		if _, exists := quantities[gauge.Priority]; exists {
			continue
		}

		quantities[gauge.Priority] = make([]quantityOverTime, 0, quantitiesCapacity)
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
			case gaugeKindReceived:
				receivedQuantities[gauge.Priority][gauge.Data]++
			case gaugeKindCompleted:
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

			item := quantityOverTime{
				RelativeTime: relativeTime,
				Quantity:     quantity,
			}

			quantities[priority] = append(quantities[priority], item)
		}

		for priority := range quantities {
			if _, exists := receivedQuantities[priority]; exists {
				continue
			}

			item := quantityOverTime{
				RelativeTime: relativeTime,
				Quantity:     0,
			}

			quantities[priority] = append(quantities[priority], item)
		}
	}

	return quantities
}

func calcWriteToFeedbackLatency(
	gauges []gauge,
	interval time.Duration,
) map[uint][]quantityOverTime {
	if len(gauges) == 0 {
		return nil
	}

	sortByData(gauges)

	latencies := make(map[uint][]time.Duration)
	pairs := make(map[uint]gauge)

	for _, gauge := range gauges {
		switch gauge.Kind {
		case gaugeKindCompleted:
			if _, exists := pairs[gauge.Priority]; !exists {
				pairs[gauge.Priority] = gauge
				continue
			}

			latency := gauge.RelativeTime - pairs[gauge.Priority].RelativeTime

			latencies[gauge.Priority] = append(latencies[gauge.Priority], latency)

			delete(pairs, gauge.Priority)
		case gaugeKindProcessed:
			if _, exists := pairs[gauge.Priority]; !exists {
				pairs[gauge.Priority] = gauge
				continue
			}

			latency := pairs[gauge.Priority].RelativeTime - gauge.RelativeTime

			latencies[gauge.Priority] = append(latencies[gauge.Priority], latency)

			delete(pairs, gauge.Priority)
		}
	}

	return processLatencies(latencies, interval)
}

func processLatencies(
	latencies map[uint][]time.Duration,
	interval time.Duration,
) map[uint][]quantityOverTime {
	for priority := range latencies {
		sortDurations(latencies[priority])
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

	quantities := make(map[uint][]quantityOverTime)

	for priority := range latencies {
		if _, exists := quantities[priority]; exists {
			continue
		}

		quantities[priority] = make([]quantityOverTime, 0, quantitiesCapacity)
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
			item := quantityOverTime{
				RelativeTime: intervalLatency,
				Quantity:     quantity,
			}

			quantities[priority] = append(quantities[priority], item)
		}

		for priority := range quantities {
			if _, exists := intervalQuantities[priority]; exists {
				continue
			}

			item := quantityOverTime{
				RelativeTime: intervalLatency,
				Quantity:     0,
			}

			quantities[priority] = append(quantities[priority], item)
		}
	}

	return quantities
}

func convertQuantityOverTimeToLineEcharts(
	quantities map[uint][]quantityOverTime,
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

func convertQuantityOverTimeToBarEcharts(
	quantities map[uint][]quantityOverTime,
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

// Internal package with research functions that are used for testing
package research

import (
	"sort"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/durations"
	"github.com/akramarenkov/cqos/v2/internal/qot"
	"github.com/akramarenkov/cqos/v2/priority/internal/measurer"

	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
)

func FilterByKind(measures []measurer.Measure, kind measurer.MeasureKind) []measurer.Measure {
	out := make([]measurer.Measure, 0, len(measures))

	for _, measure := range measures {
		if measure.Kind != kind {
			continue
		}

		out = append(out, measure)
	}

	return out
}

func sortByData(measures []measurer.Measure) {
	less := func(i int, j int) bool {
		return measures[i].Data < measures[j].Data
	}

	sort.SliceStable(measures, less)
}

func sortByRelativeTime(measures []measurer.Measure) {
	less := func(i int, j int) bool {
		return measures[i].RelativeTime < measures[j].RelativeTime
	}

	sort.SliceStable(measures, less)
}

func CalcDataQuantity(
	measures []measurer.Measure,
	resolution time.Duration,
) map[uint][]qot.QuantityOverTime {
	if len(measures) == 0 {
		return nil
	}

	sortByRelativeTime(measures)

	min := -resolution
	max := measures[len(measures)-1].RelativeTime

	capacity := (max - min) / resolution

	quantities := make(map[uint][]qot.QuantityOverTime)

	for _, measure := range measures {
		if _, exists := quantities[measure.Priority]; exists {
			continue
		}

		quantities[measure.Priority] = make([]qot.QuantityOverTime, 0, capacity)
	}

	measuresEdge := 0

	for relativeTime := min + resolution; relativeTime <= max+resolution; relativeTime += resolution {
		intervalQuantities := make(map[uint]uint)

		for id, measure := range measures[measuresEdge:] {
			if measure.RelativeTime >= relativeTime {
				measuresEdge += id
				break
			}

			intervalQuantities[measure.Priority]++

			if id == len(measures[measuresEdge:])-1 {
				measuresEdge += id + 1
			}
		}

		for priority, quantity := range intervalQuantities {
			item := qot.QuantityOverTime{
				Quantity:     quantity,
				RelativeTime: relativeTime - resolution,
			}

			quantities[priority] = append(quantities[priority], item)
		}

		for priority := range quantities {
			if _, exists := intervalQuantities[priority]; exists {
				continue
			}

			item := qot.QuantityOverTime{
				Quantity:     0,
				RelativeTime: relativeTime - resolution,
			}

			quantities[priority] = append(quantities[priority], item)
		}
	}

	return quantities
}

func CalcInProcessing(
	measures []measurer.Measure,
	resolution time.Duration,
) map[uint][]qot.QuantityOverTime {
	if len(measures) == 0 {
		return nil
	}

	sortByRelativeTime(measures)

	min := -resolution
	max := measures[len(measures)-1].RelativeTime

	capacity := (max - min) / resolution

	quantities := make(map[uint][]qot.QuantityOverTime)

	for _, measure := range measures {
		if _, exists := quantities[measure.Priority]; exists {
			continue
		}

		quantities[measure.Priority] = make([]qot.QuantityOverTime, 0, capacity)
	}

	measuresEdge := 0

	for relativeTime := min + resolution; relativeTime <= max+resolution; relativeTime += resolution {
		receivedQuantities := make(map[uint]map[uint]uint)
		completedQuantities := make(map[uint]map[uint]uint)

		for priority := range quantities {
			receivedQuantities[priority] = make(map[uint]uint)
			completedQuantities[priority] = make(map[uint]uint)
		}

		for id, measure := range measures[measuresEdge:] {
			if measure.RelativeTime >= relativeTime {
				measuresEdge += id
				break
			}

			switch measure.Kind {
			case measurer.MeasureKindReceived:
				receivedQuantities[measure.Priority][measure.Data]++
			case measurer.MeasureKindCompleted:
				completedQuantities[measure.Priority][measure.Data]++
			}

			if id == len(measures[measuresEdge:])-1 {
				measuresEdge += id + 1
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

			item := qot.QuantityOverTime{
				Quantity:     quantity,
				RelativeTime: relativeTime - resolution,
			}

			quantities[priority] = append(quantities[priority], item)
		}

		for priority := range quantities {
			if _, exists := receivedQuantities[priority]; exists {
				continue
			}

			item := qot.QuantityOverTime{
				Quantity:     0,
				RelativeTime: relativeTime - resolution,
			}

			quantities[priority] = append(quantities[priority], item)
		}
	}

	return quantities
}

func CalcWriteToFeedbackLatency(
	measures []measurer.Measure,
	interval time.Duration,
) map[uint][]qot.QuantityOverTime {
	if len(measures) == 0 {
		return nil
	}

	sortByData(measures)

	latencies := make(map[uint][]time.Duration)
	pairs := make(map[uint]measurer.Measure)

	for _, measure := range measures {
		switch measure.Kind {
		case measurer.MeasureKindCompleted:
			if _, exists := pairs[measure.Priority]; !exists {
				pairs[measure.Priority] = measure
				continue
			}

			latency := measure.RelativeTime - pairs[measure.Priority].RelativeTime

			latencies[measure.Priority] = append(latencies[measure.Priority], latency)

			delete(pairs, measure.Priority)
		case measurer.MeasureKindProcessed:
			if _, exists := pairs[measure.Priority]; !exists {
				pairs[measure.Priority] = measure
				continue
			}

			latency := pairs[measure.Priority].RelativeTime - measure.RelativeTime

			latencies[measure.Priority] = append(latencies[measure.Priority], latency)

			delete(pairs, measure.Priority)
		}
	}

	return processLatencies(latencies, interval)
}

func processLatencies(
	latencies map[uint][]time.Duration,
	interval time.Duration,
) map[uint][]qot.QuantityOverTime {
	for priority := range latencies {
		durations.Sort(latencies[priority])
	}

	min := -interval
	max := time.Duration(0)

	for priority := range latencies {
		if len(latencies[priority]) == 0 {
			continue
		}

		latency := latencies[priority][len(latencies[priority])-1]

		if latency > max {
			max = latency
		}
	}

	capacity := (max - min) / interval

	quantities := make(map[uint][]qot.QuantityOverTime)

	for priority := range latencies {
		if _, exists := quantities[priority]; exists {
			continue
		}

		quantities[priority] = make([]qot.QuantityOverTime, 0, capacity)
	}

	latenciesEdge := make(map[uint]int)

	for intervalLatency := min + interval; intervalLatency <= max+interval; intervalLatency += interval {
		intervalQuantities := make(map[uint]uint)

		for priority := range latencies {
			for id, latency := range latencies[priority][latenciesEdge[priority]:] {
				if latency >= intervalLatency {
					latenciesEdge[priority] += id
					break
				}

				intervalQuantities[priority]++

				if id == len(latencies[priority][latenciesEdge[priority]:])-1 {
					latenciesEdge[priority] += id + 1
				}
			}
		}

		for priority, quantity := range intervalQuantities {
			item := qot.QuantityOverTime{
				Quantity:     quantity,
				RelativeTime: intervalLatency - interval,
			}

			quantities[priority] = append(quantities[priority], item)
		}

		for priority := range quantities {
			if _, exists := intervalQuantities[priority]; exists {
				continue
			}

			item := qot.QuantityOverTime{
				Quantity:     0,
				RelativeTime: intervalLatency - interval,
			}

			quantities[priority] = append(quantities[priority], item)
		}
	}

	return quantities
}

func ConvertToLineEcharts(
	quantities map[uint][]qot.QuantityOverTime,
	relativeTimeUnit time.Duration,
) (map[uint][]chartsopts.LineData, []int) {
	serieses := make(map[uint][]chartsopts.LineData)
	xaxis := []int(nil)

	for priority := range quantities {
		if _, exists := serieses[priority]; !exists {
			serieses[priority] = make([]chartsopts.LineData, 0, len(quantities[priority]))

			if xaxis == nil {
				xaxis = make([]int, 0, len(quantities[priority]))
			}
		}

		for _, quantity := range quantities[priority] {
			item := chartsopts.LineData{
				Name:  quantity.RelativeTime.String(),
				Value: quantity.Quantity,
			}

			serieses[priority] = append(serieses[priority], item)

			if len(xaxis) < len(quantities[priority]) {
				xaxis = append(xaxis, int(quantity.RelativeTime/relativeTimeUnit))
			}
		}
	}

	return serieses, xaxis
}

func ConvertToBarEcharts(
	quantities map[uint][]qot.QuantityOverTime,
) (map[uint][]chartsopts.BarData, []int) {
	serieses := make(map[uint][]chartsopts.BarData)
	xaxis := []int(nil)

	for priority := range quantities {
		if _, exists := serieses[priority]; !exists {
			serieses[priority] = make([]chartsopts.BarData, 0, len(quantities[priority]))

			if xaxis == nil {
				xaxis = make([]int, 0, len(quantities[priority]))
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
				xaxis = append(xaxis, id)
			}
		}
	}

	return serieses, xaxis
}

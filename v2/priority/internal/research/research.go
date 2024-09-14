// Internal package with research functions that are used for testing.
package research

import (
	"slices"
	"sort"
	"time"

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
) map[uint][]qot.QOT {
	if len(measures) == 0 {
		return nil
	}

	sortByRelativeTime(measures)

	// To see the initial zero values ​​on the graph
	min := -resolution
	// One resolution added to max span value to get max relative time into
	// the last span
	// And one more resolution added to max span value to see the final zero values ​on
	// ​the graph
	expansion := 2
	max := measures[len(measures)-1].RelativeTime + time.Duration(expansion)*resolution

	capacity := (max - min) / resolution

	quantities := make(map[uint][]qot.QOT)

	for _, measure := range measures {
		if _, exists := quantities[measure.Priority]; exists {
			continue
		}

		quantities[measure.Priority] = make([]qot.QOT, 0, capacity)
	}

	measuresEdge := 0

	for span := min + resolution; span <= max; span += resolution {
		intervalQuantities := make(map[uint]uint)

		for id, measure := range measures[measuresEdge:] {
			if measure.RelativeTime >= span {
				measuresEdge += id
				break
			}

			intervalQuantities[measure.Priority]++

			// Prevent use of data from the last slice for spans
			// greater than max relative time + resolution
			if id == len(measures[measuresEdge:])-1 {
				measuresEdge += id + 1
			}
		}

		for priority, quantity := range intervalQuantities {
			item := qot.QOT{
				Quantity:     quantity,
				RelativeTime: span - resolution,
			}

			quantities[priority] = append(quantities[priority], item)
		}

		for priority := range quantities {
			if _, exists := intervalQuantities[priority]; exists {
				continue
			}

			item := qot.QOT{
				Quantity:     0,
				RelativeTime: span - resolution,
			}

			quantities[priority] = append(quantities[priority], item)
		}
	}

	return quantities
}

func CalcInProcessing(
	measures []measurer.Measure,
	resolution time.Duration,
) map[uint][]qot.QOT {
	if len(measures) == 0 {
		return nil
	}

	sortByRelativeTime(measures)

	// To see the initial zero values ​​on the graph
	min := -resolution
	// One resolution added to max span value to get max relative time into
	// the last span
	// And one more resolution added to max span value to see the final zero values ​on
	// ​the graph
	expansion := 2
	max := measures[len(measures)-1].RelativeTime + time.Duration(expansion)*resolution

	capacity := (max - min) / resolution

	quantities := make(map[uint][]qot.QOT)

	for _, measure := range measures {
		if _, exists := quantities[measure.Priority]; exists {
			continue
		}

		quantities[measure.Priority] = make([]qot.QOT, 0, capacity)
	}

	measuresEdge := 0

	receivedQuantities := make(map[uint]map[uint]uint)

	for priority := range quantities {
		receivedQuantities[priority] = make(map[uint]uint)
	}

	for span := min + resolution; span <= max; span += resolution {
		for id, measure := range measures[measuresEdge:] {
			if measure.RelativeTime >= span {
				measuresEdge += id
				break
			}

			switch measure.Kind {
			case measurer.MeasureKindReceived:
				receivedQuantities[measure.Priority][measure.Data]++
			case measurer.MeasureKindCompleted:
				receivedQuantities[measure.Priority][measure.Data]--
			}

			// Prevent use of data from the last slice for spans
			// greater than max relative time + resolution
			if id == len(measures[measuresEdge:])-1 {
				measuresEdge += id + 1
			}
		}

		for priority, subset := range receivedQuantities {
			quantity := uint(0)

			for _, amount := range subset {
				quantity += amount
			}

			item := qot.QOT{
				Quantity:     quantity,
				RelativeTime: span - resolution,
			}

			quantities[priority] = append(quantities[priority], item)
		}
	}

	return quantities
}

func CalcWriteToFeedbackLatency(
	measures []measurer.Measure,
	interval time.Duration,
) map[uint][]qot.QOT {
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
) map[uint][]qot.QOT {
	for priority := range latencies {
		slices.Sort(latencies[priority])
	}

	min := time.Duration(0)
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

	// One resolution added to max span value to get max latency into the last span
	max += interval

	capacity := (max - min) / interval

	quantities := make(map[uint][]qot.QOT)

	for priority := range latencies {
		quantities[priority] = make([]qot.QOT, 0, capacity)
	}

	edges := make(map[uint]int)

	for span := min + interval; span <= max; span += interval {
		spanQuantities := make(map[uint]uint)

		for priority := range latencies {
			for id, latency := range latencies[priority][edges[priority]:] {
				if latency >= span {
					edges[priority] += id
					break
				}

				spanQuantities[priority]++

				// Prevent use of data from the last slice for spans
				// greater than max latency + interval
				if id == len(latencies[priority][edges[priority]:])-1 {
					edges[priority] += id + 1
				}
			}
		}

		for priority, quantity := range spanQuantities {
			item := qot.QOT{
				Quantity:     quantity,
				RelativeTime: span - interval,
			}

			quantities[priority] = append(quantities[priority], item)
		}

		for priority := range quantities {
			if _, exists := spanQuantities[priority]; exists {
				continue
			}

			item := qot.QOT{
				Quantity:     0,
				RelativeTime: span - interval,
			}

			quantities[priority] = append(quantities[priority], item)
		}
	}

	return quantities
}

func ConvertToLineEcharts(
	quantities map[uint][]qot.QOT,
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
	quantities map[uint][]qot.QOT,
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

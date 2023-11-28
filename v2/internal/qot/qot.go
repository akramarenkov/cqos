package qot

import (
	"sort"
	"time"

	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
)

type QuantityOverTime struct {
	Quantity     uint
	RelativeTime time.Duration
}

func SortDurations(durations []time.Duration) {
	less := func(i int, j int) bool {
		return durations[i] < durations[j]
	}

	sort.SliceStable(durations, less)
}

func ConvertToLineEcharts(
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

func ConvertToBarEcharts(
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

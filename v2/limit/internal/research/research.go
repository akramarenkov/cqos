// Internal package with research functions that are used for testing.
package research

import (
	"math"
	"slices"
	"sort"
	"strconv"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/consts"
	"github.com/akramarenkov/cqos/v2/internal/qot"

	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
)

func CalcIntervalQuantities(
	relativeTimes []time.Duration,
	intervalsQuantity int,
	interval time.Duration,
) ([]qot.QOT, time.Duration) {
	if len(relativeTimes) == 0 {
		return nil, 0
	}

	slices.Sort(relativeTimes)

	maxRelativeTime := relativeTimes[len(relativeTimes)-1]

	if interval == 0 {
		if intervalsQuantity == 0 {
			return nil, 0
		}

		// It is necessary that max relative time falls into the last span,
		// so the interval is rounded up
		interval = time.Duration(math.Ceil(float64(maxRelativeTime) / float64(intervalsQuantity)))

		maxRelativeTimesRecalculated := interval * time.Duration(intervalsQuantity)

		// If max relative time is divided entirely, then add one nanosecond
		// so that it falls into the last span
		if maxRelativeTimesRecalculated == maxRelativeTime {
			interval += time.Nanosecond
		}
	} else {
		// Intervals quantity always turns out to be more by one
		// Due to rounding down during integer division and, if max relative
		// time is divided entirely, due to the fact that the span takes
		// into account elements strictly smaller than it
		intervalsQuantity = int(maxRelativeTime/interval) + 1
	}

	quantities := make([]qot.QOT, 0, intervalsQuantity)

	edge := 0

	// Interval is added to max span value to be sure that max relative time falls
	// into the last span
	for span := interval; span <= maxRelativeTime+interval; span += interval {
		spanQuantities := uint(0)

		for id, relativeTime := range relativeTimes[edge:] {
			if relativeTime >= span {
				edge += id
				break
			}

			spanQuantities++

			// Prevent use of data from the last slice for spans
			// greater than max relative time + interval
			if id == len(relativeTimes[edge:])-1 {
				edge += id + 1
			}
		}

		item := qot.QOT{
			Quantity:     spanQuantities,
			RelativeTime: span - interval,
		}

		quantities = append(quantities, item)
	}

	// Padding with zero values ​​in case intervals quantity multiplied by
	// interval is greater than max relative time
	for addition := range intervalsQuantity - len(quantities) {
		item := qot.QOT{
			Quantity:     0,
			RelativeTime: maxRelativeTime + interval*time.Duration(addition+1),
		}

		quantities = append(quantities, item)
	}

	return quantities, interval
}

func ConvertQuantityOverTimeToBarEcharts(
	quantities []qot.QOT,
) ([]chartsopts.BarData, []int) {
	serieses := make([]chartsopts.BarData, 0, len(quantities))
	xaxis := make([]int, 0, len(quantities))

	for id, quantity := range quantities {
		item := chartsopts.BarData{
			Name: quantity.RelativeTime.String(),
			Tooltip: &chartsopts.Tooltip{
				Show: true,
			},
			Value: quantity.Quantity,
		}

		serieses = append(serieses, item)
		xaxis = append(xaxis, id)
	}

	return serieses, xaxis
}

func CalcRelativeDeviations(
	relativeTimes []time.Duration,
	expected time.Duration,
) map[int]int {
	const (
		deviationsMin = -100
		deviationsMax = 100
		// from -100% to 100% with 1% step and plus zero
		deviationsLength = deviationsMax - deviationsMin + 1
	)

	if len(relativeTimes) == 0 {
		return nil
	}

	slices.Sort(relativeTimes)

	deviations := make(map[int]int, deviationsLength)

	for percent := deviationsMin; percent <= deviationsMax; percent++ {
		deviations[percent] = 0
	}

	calc := func(next time.Duration, current time.Duration) {
		diff := next - current

		deviation := ((diff - expected) * consts.HundredPercent) / expected

		if deviation > consts.HundredPercent {
			deviation = consts.HundredPercent
		}

		if deviation < -consts.HundredPercent {
			deviation = -consts.HundredPercent
		}

		deviations[int(deviation)]++
	}

	calc(relativeTimes[0], 0)

	for id := range relativeTimes {
		if id == len(relativeTimes)-1 {
			break
		}

		calc(relativeTimes[id+1], relativeTimes[id])
	}

	return deviations
}

func ConvertRelativeDeviationsToBarEcharts(
	deviations map[int]int,
) ([]chartsopts.BarData, []int) {
	serieses := make([]chartsopts.BarData, 0, len(deviations))
	xaxis := make([]int, 0, len(deviations))

	for percent := range deviations {
		xaxis = append(xaxis, percent)
	}

	sort.Ints(xaxis)

	for _, percent := range xaxis {
		item := chartsopts.BarData{
			Name: strconv.Itoa(percent) + "%",
			Tooltip: &chartsopts.Tooltip{
				Show: true,
			},
			Value: deviations[percent],
		}

		serieses = append(serieses, item)
	}

	return serieses, xaxis
}

func CalcTotalDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	copied := slices.Clone(durations)

	slices.Sort(copied)

	return copied[len(copied)-1]
}

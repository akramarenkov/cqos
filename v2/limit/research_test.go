package limit

import (
	"math"
	"sort"
	"strconv"
	"time"

	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
)

const (
	maxPercentValue = 100
)

type quantityOverTime struct {
	Quantity     uint
	RelativeTime time.Duration
}

func sortDurations(durations []time.Duration) {
	less := func(i int, j int) bool {
		return durations[i] < durations[j]
	}

	sort.SliceStable(durations, less)
}

func IsSortedDurations(durations []time.Duration) bool {
	less := func(i int, j int) bool {
		return durations[i] < durations[j]
	}

	return sort.SliceIsSorted(durations, less)
}

func calcTotalDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	sortDurations(durations)

	return durations[len(durations)-1]
}

func calcIntervalQuantities(
	relativeTimes []time.Duration,
	intervalsQuantity int,
	interval time.Duration,
) []quantityOverTime {
	if len(relativeTimes) == 0 {
		return nil
	}

	sortDurations(relativeTimes)

	maxRelativeTimes := relativeTimes[len(relativeTimes)-1]

	if interval == 0 {
		if intervalsQuantity == 0 {
			return nil
		}

		interval = maxRelativeTimes / time.Duration(intervalsQuantity)
	} else {
		intervalsQuantity = int(maxRelativeTimes / interval)
	}

	if interval == 0 {
		interval = time.Nanosecond
	}

	quantities := make([]quantityOverTime, 0, intervalsQuantity+1)

	edge := 0

	for span := interval; span <= maxRelativeTimes+interval; span += interval {
		spanQuantities := uint(0)

		for id, relativeTime := range relativeTimes[edge:] {
			if relativeTime > span {
				edge += id
				break
			}

			spanQuantities++

			if id == len(relativeTimes[edge:])-1 {
				edge += id + 1
			}
		}

		item := quantityOverTime{
			Quantity:     spanQuantities,
			RelativeTime: span - interval,
		}

		quantities = append(quantities, item)
	}

	return quantities
}

func calcSelfDeviations(
	relativeTimes []time.Duration,
	intervalsQuantity int,
	interval time.Duration,
) ([]quantityOverTime, time.Duration, time.Duration, time.Duration) {
	if len(relativeTimes) == 0 {
		return nil, 0, 0, 0
	}

	sortDurations(relativeTimes)

	deviations := make([]time.Duration, 0, len(relativeTimes))

	min := time.Duration(math.MaxInt)
	max := time.Duration(math.MinInt)
	avg := time.Duration(0)

	calc := func(next time.Duration, current time.Duration) {
		deviation := next - current

		if deviation < min {
			min = deviation
		}

		if deviation > max {
			max = deviation
		}

		avg += deviation

		deviations = append(deviations, deviation)
	}

	calc(relativeTimes[0], 0)

	for id := range relativeTimes {
		if id+1 > len(relativeTimes)-1 {
			break
		}

		calc(relativeTimes[id+1], relativeTimes[id])
	}

	avg /= time.Duration(len(deviations))

	return calcIntervalQuantities(deviations, intervalsQuantity, interval), min, max, avg
}

func convertQuantityOverTimeToBarEcharts(
	quantities []quantityOverTime,
) ([]chartsopts.BarData, []uint) {
	serieses := make([]chartsopts.BarData, 0, len(quantities))
	xaxis := make([]uint, 0, len(quantities))

	for id, quantity := range quantities {
		item := chartsopts.BarData{
			Name: quantity.RelativeTime.String(),
			Tooltip: &chartsopts.Tooltip{
				Show: true,
			},
			Value: quantity.Quantity,
		}

		serieses = append(serieses, item)
		xaxis = append(xaxis, uint(id))
	}

	return serieses, xaxis
}

func calcRelativeDeviations(
	durations []time.Duration,
	expected time.Duration,
) map[int]int {
	if len(durations) == 0 {
		return nil
	}

	sortDurations(durations)

	deviations := make(map[int]int, len(durations)/2)

	calc := func(next time.Duration, current time.Duration) {
		diff := next - current

		deviation := ((diff - expected) * maxPercentValue) / expected

		if deviation > maxPercentValue {
			deviation = maxPercentValue
		}

		deviations[int(deviation)]++
	}

	calc(durations[0], 0)

	for id := range durations {
		if id+1 > len(durations)-1 {
			break
		}

		calc(durations[id+1], durations[id])
	}

	return deviations
}

func convertRelativeDeviationsToBarEcharts(
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

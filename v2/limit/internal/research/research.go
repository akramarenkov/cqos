// Internal package with research functions that are used for testing
package research

import (
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
) ([]qot.QuantityOverTime, time.Duration) {
	if len(relativeTimes) == 0 {
		return nil, 0
	}

	slices.Sort(relativeTimes)

	maxRelativeTimes := relativeTimes[len(relativeTimes)-1]

	if interval == 0 {
		if intervalsQuantity == 0 {
			return nil, 0
		}

		interval = maxRelativeTimes / time.Duration(intervalsQuantity)
	} else {
		intervalsQuantity = int(maxRelativeTimes / interval)
	}

	if interval == 0 {
		interval = time.Nanosecond
	}

	quantities := make([]qot.QuantityOverTime, 0, intervalsQuantity+1)

	edge := 0

	for span := interval; span <= maxRelativeTimes+interval; span += interval {
		spanQuantities := uint(0)

		for id, relativeTime := range relativeTimes[edge:] {
			if relativeTime >= span {
				edge += id
				break
			}

			spanQuantities++

			if id == len(relativeTimes[edge:])-1 {
				edge += id + 1
			}
		}

		item := qot.QuantityOverTime{
			Quantity:     spanQuantities,
			RelativeTime: span - interval,
		}

		quantities = append(quantities, item)
	}

	return quantities, interval
}

func ConvertQuantityOverTimeToBarEcharts(
	quantities []qot.QuantityOverTime,
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

		deviation := ((diff - expected) * consts.OneHundredPercent) / expected

		if deviation > consts.OneHundredPercent {
			deviation = consts.OneHundredPercent
		}

		if deviation < -consts.OneHundredPercent {
			deviation = -consts.OneHundredPercent
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

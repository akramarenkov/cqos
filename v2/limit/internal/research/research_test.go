package research

import (
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/qot"

	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
	"github.com/stretchr/testify/require"
)

func TestCalcIntervalQuantitiesSplitByInterval(t *testing.T) {
	relativeTimes := []time.Duration{
		0,
		time.Millisecond,
		2 * time.Millisecond,
		5 * time.Millisecond,
		9 * time.Millisecond,
		11 * time.Millisecond,
		13 * time.Millisecond,
		17 * time.Millisecond,
	}

	interval := 10 * time.Millisecond

	expectedQuantities := []qot.QuantityOverTime{
		{
			Quantity:     5,
			RelativeTime: 0,
		},
		{
			Quantity:     3,
			RelativeTime: 10 * time.Millisecond,
		},
	}

	expectedBarData := []chartsopts.BarData{
		{
			Name:  "0s",
			Value: uint(5),
			Tooltip: &chartsopts.Tooltip{
				Show: true,
			},
		},
		{
			Name:  "10ms",
			Value: uint(3),
			Tooltip: &chartsopts.Tooltip{
				Show: true,
			},
		},
	}

	quantities, calcInterval := CalcIntervalQuantities(
		relativeTimes,
		0,
		interval,
	)
	require.Equal(t, expectedQuantities, quantities)
	require.Equal(t, interval, calcInterval)

	axisY, axisX := ConvertQuantityOverTimeToBarEcharts(quantities)
	require.Equal(t, expectedBarData, axisY)
	require.Equal(t, []int{0, 1}, axisX)
}

func TestCalcIntervalQuantitiesSplitByIntervalEntirely(t *testing.T) {
	relativeTimes := []time.Duration{
		0,
		time.Millisecond,
		2 * time.Millisecond,
		5 * time.Millisecond,
		9 * time.Millisecond,
		11 * time.Millisecond,
		13 * time.Millisecond,
		20 * time.Millisecond,
	}

	interval := 10 * time.Millisecond

	expected := []qot.QuantityOverTime{
		{
			Quantity:     5,
			RelativeTime: 0,
		},
		{
			Quantity:     2,
			RelativeTime: 10 * time.Millisecond,
		},
		{
			Quantity:     1,
			RelativeTime: 20 * time.Millisecond,
		},
	}

	quantities, calcInterval := CalcIntervalQuantities(
		relativeTimes,
		0,
		interval,
	)
	require.Equal(t, expected, quantities)
	require.Equal(t, interval, calcInterval)
}

func TestCalcIntervalQuantitiesSplitByIntervalsQuantity(t *testing.T) {
	relativeTimes := []time.Duration{
		0,
		time.Millisecond,
		2 * time.Millisecond,
		5 * time.Millisecond,
		9 * time.Millisecond,
		11 * time.Millisecond,
		13 * time.Millisecond,
		17 * time.Millisecond,
	}

	intervalsQuantity := 2

	expectedCalcInterval := 8*time.Millisecond + 500*time.Microsecond + time.Nanosecond

	expected := []qot.QuantityOverTime{
		{
			Quantity:     4,
			RelativeTime: 0,
		},
		{
			Quantity:     4,
			RelativeTime: expectedCalcInterval,
		},
	}

	quantities, calcInterval := CalcIntervalQuantities(
		relativeTimes,
		intervalsQuantity,
		0,
	)
	require.Equal(t, expected, quantities)
	require.Equal(t, expectedCalcInterval, calcInterval)
}

func TestCalcIntervalQuantitiesZeroInput(t *testing.T) {
	quantities, calcInterval := CalcIntervalQuantities(
		nil,
		0,
		time.Second,
	)
	require.Equal(t, []qot.QuantityOverTime(nil), quantities)
	require.Equal(t, time.Duration(0), calcInterval)

	quantities, calcInterval = CalcIntervalQuantities(
		[]time.Duration{},
		0,
		time.Second,
	)
	require.Equal(t, []qot.QuantityOverTime(nil), quantities)
	require.Equal(t, time.Duration(0), calcInterval)
}

func TestCalcIntervalQuantitiesZeroSplit(t *testing.T) {
	quantities, calcInterval := CalcIntervalQuantities(
		[]time.Duration{1, 2},
		0,
		0,
	)
	require.Equal(t, []qot.QuantityOverTime(nil), quantities)
	require.Equal(t, time.Duration(0), calcInterval)
}

func TestCalcIntervalQuantitiesSmallRatio(t *testing.T) {
	relativeTimes := []time.Duration{
		0,
		time.Nanosecond,
		2 * time.Nanosecond,
		5 * time.Nanosecond,
	}

	intervalsQuantity := 10

	expectedCalcInterval := time.Nanosecond

	expected := []qot.QuantityOverTime{
		{
			Quantity:     1,
			RelativeTime: 0,
		},
		{
			Quantity:     1,
			RelativeTime: 1,
		},
		{
			Quantity:     1,
			RelativeTime: 2,
		},
		{
			Quantity:     0,
			RelativeTime: 3,
		},
		{
			Quantity:     0,
			RelativeTime: 4,
		},
		{
			Quantity:     1,
			RelativeTime: 5,
		},
		{
			Quantity:     0,
			RelativeTime: 6,
		},
		{
			Quantity:     0,
			RelativeTime: 7,
		},
		{
			Quantity:     0,
			RelativeTime: 8,
		},
		{
			Quantity:     0,
			RelativeTime: 9,
		},
	}

	quantities, calcInterval := CalcIntervalQuantities(
		relativeTimes,
		intervalsQuantity,
		0,
	)
	require.Equal(t, expected, quantities)
	require.Equal(t, expectedCalcInterval, calcInterval)
}

func TestConvertQuantityOverTimeToBarEcharts(t *testing.T) {
	quantities := []qot.QuantityOverTime{
		{
			Quantity:     5,
			RelativeTime: 0,
		},
		{
			Quantity:     3,
			RelativeTime: 10 * time.Millisecond,
		},
	}

	expectedY := []chartsopts.BarData{
		{
			Name:  "0s",
			Value: uint(5),
			Tooltip: &chartsopts.Tooltip{
				Show: true,
			},
		},
		{
			Name:  "10ms",
			Value: uint(3),
			Tooltip: &chartsopts.Tooltip{
				Show: true,
			},
		},
	}

	expectedX := []int{
		0,
		1,
	}

	axisY, axisX := ConvertQuantityOverTimeToBarEcharts(quantities)
	require.Equal(t, expectedY, axisY)
	require.Equal(t, expectedX, axisX)
}

func TestCalcRelativeDeviations(t *testing.T) {
	relativeTimes := []time.Duration{
		-100 * time.Microsecond,
		900 * time.Microsecond,
		900 * time.Microsecond,
		2000 * time.Microsecond,
		2800 * time.Microsecond,
		3700 * time.Microsecond,
		4700 * time.Microsecond,
		5500 * time.Microsecond,
		6600 * time.Microsecond,
		7600 * time.Microsecond,
		8600 * time.Microsecond,
		9600 * time.Microsecond,
		10700 * time.Microsecond,
		12700 * time.Microsecond,
		14800 * time.Microsecond,
	}

	expected := make(map[int]int, 201)

	for deviation := -100; deviation <= 100; deviation++ {
		expected[deviation] = 0
	}

	expected[-100] = 2
	expected[-20] = 2
	expected[-10] = 1
	expected[0] = 5
	expected[10] = 3
	expected[100] = 2

	deviations := CalcRelativeDeviations(relativeTimes, time.Millisecond)
	require.Equal(t, expected, deviations)
}

func TestCalcRelativeDeviationsZeroInput(t *testing.T) {
	deviations := CalcRelativeDeviations(nil, time.Millisecond)
	require.Equal(t, map[int]int(nil), deviations)

	deviations = CalcRelativeDeviations([]time.Duration{}, time.Millisecond)
	require.Equal(t, map[int]int(nil), deviations)
}

package research

import (
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/qot"
	"github.com/akramarenkov/cqos/v2/priority/internal/measurer"
	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
	"github.com/stretchr/testify/require"
)

func TestCalcWriteToFeedbackLatency(t *testing.T) {
	measures := []measurer.Measure{
		// first priority
		{
			Data:         0,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     1,
			RelativeTime: 11 * time.Microsecond,
		},
		{
			Data:         0,
			Kind:         measurer.MeasureKindProcessed,
			Priority:     1,
			RelativeTime: 10 * time.Microsecond,
		},
		{
			Data:         0,
			Kind:         measurer.MeasureKindReceived,
			Priority:     1,
			RelativeTime: 0,
		},
		{
			Data:         1,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     1,
			RelativeTime: 10 * time.Microsecond,
		},
		{
			Data:         1,
			Kind:         measurer.MeasureKindProcessed,
			Priority:     1,
			RelativeTime: 2 * time.Microsecond,
		},
		{
			Data:         1,
			Kind:         measurer.MeasureKindReceived,
			Priority:     1,
			RelativeTime: 0,
		},
		{
			Data:         2,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     1,
			RelativeTime: 28 * time.Microsecond,
		},
		{
			Data:         2,
			Kind:         measurer.MeasureKindProcessed,
			Priority:     1,
			RelativeTime: 25 * time.Microsecond,
		},
		{
			Data:         2,
			Kind:         measurer.MeasureKindReceived,
			Priority:     1,
			RelativeTime: 20,
		},
		{
			Data:         3,
			Kind:         measurer.MeasureKindProcessed,
			Priority:     1,
			RelativeTime: 35 * time.Microsecond,
		},
		{
			Data:         3,
			Kind:         measurer.MeasureKindReceived,
			Priority:     1,
			RelativeTime: 30,
		},
		{
			Data:         3,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     1,
			RelativeTime: 40 * time.Microsecond,
		},
		// third priority
		{
			Data:         0,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     3,
			RelativeTime: 4 * time.Microsecond,
		},
		{
			Data:         0,
			Kind:         measurer.MeasureKindProcessed,
			Priority:     3,
			RelativeTime: 3 * time.Microsecond,
		},
		{
			Data:         0,
			Kind:         measurer.MeasureKindReceived,
			Priority:     3,
			RelativeTime: 0,
		},
		{
			Data:         1,
			Kind:         measurer.MeasureKindReceived,
			Priority:     3,
			RelativeTime: 0,
		},
		{
			Data:         1,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     3,
			RelativeTime: 5 * time.Microsecond,
		},
		{
			Data:         1,
			Kind:         measurer.MeasureKindProcessed,
			Priority:     3,
			RelativeTime: 3 * time.Microsecond,
		},
		{
			Data:         2,
			Kind:         measurer.MeasureKindProcessed,
			Priority:     3,
			RelativeTime: 3 * time.Microsecond,
		},
		{
			Data:         2,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     3,
			RelativeTime: 7 * time.Microsecond,
		},
		{
			Data:         2,
			Kind:         measurer.MeasureKindReceived,
			Priority:     3,
			RelativeTime: 0,
		},
		{
			Data:         3,
			Kind:         measurer.MeasureKindProcessed,
			Priority:     3,
			RelativeTime: 3 * time.Microsecond,
		},
		{
			Data:         3,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     3,
			RelativeTime: 8 * time.Microsecond,
		},
		{
			Data:         3,
			Kind:         measurer.MeasureKindReceived,
			Priority:     3,
			RelativeTime: 0,
		},
		{
			Data:         4,
			Kind:         measurer.MeasureKindProcessed,
			Priority:     3,
			RelativeTime: 3 * time.Microsecond,
		},
		{
			Data:         4,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     3,
			RelativeTime: 9 * time.Microsecond,
		},
		{
			Data:         4,
			Kind:         measurer.MeasureKindReceived,
			Priority:     3,
			RelativeTime: 0,
		},
		{
			Data:         5,
			Kind:         measurer.MeasureKindProcessed,
			Priority:     3,
			RelativeTime: 3 * time.Microsecond,
		},
		{
			Data:         5,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     3,
			RelativeTime: 19 * time.Microsecond,
		},
		{
			Data:         5,
			Kind:         measurer.MeasureKindReceived,
			Priority:     3,
			RelativeTime: 0,
		},
	}

	interval := 5 * time.Microsecond

	expected := map[uint][]qot.QuantityOverTime{
		1: {
			{
				RelativeTime: 0,
				Quantity:     2,
			},
			{
				RelativeTime: interval,
				Quantity:     2,
			},
			{
				RelativeTime: 2 * interval,
				Quantity:     0,
			},
			{
				RelativeTime: 3 * interval,
				Quantity:     0,
			},
		},
		3: {
			{
				RelativeTime: 0,
				Quantity:     3,
			},
			{
				RelativeTime: interval,
				Quantity:     2,
			},
			{
				RelativeTime: 2 * interval,
				Quantity:     0,
			},
			{
				RelativeTime: 3 * interval,
				Quantity:     1,
			},
		},
	}

	quantities := CalcWriteToFeedbackLatency(measures, interval)
	require.Equal(t, expected, quantities)
}

func TestCalcWriteToFeedbackLatencyInput(t *testing.T) {
	quantities := CalcWriteToFeedbackLatency(nil, 5*time.Microsecond)
	require.Equal(t, map[uint][]qot.QuantityOverTime(nil), quantities)

	quantities = CalcWriteToFeedbackLatency([]measurer.Measure{}, 5*time.Microsecond)
	require.Equal(t, map[uint][]qot.QuantityOverTime(nil), quantities)
}

func TestProcessLatencies(t *testing.T) {
	latencies := map[uint][]time.Duration{
		1: {
			time.Microsecond,
			8 * time.Microsecond,
			3 * time.Microsecond,
			5 * time.Microsecond,
		},
		2: {},
		3: {
			1 * time.Microsecond,
			2 * time.Microsecond,
			4 * time.Microsecond,
			5 * time.Microsecond,
			6 * time.Microsecond,
			16 * time.Microsecond,
		},
	}

	interval := 5 * time.Microsecond

	expected := map[uint][]qot.QuantityOverTime{
		1: {
			{
				RelativeTime: 0,
				Quantity:     2,
			},
			{
				RelativeTime: interval,
				Quantity:     2,
			},
			{
				RelativeTime: 2 * interval,
				Quantity:     0,
			},
			{
				RelativeTime: 3 * interval,
				Quantity:     0,
			},
		},
		2: {
			{
				RelativeTime: 0,
				Quantity:     0,
			},
			{
				RelativeTime: interval,
				Quantity:     0,
			},
			{
				RelativeTime: 2 * interval,
				Quantity:     0,
			},
			{
				RelativeTime: 3 * interval,
				Quantity:     0,
			},
		},
		3: {
			{
				RelativeTime: 0,
				Quantity:     3,
			},
			{
				RelativeTime: interval,
				Quantity:     2,
			},
			{
				RelativeTime: 2 * interval,
				Quantity:     0,
			},
			{
				RelativeTime: 3 * interval,
				Quantity:     1,
			},
		},
	}

	quantities := processLatencies(latencies, interval)
	require.Equal(t, expected, quantities)
}

func TestConvertToBarEcharts(t *testing.T) {
	resolution := 5 * time.Microsecond

	quantities := map[uint][]qot.QuantityOverTime{
		1: {
			{
				RelativeTime: -resolution,
				Quantity:     0,
			},
			{
				RelativeTime: 0,
				Quantity:     1,
			},
			{
				RelativeTime: resolution,
				Quantity:     0,
			},
		},
		2: {
			{
				RelativeTime: -resolution,
				Quantity:     0,
			},
			{
				RelativeTime: 0,
				Quantity:     2,
			},
			{
				RelativeTime: resolution,
				Quantity:     0,
			},
		},
		3: {
			{
				RelativeTime: -resolution,
				Quantity:     0,
			},
			{
				RelativeTime: 0,
				Quantity:     3,
			},
			{
				RelativeTime: resolution,
				Quantity:     0,
			},
		},
	}

	expectedY := map[uint][]chartsopts.BarData{
		1: {
			{
				Name: "-5µs",
				Tooltip: &chartsopts.Tooltip{
					Show: true,
				},
				Value: uint(0),
			},
			{
				Name: "0s",
				Tooltip: &chartsopts.Tooltip{
					Show: true,
				},
				Value: uint(1),
			},
			{
				Name: "5µs",
				Tooltip: &chartsopts.Tooltip{
					Show: true,
				},
				Value: uint(0),
			},
		},
		2: {
			{
				Name: "-5µs",
				Tooltip: &chartsopts.Tooltip{
					Show: true,
				},
				Value: uint(0),
			},
			{
				Name: "0s",
				Tooltip: &chartsopts.Tooltip{
					Show: true,
				},
				Value: uint(2),
			},
			{
				Name: "5µs",
				Tooltip: &chartsopts.Tooltip{
					Show: true,
				},
				Value: uint(0),
			},
		},
		3: {
			{
				Name: "-5µs",
				Tooltip: &chartsopts.Tooltip{
					Show: true,
				},
				Value: uint(0),
			},
			{
				Name: "0s",
				Tooltip: &chartsopts.Tooltip{
					Show: true,
				},
				Value: uint(3),
			},
			{
				Name: "5µs",
				Tooltip: &chartsopts.Tooltip{
					Show: true,
				},
				Value: uint(0),
			},
		},
	}

	expectedX := []int{0, 1, 2}

	axisY, axisX := ConvertToBarEcharts(quantities)
	require.Equal(t, expectedY, axisY)
	require.Equal(t, expectedX, axisX)
}

package research

import (
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/qot"
	"github.com/akramarenkov/cqos/v2/priority/internal/measurer"
	"github.com/stretchr/testify/require"
)

func TestCalcDataQuantity(t *testing.T) {
	measures := []measurer.Measure{
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
			Priority:     2,
			RelativeTime: 25 * time.Microsecond,
		},
		{
			Data:         1,
			Kind:         measurer.MeasureKindProcessed,
			Priority:     2,
			RelativeTime: 20 * time.Microsecond,
		},
		{
			Data:         1,
			Kind:         measurer.MeasureKindReceived,
			Priority:     2,
			RelativeTime: time.Microsecond,
		},
		{
			Data:         2,
			Kind:         measurer.MeasureKindProcessed,
			Priority:     3,
			RelativeTime: 30 * time.Microsecond,
		},
		{
			Data:         2,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     3,
			RelativeTime: 33 * time.Microsecond,
		},
		{
			Data:         2,
			Kind:         measurer.MeasureKindReceived,
			Priority:     3,
			RelativeTime: 0,
		},
	}

	resolution := 5 * time.Microsecond

	expected := map[uint][]qot.QuantityOverTime{
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
			{
				RelativeTime: 2 * resolution,
				Quantity:     2,
			},
			{
				RelativeTime: 3 * resolution,
				Quantity:     0,
			},
			{
				RelativeTime: 4 * resolution,
				Quantity:     0,
			},
			{
				RelativeTime: 5 * resolution,
				Quantity:     0,
			},
			{
				RelativeTime: 6 * resolution,
				Quantity:     0,
			},
			{
				RelativeTime: 7 * resolution,
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
				Quantity:     1,
			},
			{
				RelativeTime: resolution,
				Quantity:     0,
			},
			{
				RelativeTime: 2 * resolution,
				Quantity:     0,
			},
			{
				RelativeTime: 3 * resolution,
				Quantity:     0,
			},
			{
				RelativeTime: 4 * resolution,
				Quantity:     1,
			},
			{
				RelativeTime: 5 * resolution,
				Quantity:     1,
			},
			{
				RelativeTime: 6 * resolution,
				Quantity:     0,
			},
			{
				RelativeTime: 7 * resolution,
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
				Quantity:     1,
			},
			{
				RelativeTime: resolution,
				Quantity:     0,
			},
			{
				RelativeTime: 2 * resolution,
				Quantity:     0,
			},
			{
				RelativeTime: 3 * resolution,
				Quantity:     0,
			},
			{
				RelativeTime: 4 * resolution,
				Quantity:     0,
			},
			{
				RelativeTime: 5 * resolution,
				Quantity:     0,
			},
			{
				RelativeTime: 6 * resolution,
				Quantity:     2,
			},
			{
				RelativeTime: 7 * resolution,
				Quantity:     0,
			},
		},
	}

	quantities := CalcDataQuantity(measures, resolution)
	require.Equal(t, expected, quantities)
}

func TestCalcDataQuantityZeroInput(t *testing.T) {
	quantities := CalcDataQuantity(nil, 5*time.Microsecond)
	require.Equal(t, map[uint][]qot.QuantityOverTime(nil), quantities)

	quantities = CalcDataQuantity([]measurer.Measure{}, 5*time.Microsecond)
	require.Equal(t, map[uint][]qot.QuantityOverTime(nil), quantities)
}

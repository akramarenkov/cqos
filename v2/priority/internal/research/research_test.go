package research

import (
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/priority/internal/measurer"
	"github.com/stretchr/testify/require"
)

func TestFilterByKind(t *testing.T) {
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
			Kind:         measurer.MeasureKindProcessed,
			Priority:     2,
			RelativeTime: 20 * time.Microsecond,
		},
		{
			Data:         1,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     2,
			RelativeTime: 25 * time.Microsecond,
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
			Kind:         measurer.MeasureKindReceived,
			Priority:     3,
			RelativeTime: 0,
		},
		{
			Data:         2,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     3,
			RelativeTime: 33 * time.Microsecond,
		},
	}

	expected := []measurer.Measure{
		{
			Data:         0,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     1,
			RelativeTime: 11 * time.Microsecond,
		},
		{
			Data:         1,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     2,
			RelativeTime: 25 * time.Microsecond,
		},
		{
			Data:         2,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     3,
			RelativeTime: 33 * time.Microsecond,
		},
	}

	filtered := FilterByKind(measures, measurer.MeasureKindCompleted)
	require.Equal(t, expected, filtered)
}

func TestSortByData(t *testing.T) {
	measures := []measurer.Measure{
		{
			Data:         0,
			Kind:         measurer.MeasureKindProcessed,
			Priority:     1,
			RelativeTime: 10 * time.Microsecond,
		},
		{
			Data:         1,
			Kind:         measurer.MeasureKindProcessed,
			Priority:     2,
			RelativeTime: 20 * time.Microsecond,
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
			Data:         0,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     1,
			RelativeTime: 11 * time.Microsecond,
		},
		{
			Data:         2,
			Kind:         measurer.MeasureKindProcessed,
			Priority:     3,
			RelativeTime: 30 * time.Microsecond,
		},
		{
			Data:         2,
			Kind:         measurer.MeasureKindReceived,
			Priority:     3,
			RelativeTime: 0,
		},
		{
			Data:         1,
			Kind:         measurer.MeasureKindReceived,
			Priority:     2,
			RelativeTime: time.Microsecond,
		},
		{
			Data:         2,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     3,
			RelativeTime: 33 * time.Microsecond,
		},
	}

	expected := []measurer.Measure{
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
			Data:         0,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     1,
			RelativeTime: 11 * time.Microsecond,
		},
		{
			Data:         1,
			Kind:         measurer.MeasureKindProcessed,
			Priority:     2,
			RelativeTime: 20 * time.Microsecond,
		},
		{
			Data:         1,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     2,
			RelativeTime: 25 * time.Microsecond,
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
			Kind:         measurer.MeasureKindReceived,
			Priority:     3,
			RelativeTime: 0,
		},
		{
			Data:         2,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     3,
			RelativeTime: 33 * time.Microsecond,
		},
	}

	sortByData(measures)
	require.Equal(t, expected, measures)
}

func TestSortByRelativeTime(t *testing.T) {
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

	expected := []measurer.Measure{
		{
			Data:         0,
			Kind:         measurer.MeasureKindReceived,
			Priority:     1,
			RelativeTime: 0,
		},
		{
			Data:         2,
			Kind:         measurer.MeasureKindReceived,
			Priority:     3,
			RelativeTime: 0,
		},
		{
			Data:         1,
			Kind:         measurer.MeasureKindReceived,
			Priority:     2,
			RelativeTime: time.Microsecond,
		},
		{
			Data:         0,
			Kind:         measurer.MeasureKindProcessed,
			Priority:     1,
			RelativeTime: 10 * time.Microsecond,
		},
		{
			Data:         0,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     1,
			RelativeTime: 11 * time.Microsecond,
		},
		{
			Data:         1,
			Kind:         measurer.MeasureKindProcessed,
			Priority:     2,
			RelativeTime: 20 * time.Microsecond,
		},
		{
			Data:         1,
			Kind:         measurer.MeasureKindCompleted,
			Priority:     2,
			RelativeTime: 25 * time.Microsecond,
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
	}

	sortByRelativeTime(measures)
	require.Equal(t, expected, measures)
}

package limit

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/durations"
	"github.com/akramarenkov/cqos/v2/limit/internal/research"
	"github.com/akramarenkov/cqos/v2/limit/internal/stress"

	"github.com/go-echarts/go-echarts/v2/charts"
	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
	"github.com/stretchr/testify/require"
)

func createTimeQuantitiesGraph(
	t *testing.T,
	relativeTimes []time.Duration,
	intervalsQuantity int,
	stressSystem bool,
) {
	quantities := research.CalcIntervalQuantities(relativeTimes, intervalsQuantity, 0)

	axisY, axisX := research.ConvertQuantityOverTimeToBarEcharts(quantities)

	chart := charts.NewBar()

	subtitle := fmt.Sprintf(
		"Total quantity: %d, "+
			"total duration: %s, "+
			"stress system: %t, "+
			"time: %s",
		len(relativeTimes),
		durations.CalcTotalDuration(relativeTimes),
		stressSystem,
		time.Now().Format(time.RFC3339),
	)

	chart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "Time quantities",
				Subtitle: subtitle,
			},
		),
	)

	chart.SetXAxis(axisX).AddSeries("quantities", axisY)

	baseName := "graph_time_quantities_" +
		strconv.Itoa(len(relativeTimes)) +
		"_stress_" +
		strconv.FormatBool(stressSystem)

	chartFile, err := os.Create(baseName + ".html")
	require.NoError(t, err)

	err = chart.Render(chartFile)
	require.NoError(t, err)
}

func createTimeDeviationsGraph(
	t *testing.T,
	relativeTimes []time.Duration,
	intervalsQuantity int,
	stressSystem bool,
) {
	deviations, min, max, avg := research.CalcSelfDeviations(relativeTimes, intervalsQuantity, 0)

	axisY, axisX := research.ConvertQuantityOverTimeToBarEcharts(deviations)

	chart := charts.NewBar()

	subtitle := fmt.Sprintf(
		"Total quantity: %d, "+
			"min: %s, "+
			"max: %s, "+
			"avg: %s, "+
			"stress system: %t, "+
			"time: %s",
		len(relativeTimes),
		min,
		max,
		avg,
		stressSystem,
		time.Now().Format(time.RFC3339),
	)

	chart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "Time deviations",
				Subtitle: subtitle,
			},
		),
	)

	chart.SetXAxis(axisX).AddSeries("deviations", axisY)

	baseName := "graph_time_deviations_" +
		strconv.Itoa(len(relativeTimes)) +
		"_stress_" +
		strconv.FormatBool(stressSystem)

	chartFile, err := os.Create(baseName + ".html")
	require.NoError(t, err)

	err = chart.Render(chartFile)
	require.NoError(t, err)
}

func testGraphTime(t *testing.T, quantity int, stressSystem bool) {
	if os.Getenv("CQOS_ENABLE_GRAPHS") == "" {
		t.SkipNow()
	}

	if stressSystem {
		stress, err := stress.New(0)
		require.NoError(t, err)

		defer stress.Stop()
	}

	relativeTimes := make([]time.Duration, quantity)

	startedAt := time.Now()

	for id := 0; id < quantity; id++ {
		relativeTimes[id] = time.Since(startedAt)
	}

	require.Equal(t, true, durations.IsSorted(relativeTimes))

	createTimeQuantitiesGraph(t, relativeTimes, 100, stressSystem)
	createTimeDeviationsGraph(t, relativeTimes, 100, stressSystem)
}

func TestGraphTime(t *testing.T) {
	testGraphTime(t, 1000, false)
	testGraphTime(t, 10000000, false)

	testGraphTime(t, 1000, true)
	testGraphTime(t, 10000000, true)
}

func createTickerTickQuantitiesGraph(
	t *testing.T,
	relativeTimes []time.Duration,
	duration time.Duration,
	buffered bool,
	stressSystem bool,
) {
	quantities := research.CalcIntervalQuantities(relativeTimes, 0, duration)

	axisY, axisX := research.ConvertQuantityOverTimeToBarEcharts(quantities)

	chart := charts.NewBar()

	subtitle := fmt.Sprintf(
		"Total quantity: %d, "+
			"ticker duration: %s, "+
			"total duration: {expected:  %s, actual: %s}, "+
			"buffered: %t, "+
			"stress system: %t, "+
			"time: %s",
		len(relativeTimes),
		duration,
		time.Duration(len(relativeTimes))*duration,
		durations.CalcTotalDuration(relativeTimes),
		buffered,
		stressSystem,
		time.Now().Format(time.RFC3339),
	)

	chart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "Ticker tick quantities over time",
				Subtitle: subtitle,
			},
		),
	)

	chart.SetXAxis(axisX).AddSeries("quantities", axisY)

	baseName := "graph_ticker_tick_quantities_" +
		strconv.Itoa(len(relativeTimes)) +
		"_ticker_duration_" +
		duration.String() +
		"_buffered_" +
		strconv.FormatBool(buffered) +
		"_stress_" +
		strconv.FormatBool(stressSystem)

	chartFile, err := os.Create(baseName + ".html")
	require.NoError(t, err)

	err = chart.Render(chartFile)
	require.NoError(t, err)
}

func createTickerTickDeviationsGraph(
	t *testing.T,
	relativeTimes []time.Duration,
	duration time.Duration,
	buffered bool,
	stressSystem bool,
) {
	deviations := research.CalcRelativeDeviations(relativeTimes, duration)

	axisY, axisX := research.ConvertRelativeDeviationsToBarEcharts(deviations)

	chart := charts.NewBar()

	subtitle := fmt.Sprintf(
		"Total quantity: %d, "+
			"ticker duration: %s, "+
			"buffered: %t, "+
			"stress system: %t, "+
			"time: %s",
		len(relativeTimes),
		duration,
		buffered,
		stressSystem,
		time.Now().Format(time.RFC3339),
	)

	chart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "Ticker tick deviations",
				Subtitle: subtitle,
			},
		),
	)

	chart.SetXAxis(axisX).AddSeries("deviations", axisY)

	baseName := "graph_ticker_tick_deviations_" +
		strconv.Itoa(len(relativeTimes)) +
		"_ticker_duration_" +
		duration.String() +
		"_buffered_" +
		strconv.FormatBool(buffered) +
		"_stress_" +
		strconv.FormatBool(stressSystem)

	chartFile, err := os.Create(baseName + ".html")
	require.NoError(t, err)

	err = chart.Render(chartFile)
	require.NoError(t, err)
}

func testGraphTicker(
	t *testing.T,
	quantity uint,
	duration time.Duration,
	buffered bool,
	stressSystem bool,
) {
	if os.Getenv("CQOS_ENABLE_GRAPHS") == "" {
		t.SkipNow()
	}

	if stressSystem {
		stress, err := stress.New(0)
		require.NoError(t, err)

		defer stress.Stop()
	}

	relativeTimes := make([]time.Duration, 0, quantity)

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	ticks := ticker.C

	if buffered {
		buffer := make(chan time.Time, quantity)
		defer close(buffer)

		ticks = buffer

		wg := &sync.WaitGroup{}
		defer wg.Wait()

		breaker := make(chan bool)
		defer close(breaker)

		wg.Add(1)

		go func() {
			defer wg.Done()

			for {
				select {
				case <-breaker:
					return
				case time := <-ticker.C:
					select {
					case <-breaker:
						return
					case buffer <- time:
					}
				}
			}
		}()
	}

	startedAt := time.Now()

	for range ticks {
		relativeTimes = append(relativeTimes, time.Since(startedAt))

		if uint(len(relativeTimes)) == quantity {
			break
		}
	}

	require.Equal(t, true, durations.IsSorted(relativeTimes))

	createTickerTickQuantitiesGraph(t, relativeTimes, duration, buffered, stressSystem)
	createTickerTickDeviationsGraph(t, relativeTimes, duration, buffered, stressSystem)
}

func TestGraphTicker(t *testing.T) {
	testGraphTicker(t, 100, time.Second, false, false)
	testGraphTicker(t, 100, 100*time.Millisecond, false, false)
	testGraphTicker(t, 1000, 10*time.Millisecond, false, false)
	testGraphTicker(t, 1000, time.Millisecond, false, false)
	testGraphTicker(t, 1000, 100*time.Microsecond, false, false)
	testGraphTicker(t, 1000, 10*time.Microsecond, false, false)
	testGraphTicker(t, 1000, time.Microsecond, false, false)
	testGraphTicker(t, 1000, 100*time.Nanosecond, false, false)
	testGraphTicker(t, 1000, 10*time.Nanosecond, false, false)
	testGraphTicker(t, 1000, time.Nanosecond, false, false)

	testGraphTicker(t, 100, time.Second, true, false)
	testGraphTicker(t, 100, 100*time.Millisecond, true, false)
	testGraphTicker(t, 1000, 10*time.Millisecond, true, false)
	testGraphTicker(t, 1000, time.Millisecond, true, false)
	testGraphTicker(t, 1000, 100*time.Microsecond, true, false)
	testGraphTicker(t, 1000, 10*time.Microsecond, true, false)
	testGraphTicker(t, 1000, time.Microsecond, true, false)
	testGraphTicker(t, 1000, 100*time.Nanosecond, true, false)
	testGraphTicker(t, 1000, 10*time.Nanosecond, true, false)
	testGraphTicker(t, 1000, time.Nanosecond, true, false)

	testGraphTicker(t, 100, time.Second, false, true)
	testGraphTicker(t, 100, 100*time.Millisecond, false, true)
	testGraphTicker(t, 1000, 10*time.Millisecond, false, true)
	testGraphTicker(t, 1000, time.Millisecond, false, true)
	testGraphTicker(t, 1000, 100*time.Microsecond, false, true)
	testGraphTicker(t, 1000, 10*time.Microsecond, false, true)
	testGraphTicker(t, 1000, time.Microsecond, false, true)
	testGraphTicker(t, 1000, 100*time.Nanosecond, false, true)
	testGraphTicker(t, 1000, 10*time.Nanosecond, false, true)
	testGraphTicker(t, 1000, time.Nanosecond, false, true)

	testGraphTicker(t, 100, time.Second, true, true)
	testGraphTicker(t, 100, 100*time.Millisecond, true, true)
	testGraphTicker(t, 1000, 10*time.Millisecond, true, true)
	testGraphTicker(t, 1000, time.Millisecond, true, true)
	testGraphTicker(t, 1000, 100*time.Microsecond, true, true)
	testGraphTicker(t, 1000, 10*time.Microsecond, true, true)
	testGraphTicker(t, 1000, time.Microsecond, true, true)
	testGraphTicker(t, 1000, 100*time.Nanosecond, true, true)
	testGraphTicker(t, 1000, 10*time.Nanosecond, true, true)
	testGraphTicker(t, 1000, time.Nanosecond, true, true)
}

func createQuantitiesGraph(
	t *testing.T,
	relativeTimes []time.Duration,
	limit Rate,
	stressSystem bool,
	kind string,
) {
	quantities := research.CalcIntervalQuantities(relativeTimes, 0, limit.Interval)

	axisY, axisX := research.ConvertQuantityOverTimeToBarEcharts(quantities)

	chart := charts.NewBar()

	subtitle := fmt.Sprintf(
		"Total quantity: %d, "+
			"limit: {quantity: %d, interval: %s}, "+
			"total duration: {expected:  %s, actual: %s}, "+
			"stress system: %t, "+
			"kind: %s, "+
			"time: %s",
		len(relativeTimes),
		limit.Quantity,
		limit.Interval,
		time.Duration(len(relativeTimes))*limit.Interval/time.Duration(limit.Quantity),
		durations.CalcTotalDuration(relativeTimes),
		stressSystem,
		kind,
		time.Now().Format(time.RFC3339),
	)

	chart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "Quantities over time",
				Subtitle: subtitle,
			},
		),
	)

	chart.SetXAxis(axisX).AddSeries("quantities", axisY)

	baseName := "graph_quantities_" +
		strconv.Itoa(len(relativeTimes)) +
		"_limit_quantity_" +
		strconv.Itoa(int(limit.Quantity)) +
		"_limit_interval_" +
		limit.Interval.String() +
		"_stress_" +
		strconv.FormatBool(stressSystem)

	chartFile, err := os.Create(baseName + "_" + kind + ".html")
	require.NoError(t, err)

	err = chart.Render(chartFile)
	require.NoError(t, err)
}

func createDeviationsGraph(
	t *testing.T,
	relativeTimes []time.Duration,
	limit Rate,
	stressSystem bool,
	kind string,
) {
	flattenLimit, done := limit.Flatten()
	require.Equal(t, true, done)

	deviations := research.CalcRelativeDeviations(relativeTimes, flattenLimit.Interval)

	axisY, axisX := research.ConvertRelativeDeviationsToBarEcharts(deviations)

	chart := charts.NewBar()

	subtitle := fmt.Sprintf(
		"Total quantity: %d, "+
			"limit {quantity: %d, interval: %s}, "+
			"stress system: %t, "+
			"kind: %s, "+
			"time: %s",
		len(relativeTimes),
		limit.Quantity,
		limit.Interval,
		stressSystem,
		kind,
		time.Now().Format(time.RFC3339),
	)

	chart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "Deviations",
				Subtitle: subtitle,
			},
		),
	)

	chart.SetXAxis(axisX).AddSeries("deviations", axisY)

	baseName := "graph_deviations_" +
		strconv.Itoa(len(relativeTimes)) +
		"_limit_quantity_" +
		strconv.Itoa(int(limit.Quantity)) +
		"_limit_interval_" +
		limit.Interval.String() +
		"_stress_" +
		strconv.FormatBool(stressSystem)

	chartFile, err := os.Create(baseName + "_" + kind + ".html")
	require.NoError(t, err)

	err = chart.Render(chartFile)
	require.NoError(t, err)
}

func testGraphDisciplineSynthetic(
	t *testing.T,
	quantity uint,
	limit Rate,
	stressSystem bool,
) {
	if os.Getenv("CQOS_ENABLE_GRAPHS") == "" {
		t.SkipNow()
	}

	if stressSystem {
		stress, err := stress.New(0)
		require.NoError(t, err)

		defer stress.Stop()
	}

	input := make(chan uint, quantity)
	relativeTimes := make([]time.Duration, 0, quantity)

	for stage := uint(1); stage <= quantity; stage++ {
		input <- stage
	}

	close(input)

	opts := Opts[uint]{
		Input:     input,
		Limit:     limit,
		OutputCap: quantity,
	}

	startedAt := time.Now()

	discipline, err := New(opts)
	require.NoError(t, err)

	for range discipline.Output() {
		relativeTimes = append(relativeTimes, time.Since(startedAt))
	}

	require.Equal(t, true, durations.IsSorted(relativeTimes))

	createQuantitiesGraph(t, relativeTimes, limit, stressSystem, "synthetic")
	createDeviationsGraph(t, relativeTimes, limit, stressSystem, "synthetic")
}

func TestGraphDisciplineSynthetic(t *testing.T) {
	testGraphDisciplineSynthetic(
		t,
		10000,
		Rate{Interval: time.Millisecond, Quantity: 1},
		false,
	)
	testGraphDisciplineSynthetic(
		t,
		10000,
		Rate{Interval: time.Millisecond, Quantity: 1},
		true,
	)

	testGraphDisciplineSynthetic(
		t,
		10000,
		Rate{Interval: 10 * time.Millisecond, Quantity: 10},
		false,
	)
	testGraphDisciplineSynthetic(
		t,
		10000,
		Rate{Interval: 10 * time.Millisecond, Quantity: 10},
		true,
	)

	testGraphDisciplineSynthetic(
		t,
		10000,
		Rate{Interval: 100 * time.Millisecond, Quantity: 100},
		false,
	)
	testGraphDisciplineSynthetic(
		t,
		10000,
		Rate{Interval: 100 * time.Millisecond, Quantity: 100},
		true,
	)

	testGraphDisciplineSynthetic(
		t,
		10000,
		Rate{Interval: 1000 * time.Millisecond, Quantity: 1000},
		false,
	)
	testGraphDisciplineSynthetic(
		t,
		10000,
		Rate{Interval: 1000 * time.Millisecond, Quantity: 1000},
		true,
	)
}

func testGraphDisciplineRegular(
	t *testing.T,
	quantity uint,
	limit Rate,
	stressSystem bool,
) {
	if os.Getenv("CQOS_ENABLE_GRAPHS") == "" {
		t.SkipNow()
	}

	if stressSystem {
		stress, err := stress.New(0)
		require.NoError(t, err)

		defer stress.Stop()
	}

	input := make(chan uint)

	opts := Opts[uint]{
		Input: input,
		Limit: limit,
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	wg := &sync.WaitGroup{}

	wg.Add(2)

	inSequence := make([]uint, 0, quantity)
	outSequence := make([]uint, 0, quantity)
	relativeTimes := make([]time.Duration, 0, quantity)

	startedAt := time.Now()

	go func() {
		defer wg.Done()
		defer close(input)

		for stage := uint(1); stage <= quantity; stage++ {
			inSequence = append(inSequence, stage)

			input <- stage
		}
	}()

	go func() {
		defer wg.Done()

		for item := range discipline.Output() {
			relativeTimes = append(relativeTimes, time.Since(startedAt))
			outSequence = append(outSequence, item)
		}
	}()

	wg.Wait()

	require.Equal(t, inSequence, outSequence)
	require.Equal(t, true, durations.IsSorted(relativeTimes))

	createQuantitiesGraph(t, relativeTimes, limit, stressSystem, "regular")
	createDeviationsGraph(t, relativeTimes, limit, stressSystem, "regular")
}

func TestGraphDisciplineRegular1000(t *testing.T) {
	testGraphDisciplineRegular(
		t,
		10000,
		Rate{Interval: 1 * time.Millisecond, Quantity: 1},
		false,
	)
	testGraphDisciplineRegular(
		t,
		10000,
		Rate{Interval: 1 * time.Millisecond, Quantity: 1},
		true,
	)

	testGraphDisciplineRegular(
		t,
		10000,
		Rate{Interval: 10 * time.Millisecond, Quantity: 10},
		false,
	)
	testGraphDisciplineRegular(
		t,
		10000,
		Rate{Interval: 10 * time.Millisecond, Quantity: 10},
		true,
	)

	testGraphDisciplineRegular(
		t,
		10000,
		Rate{Interval: 100 * time.Millisecond, Quantity: 100},
		false,
	)
	testGraphDisciplineRegular(
		t,
		10000,
		Rate{Interval: 100 * time.Millisecond, Quantity: 100},
		true,
	)

	testGraphDisciplineRegular(
		t,
		10000,
		Rate{Interval: 1000 * time.Millisecond, Quantity: 1000},
		false,
	)
	testGraphDisciplineRegular(
		t,
		10000,
		Rate{Interval: 1000 * time.Millisecond, Quantity: 1000},
		true,
	)
}

func TestGraphDisciplineRegular10000(t *testing.T) {
	testGraphDisciplineRegular(
		t,
		100000,
		Rate{Interval: 1 * time.Millisecond, Quantity: 10},
		false,
	)
	testGraphDisciplineRegular(
		t,
		100000,
		Rate{Interval: 1 * time.Millisecond, Quantity: 10},
		true,
	)

	testGraphDisciplineRegular(
		t,
		100000,
		Rate{Interval: 10 * time.Millisecond, Quantity: 100},
		false,
	)
	testGraphDisciplineRegular(
		t,
		100000,
		Rate{Interval: 10 * time.Millisecond, Quantity: 100},
		true,
	)

	testGraphDisciplineRegular(
		t,
		100000,
		Rate{Interval: 100 * time.Millisecond, Quantity: 1000},
		false,
	)
	testGraphDisciplineRegular(
		t,
		100000,
		Rate{Interval: 100 * time.Millisecond, Quantity: 1000},
		true,
	)

	testGraphDisciplineRegular(
		t,
		100000,
		Rate{Interval: 1000 * time.Millisecond, Quantity: 10000},
		false,
	)
	testGraphDisciplineRegular(
		t,
		100000,
		Rate{Interval: 1000 * time.Millisecond, Quantity: 10000},
		true,
	)
}

func TestGraphDisciplineRegular100000(t *testing.T) {
	testGraphDisciplineRegular(
		t,
		1000000,
		Rate{Interval: 1 * time.Millisecond, Quantity: 100},
		false,
	)
	testGraphDisciplineRegular(
		t,
		1000000,
		Rate{Interval: 1 * time.Millisecond, Quantity: 100},
		true,
	)

	testGraphDisciplineRegular(
		t,
		1000000,
		Rate{Interval: 10 * time.Millisecond, Quantity: 1000},
		false,
	)
	testGraphDisciplineRegular(
		t,
		1000000,
		Rate{Interval: 10 * time.Millisecond, Quantity: 1000},
		true,
	)

	testGraphDisciplineRegular(
		t,
		1000000,
		Rate{Interval: 100 * time.Millisecond, Quantity: 10000},
		false,
	)
	testGraphDisciplineRegular(
		t,
		1000000,
		Rate{Interval: 100 * time.Millisecond, Quantity: 10000},
		true,
	)

	testGraphDisciplineRegular(
		t,
		1000000,
		Rate{Interval: 1000 * time.Millisecond, Quantity: 100000},
		false,
	)
	testGraphDisciplineRegular(
		t,
		1000000,
		Rate{Interval: 1000 * time.Millisecond, Quantity: 100000},
		true,
	)
}

func TestGraphDisciplineRegular1000000(t *testing.T) {
	testGraphDisciplineRegular(
		t,
		10000000,
		Rate{Interval: 1 * time.Millisecond, Quantity: 1000},
		false,
	)
	testGraphDisciplineRegular(
		t,
		10000000,
		Rate{Interval: 1 * time.Millisecond, Quantity: 1000},
		true,
	)

	testGraphDisciplineRegular(
		t,
		10000000,
		Rate{Interval: 10 * time.Millisecond, Quantity: 10000},
		false,
	)
	testGraphDisciplineRegular(
		t,
		10000000,
		Rate{Interval: 10 * time.Millisecond, Quantity: 10000},
		true,
	)

	testGraphDisciplineRegular(
		t,
		10000000,
		Rate{Interval: 100 * time.Millisecond, Quantity: 100000},
		false,
	)
	testGraphDisciplineRegular(
		t,
		10000000,
		Rate{Interval: 100 * time.Millisecond, Quantity: 100000},
		true,
	)

	testGraphDisciplineRegular(
		t,
		10000000,
		Rate{Interval: 1000 * time.Millisecond, Quantity: 1000000},
		false,
	)
	testGraphDisciplineRegular(
		t,
		10000000,
		Rate{Interval: 1000 * time.Millisecond, Quantity: 1000000},
		true,
	)
}
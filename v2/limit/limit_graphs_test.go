package limit

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/consts"
	"github.com/akramarenkov/cqos/v2/internal/durations"
	"github.com/akramarenkov/cqos/v2/limit/internal/research"
	"github.com/akramarenkov/cqos/v2/limit/internal/stress"

	"github.com/go-echarts/go-echarts/v2/charts"
	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
	"github.com/stretchr/testify/require"
)

func TestGraphTime(t *testing.T) {
	// we use the same calc (graph) interval without stress system and
	// under stress system
	quantitiesInterval1e3, deviationsInterval1e3 := testGraphTime(
		t,
		1e3,
		100,
		0,
		100,
		0,
		false,
	)

	quantitiesInterval1e7, deviationsInterval1e7 := testGraphTime(
		t,
		1e7,
		100,
		0,
		100,
		0,
		false,
	)

	testGraphTime(
		t,
		1e3,
		0,
		quantitiesInterval1e3,
		0,
		deviationsInterval1e3,
		true,
	)

	testGraphTime(
		t,
		1e7,
		0,
		quantitiesInterval1e7,
		0,
		deviationsInterval1e7,
		true,
	)
}

func testGraphTime(
	t *testing.T,
	quantity int,
	quantitiesIntervalsQuantity int,
	quantitiesInterval time.Duration,
	deviationsIntervalsQuantity int,
	deviationsInterval time.Duration,
	stressSystem bool,
) (time.Duration, time.Duration) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	if stressSystem {
		stress, err := stress.New(0, 0)
		require.NoError(t, err)

		defer stress.Stop()
	}

	relativeTimes := make([]time.Duration, quantity)

	startedAt := time.Now()

	for id := 0; id < quantity; id++ {
		relativeTimes[id] = time.Since(startedAt)
	}

	require.Equal(t, true, durations.IsSorted(relativeTimes))

	quantitiesInterval = createTimeQuantitiesGraph(
		t,
		relativeTimes,
		quantitiesIntervalsQuantity,
		quantitiesInterval,
		stressSystem,
	)

	deviationsInterval = createTimeDeviationsGraph(
		t,
		relativeTimes,
		deviationsIntervalsQuantity,
		deviationsInterval,
		stressSystem,
	)

	return quantitiesInterval, deviationsInterval
}

func createTimeQuantitiesGraph(
	t *testing.T,
	relativeTimes []time.Duration,
	calcIntervalsQuantity int,
	calcInterval time.Duration,
	stressSystem bool,
) time.Duration {
	quantities, calcInterval := research.CalcIntervalQuantities(
		relativeTimes,
		calcIntervalsQuantity,
		calcInterval,
	)

	axisY, axisX := research.ConvertQuantityOverTimeToBarEcharts(quantities)

	subtitleAddition := fmt.Sprintf(
		"total duration: %s",
		durations.CalcTotalDuration(relativeTimes),
	)

	createGraph(
		t,
		"Time quantities",
		subtitleAddition,
		"time_quantities",
		"quantities",
		len(relativeTimes),
		calcInterval.String(),
		stressSystem,
		axisY,
		axisX,
	)

	return calcInterval
}

func createTimeDeviationsGraph(
	t *testing.T,
	relativeTimes []time.Duration,
	calcIntervalsQuantity int,
	calcInterval time.Duration,
	stressSystem bool,
) time.Duration {
	deviations, calcInterval, min, max, avg := research.CalcSelfDeviations(
		relativeTimes,
		calcIntervalsQuantity,
		calcInterval,
	)

	axisY, axisX := research.ConvertQuantityOverTimeToBarEcharts(deviations)

	subtitleAddition := fmt.Sprintf(
		"min: %s, "+
			"max: %s, "+
			"avg: %s",
		min,
		max,
		avg,
	)

	createGraph(
		t,
		"Time deviations",
		subtitleAddition,
		"time_deviations",
		"deviations",
		len(relativeTimes),
		calcInterval.String(),
		stressSystem,
		axisY,
		axisX,
	)

	return calcInterval
}

func TestGraphExtrapolatedDuration(t *testing.T) {
	testGraphExtrapolatedDuration(
		t,
		1e2,
		false,
	)
	testGraphExtrapolatedDuration(
		t,
		1e3,
		false,
	)
	testGraphExtrapolatedDuration(
		t,
		1e4,
		false,
	)
	testGraphExtrapolatedDuration(
		t,
		1e5,
		false,
	)

	testGraphExtrapolatedDuration(
		t,
		1e2,
		true,
	)
	testGraphExtrapolatedDuration(
		t,
		1e3,
		true,
	)
	testGraphExtrapolatedDuration(
		t,
		1e4,
		true,
	)
	testGraphExtrapolatedDuration(
		t,
		1e5,
		true,
	)
}

func testGraphExtrapolatedDuration(
	t *testing.T,
	quantity int,
	stressSystem bool,
) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	if stressSystem {
		stress, err := stress.New(0, 0)
		require.NoError(t, err)

		defer stress.Stop()
	}

	durations := make([]time.Duration, quantity)

	startedAt := time.Now()

	for id := 0; id < quantity; id++ {
		durations[id] = time.Duration(float64(time.Since(startedAt)) * float64(quantity) / float64(id+1))
	}

	createExtrapolatedDurationGraph(t, durations, stressSystem)
	createExtrapolatedDurationDeviationsGraph(t, durations, stressSystem)
}

func createExtrapolatedDurationGraph(
	t *testing.T,
	durations []time.Duration,
	stressSystem bool,
) {
	axisY, axisX := research.ConvertDurationsToBarEcharts(durations)

	subtitleAddition := fmt.Sprintf(
		"expected (last) duration: %s",
		research.GetExpectedExtrapolatedDuration(durations),
	)

	createGraph(
		t,
		"Extrapolated total duration",
		subtitleAddition,
		"duration_values",
		"durations",
		len(durations),
		"1 measurement",
		stressSystem,
		axisY,
		axisX,
	)
}

func createExtrapolatedDurationDeviationsGraph(
	t *testing.T,
	durations []time.Duration,
	stressSystem bool,
) {
	deviations, expected := research.CalcExtrapolatedDurationDeviations(durations)

	axisY, axisX := research.ConvertDurationDeviationsToBarEcharts(deviations)

	subtitleAddition := fmt.Sprintf(
		"expected (last) duration: %s",
		expected,
	)

	createGraph(
		t,
		"Extrapolated total duration deviations",
		subtitleAddition,
		"duration_deviations",
		"deviations",
		len(durations),
		"1 measurement",
		stressSystem,
		axisY,
		axisX,
	)
}

func TestGraphTicker(t *testing.T) {
	testGraphTicker(t, 1e2, time.Second, false)
	testGraphTicker(t, 1e2, 100*time.Millisecond, false)
	testGraphTicker(t, 1e3, 10*time.Millisecond, false)
	testGraphTicker(t, 1e3, time.Millisecond, false)
	testGraphTicker(t, 1e3, 100*time.Microsecond, false)
	testGraphTicker(t, 1e3, 10*time.Microsecond, false)
	testGraphTicker(t, 1e3, time.Microsecond, false)
	testGraphTicker(t, 1e3, 100*time.Nanosecond, false)
	testGraphTicker(t, 1e3, 10*time.Nanosecond, false)
	testGraphTicker(t, 1e3, time.Nanosecond, false)

	testGraphTicker(t, 1e2, time.Second, true)
	testGraphTicker(t, 1e2, 100*time.Millisecond, true)
	testGraphTicker(t, 1e3, 10*time.Millisecond, true)
	testGraphTicker(t, 1e3, time.Millisecond, true)
	testGraphTicker(t, 1e3, 100*time.Microsecond, true)
	testGraphTicker(t, 1e3, 10*time.Microsecond, true)
	testGraphTicker(t, 1e3, time.Microsecond, true)
	testGraphTicker(t, 1e3, 100*time.Nanosecond, true)
	testGraphTicker(t, 1e3, 10*time.Nanosecond, true)
	testGraphTicker(t, 1e3, time.Nanosecond, true)
}

func testGraphTicker(
	t *testing.T,
	quantity int,
	tickerDuration time.Duration,
	stressSystem bool,
) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	if stressSystem {
		stress, err := stress.New(0, 0)
		require.NoError(t, err)

		defer stress.Stop()
	}

	relativeTimes := make([]time.Duration, 0, quantity)

	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	startedAt := time.Now()

	for range ticker.C {
		relativeTimes = append(relativeTimes, time.Since(startedAt))

		if len(relativeTimes) == quantity {
			break
		}
	}

	createTickerQuantitiesGraph(t, relativeTimes, tickerDuration, stressSystem)
	createTickerDeviationsGraph(t, relativeTimes, tickerDuration, stressSystem)
}

func createTickerQuantitiesGraph(
	t *testing.T,
	relativeTimes []time.Duration,
	tickerDuration time.Duration,
	stressSystem bool,
) {
	quantities, calcInterval := research.CalcIntervalQuantities(relativeTimes, 0, tickerDuration)

	axisY, axisX := research.ConvertQuantityOverTimeToBarEcharts(quantities)

	expectedDuration := time.Duration(len(relativeTimes)) * tickerDuration

	subtitleAddition := fmt.Sprintf(
		"ticker duration: %s, "+
			formatTotalDuration(expectedDuration, relativeTimes),
		tickerDuration,
	)

	fileNameAddition := "ticker_tick_quantities_" +
		"ticker_duration_" +
		tickerDuration.String()

	createGraph(
		t,
		"Ticker tick quantities over time",
		subtitleAddition,
		fileNameAddition,
		"quantities",
		len(relativeTimes),
		calcInterval.String(),
		stressSystem,
		axisY,
		axisX,
	)
}

func createTickerDeviationsGraph(
	t *testing.T,
	relativeTimes []time.Duration,
	tickerDuration time.Duration,
	stressSystem bool,
) {
	deviations := research.CalcRelativeDeviations(relativeTimes, tickerDuration)

	axisY, axisX := research.ConvertRelativeDeviationsToBarEcharts(deviations)

	subtitleAddition := fmt.Sprintf(
		"ticker duration: %s",
		tickerDuration,
	)

	fileNameAddition := "ticker_tick_deviations_" +
		"ticker_duration_" +
		tickerDuration.String()

	createGraph(
		t,
		"Ticker tick deviations",
		subtitleAddition,
		fileNameAddition,
		"deviations",
		len(relativeTimes),
		"1%",
		stressSystem,
		axisY,
		axisX,
	)
}

func TestGraphSleep(t *testing.T) {
	testGraphSleep(t, 1e2, time.Second, false)
	testGraphSleep(t, 1e2, 100*time.Millisecond, false)
	testGraphSleep(t, 1e3, 10*time.Millisecond, false)
	testGraphSleep(t, 1e3, time.Millisecond, false)
	testGraphSleep(t, 1e3, 100*time.Microsecond, false)
	testGraphSleep(t, 1e3, 10*time.Microsecond, false)
	testGraphSleep(t, 1e3, time.Microsecond, false)
	testGraphSleep(t, 1e3, 100*time.Nanosecond, false)
	testGraphSleep(t, 1e3, 10*time.Nanosecond, false)
	testGraphSleep(t, 1e3, time.Nanosecond, false)

	testGraphSleep(t, 1e2, time.Second, true)
	testGraphSleep(t, 1e2, 100*time.Millisecond, true)
	testGraphSleep(t, 1e3, 10*time.Millisecond, true)
	testGraphSleep(t, 1e3, time.Millisecond, true)
	testGraphSleep(t, 1e3, 100*time.Microsecond, true)
	testGraphSleep(t, 1e3, 10*time.Microsecond, true)
	testGraphSleep(t, 1e3, time.Microsecond, true)
	testGraphSleep(t, 1e3, 100*time.Nanosecond, true)
	testGraphSleep(t, 1e3, 10*time.Nanosecond, true)
	testGraphSleep(t, 1e3, time.Nanosecond, true)
}

func testGraphSleep(
	t *testing.T,
	quantity int,
	duration time.Duration,
	stressSystem bool,
) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	if stressSystem {
		stress, err := stress.New(0, 0)
		require.NoError(t, err)

		defer stress.Stop()
	}

	relativeTimes := make([]time.Duration, quantity)

	startedAt := time.Now()

	for id := 0; id < quantity; id++ {
		time.Sleep(duration)

		relativeTimes[id] = time.Since(startedAt)
	}

	createSleepQuantitiesGraph(t, relativeTimes, duration, stressSystem)
	createSleepDeviationsGraph(t, relativeTimes, duration, stressSystem)
}

func createSleepQuantitiesGraph(
	t *testing.T,
	relativeTimes []time.Duration,
	duration time.Duration,
	stressSystem bool,
) {
	quantities, calcInterval := research.CalcIntervalQuantities(relativeTimes, 0, duration)

	axisY, axisX := research.ConvertQuantityOverTimeToBarEcharts(quantities)

	expectedDuration := time.Duration(len(relativeTimes)) * duration

	subtitleAddition := fmt.Sprintf(
		"sleep duration: %s, "+
			formatTotalDuration(expectedDuration, relativeTimes),
		duration,
	)

	fileNameAddition := "sleep_quantities_" +
		"sleep_duration_" +
		duration.String()

	createGraph(
		t,
		"Sleep quantities over time",
		subtitleAddition,
		fileNameAddition,
		"quantities",
		len(relativeTimes),
		calcInterval.String(),
		stressSystem,
		axisY,
		axisX,
	)
}

func createSleepDeviationsGraph(
	t *testing.T,
	relativeTimes []time.Duration,
	duration time.Duration,
	stressSystem bool,
) {
	deviations := research.CalcRelativeDeviations(relativeTimes, duration)

	axisY, axisX := research.ConvertRelativeDeviationsToBarEcharts(deviations)

	subtitleAddition := fmt.Sprintf(
		"sleep duration: %s",
		duration,
	)

	fileNameAddition := "sleep_deviations_" +
		"sleep_duration_" +
		duration.String()

	createGraph(
		t,
		"Sleep deviations",
		subtitleAddition,
		fileNameAddition,
		"deviations",
		len(relativeTimes),
		"1%",
		stressSystem,
		axisY,
		axisX,
	)
}

func TestGraphDiscipline(t *testing.T) {
	testGraphDiscipline(
		t,
		1e4+1,
		Rate{Interval: time.Second, Quantity: 1e3},
		false,
	)

	testGraphDiscipline(
		t,
		1e4+1,
		Rate{Interval: time.Millisecond, Quantity: 1},
		false,
	)

	testGraphDiscipline(
		t,
		1e5+1,
		Rate{Interval: time.Second, Quantity: 1e4},
		false,
	)

	testGraphDiscipline(
		t,
		1e5+1,
		Rate{Interval: 100 * time.Microsecond, Quantity: 1},
		false,
	)

	testGraphDiscipline(
		t,
		1e6+1,
		Rate{Interval: time.Second, Quantity: 1e5},
		false,
	)

	testGraphDiscipline(
		t,
		1e6+1,
		Rate{Interval: 10 * time.Microsecond, Quantity: 1},
		false,
	)

	testGraphDiscipline(
		t,
		1e7+1,
		Rate{Interval: time.Second, Quantity: 1e6},
		false,
	)

	testGraphDiscipline(
		t,
		2.6e7+1,
		Rate{Interval: time.Second, Quantity: 2.6e6},
		false,
	)

	testGraphDiscipline(
		t,
		3e7+1,
		Rate{Interval: time.Second, Quantity: 3e6},
		false,
	)

	testGraphDiscipline(
		t,
		2.6e7+1,
		Rate{Interval: time.Second, Quantity: 2.6e6},
		true,
	)
}

func testGraphDiscipline(
	t *testing.T,
	quantity int,
	limit Rate,
	stressSystem bool,
) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	if stressSystem {
		stress, err := stress.New(0, 0)
		require.NoError(t, err)

		defer stress.Stop()
	}

	input := make(chan int)

	opts := Opts[int]{
		Input: input,
		Limit: limit,
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	relativeTimes := make([]time.Duration, 0, quantity)

	startedAt := time.Now()

	go func() {
		defer close(input)

		for stage := 0; stage < quantity; stage++ {
			input <- stage
		}
	}()

	for range discipline.Output() {
		relativeTimes = append(relativeTimes, time.Since(startedAt))
	}

	createQuantitiesGraph(t, relativeTimes, limit, stressSystem)
	createDeviationsGraph(t, relativeTimes, limit, stressSystem)
}

func createQuantitiesGraph(
	t *testing.T,
	relativeTimes []time.Duration,
	limit Rate,
	stressSystem bool,
) {
	quantities, calcInterval := research.CalcIntervalQuantities(relativeTimes, 100, 0)

	axisY, axisX := research.ConvertQuantityOverTimeToBarEcharts(quantities)

	expectedDuration := (time.Duration(len(relativeTimes)) * limit.Interval) / time.Duration(limit.Quantity)

	subtitleAddition := fmt.Sprintf(
		"limit: {quantity: %d, interval: %s}, "+
			formatTotalDuration(expectedDuration, relativeTimes),
		limit.Quantity,
		limit.Interval,
	)

	fileNameAddition := "quantities_" +
		"limit_quantity_" +
		strconv.Itoa(int(limit.Quantity)) +
		"_limit_interval_" +
		limit.Interval.String()

	createGraph(
		t,
		"Quantities over time",
		subtitleAddition,
		fileNameAddition,
		"quantities",
		len(relativeTimes),
		calcInterval.String(),
		stressSystem,
		axisY,
		axisX,
	)
}

func createDeviationsGraph(
	t *testing.T,
	relativeTimes []time.Duration,
	limit Rate,
	stressSystem bool,
) {
	flatten, err := limit.Flatten()
	require.NoError(t, err)

	deviations := research.CalcRelativeDeviations(relativeTimes, flatten.Interval)

	axisY, axisX := research.ConvertRelativeDeviationsToBarEcharts(deviations)

	subtitleAddition := fmt.Sprintf(
		"limit: {quantity: %d, interval: %s}, "+
			"flatten: {quantity: %d, interval: %s}",
		limit.Quantity,
		limit.Interval,
		flatten.Quantity,
		flatten.Interval,
	)

	fileNameAddition := "deviations_" +
		"limit_quantity_" +
		strconv.Itoa(int(limit.Quantity)) +
		"_limit_interval_" +
		limit.Interval.String()

	createGraph(
		t,
		"Deviations",
		subtitleAddition,
		fileNameAddition,
		"deviations",
		len(relativeTimes),
		"1%",
		stressSystem,
		axisY,
		axisX,
	)
}

func createGraph(
	t *testing.T,
	title string,
	subtitleAddition string,
	fileNameAddition string,
	seriesName string,
	totalQuantity int,
	graphInterval string,
	stressSystem bool,
	series []chartsopts.BarData,
	abscissa interface{},
) {
	subtitle := fmt.Sprintf(
		"Total quantity: %d, "+
			"graph interval: %s, "+
			subtitleAddition+", "+
			"stress system: %t, "+
			"time: %s",
		totalQuantity,
		graphInterval,
		stressSystem,
		time.Now().Format(time.RFC3339),
	)

	fileName := "graph_" +
		strconv.Itoa(totalQuantity) +
		"_" +
		fileNameAddition +
		"_stress_" +
		strconv.FormatBool(stressSystem) +
		".html"

	createBarGraph(
		t,
		title,
		subtitle,
		fileName,
		seriesName,
		series,
		abscissa,
	)
}

func createBarGraph(
	t *testing.T,
	title string,
	subtitle string,
	fileName string,
	seriesName string,
	series []chartsopts.BarData,
	abscissa interface{},
) {
	if len(series) == 0 {
		return
	}

	chart := charts.NewBar()

	chart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    title,
				Subtitle: subtitle,
			},
		),
	)

	chart.SetXAxis(abscissa).AddSeries(seriesName, series)

	file, err := os.Create(fileName)
	require.NoError(t, err)

	err = chart.Render(file)
	require.NoError(t, err)
}

func formatTotalDuration(expected time.Duration, relativeTimes []time.Duration) string {
	out := fmt.Sprintf(
		"total duration: {expected:  %s, actual: %s}",
		expected,
		durations.CalcTotalDuration(relativeTimes),
	)

	return out
}

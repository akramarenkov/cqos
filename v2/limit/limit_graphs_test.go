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

func TestGraphTicker(t *testing.T) {
	testGraphTicker(t, 1e2, time.Second, false)
	testGraphTicker(t, 1e2, 100*time.Millisecond, false)
	testGraphTicker(t, 1e3, 10*time.Millisecond, false)
	testGraphTicker(t, 1e3, 5*time.Millisecond, false)
	testGraphTicker(t, 1e3, 4*time.Millisecond, false)
	testGraphTicker(t, 1e3, 3*time.Millisecond, false)
	testGraphTicker(t, 1e3, 2*time.Millisecond, false)
	testGraphTicker(t, 1e3, time.Millisecond, false)
	testGraphTicker(t, 1e3, 100*time.Microsecond, false)

	testGraphTicker(t, 1e2, time.Second, true)
	testGraphTicker(t, 1e2, 100*time.Millisecond, true)
	testGraphTicker(t, 1e3, 10*time.Millisecond, true)
	testGraphTicker(t, 1e3, 5*time.Millisecond, true)
	testGraphTicker(t, 1e3, 4*time.Millisecond, true)
	testGraphTicker(t, 1e3, 3*time.Millisecond, true)
	testGraphTicker(t, 1e3, 2*time.Millisecond, true)
	testGraphTicker(t, 1e3, time.Millisecond, true)
	testGraphTicker(t, 1e3, 100*time.Microsecond, true)
}

func testGraphTicker(
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

	relativeTimes := make([]time.Duration, 0, quantity)

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	startedAt := time.Now()

	for range ticker.C {
		relativeTimes = append(relativeTimes, time.Since(startedAt))

		if len(relativeTimes) == quantity {
			break
		}
	}

	createDelayerQuantitiesGraph(t, "Ticker", "ticker", relativeTimes, duration, stressSystem)
	createDelayerDeviationsGraph(t, "Ticker", "ticker", relativeTimes, duration, stressSystem)
}

func TestGraphSleep(t *testing.T) {
	testGraphSleep(t, 1e2, time.Second, false)
	testGraphSleep(t, 1e2, 100*time.Millisecond, false)
	testGraphSleep(t, 1e3, 10*time.Millisecond, false)
	testGraphSleep(t, 1e3, 5*time.Millisecond, false)
	testGraphSleep(t, 1e3, 4*time.Millisecond, false)
	testGraphSleep(t, 1e3, 3*time.Millisecond, false)
	testGraphSleep(t, 1e3, 2*time.Millisecond, false)
	testGraphSleep(t, 1e3, time.Millisecond, false)
	testGraphSleep(t, 1e3, 100*time.Microsecond, false)

	testGraphSleep(t, 1e2, time.Second, true)
	testGraphSleep(t, 1e2, 100*time.Millisecond, true)
	testGraphSleep(t, 1e3, 10*time.Millisecond, true)
	testGraphSleep(t, 1e3, 5*time.Millisecond, true)
	testGraphSleep(t, 1e3, 4*time.Millisecond, true)
	testGraphSleep(t, 1e3, 3*time.Millisecond, true)
	testGraphSleep(t, 1e3, 2*time.Millisecond, true)
	testGraphSleep(t, 1e3, time.Millisecond, true)
	testGraphSleep(t, 1e3, 100*time.Microsecond, true)
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

	createDelayerQuantitiesGraph(t, "Sleep", "sleep", relativeTimes, duration, stressSystem)
	createDelayerDeviationsGraph(t, "Sleep", "sleep", relativeTimes, duration, stressSystem)
}

func createDelayerQuantitiesGraph(
	t *testing.T,
	titlePerfix string,
	fileNamePerfix string,
	relativeTimes []time.Duration,
	duration time.Duration,
	stressSystem bool,
) {
	quantities, calcInterval := research.CalcIntervalQuantities(relativeTimes, 0, duration)

	axisY, axisX := research.ConvertQuantityOverTimeToBarEcharts(quantities)

	expectedDuration := time.Duration(len(relativeTimes)) * duration

	subtitleAdd := fmt.Sprintf(
		"duration: %s, "+formatTotalDuration(expectedDuration, relativeTimes),
		duration,
	)

	fileNameAdd := fileNamePerfix + "_quantities_duration_" + duration.String()

	createGraph(
		t,
		titlePerfix+" quantities over time",
		subtitleAdd,
		fileNameAdd,
		"quantities",
		len(relativeTimes),
		calcInterval.String(),
		stressSystem,
		axisY,
		axisX,
	)
}

func createDelayerDeviationsGraph(
	t *testing.T,
	titlePerfix string,
	fileNamePerfix string,
	relativeTimes []time.Duration,
	duration time.Duration,
	stressSystem bool,
) {
	deviations := research.CalcRelativeDeviations(relativeTimes, duration)

	axisY, axisX := research.ConvertRelativeDeviationsToBarEcharts(deviations)

	subtitleAdd := fmt.Sprintf(
		"duration: %s",
		duration,
	)

	fileNameAdd := fileNamePerfix + "_deviations_duration_" + duration.String()

	createGraph(
		t,
		titlePerfix+" deviations",
		subtitleAdd,
		fileNameAdd,
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

	subtitleAdd := fmt.Sprintf(
		"limit: {quantity: %d, interval: %s}, "+
			formatTotalDuration(expectedDuration, relativeTimes),
		limit.Quantity,
		limit.Interval,
	)

	fileNameAdd := "quantities_" +
		"limit_quantity_" +
		strconv.Itoa(int(limit.Quantity)) +
		"_limit_interval_" +
		limit.Interval.String()

	createGraph(
		t,
		"Quantities over time",
		subtitleAdd,
		fileNameAdd,
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

	subtitleAdd := fmt.Sprintf(
		"limit: {quantity: %d, interval: %s}, "+
			"flatten: {quantity: %d, interval: %s}",
		limit.Quantity,
		limit.Interval,
		flatten.Quantity,
		flatten.Interval,
	)

	fileNameAdd := "deviations_" +
		"limit_quantity_" +
		strconv.Itoa(int(limit.Quantity)) +
		"_limit_interval_" +
		limit.Interval.String()

	createGraph(
		t,
		"Deviations",
		subtitleAdd,
		fileNameAdd,
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
	subtitleAdd string,
	fileNameAdd string,
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
			subtitleAdd+", "+
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
		fileNameAdd +
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

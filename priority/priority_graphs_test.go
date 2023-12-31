package priority

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/akramarenkov/cqos/internal/consts"
	"github.com/akramarenkov/cqos/priority/internal/common"

	"github.com/go-echarts/go-echarts/v2/charts"
	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
	"github.com/stretchr/testify/require"
)

func addLineSeries(line *charts.Line, serieses map[uint][]chartsopts.LineData) {
	priorities := make([]uint, 0, len(serieses))

	for priority := range serieses {
		priorities = append(priorities, priority)
	}

	common.SortPriorities(priorities)

	for _, priority := range priorities {
		line.AddSeries(strconv.Itoa(int(priority)), serieses[priority])
	}
}

func addBarSeries(bar *charts.Bar, serieses map[uint][]chartsopts.BarData) {
	priorities := make([]uint, 0, len(serieses))

	for priority := range serieses {
		priorities = append(priorities, priority)
	}

	common.SortPriorities(priorities)

	for _, priority := range priorities {
		bar.AddSeries(strconv.Itoa(int(priority)), serieses[priority])
	}
}

func createLineGraph(
	t *testing.T,
	title string,
	subtitle string,
	fileName string,
	serieses map[uint][]chartsopts.LineData,
	abscissa []int,
) {
	if len(serieses) == 0 {
		return
	}

	chart := charts.NewLine()

	chart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    title,
				Subtitle: subtitle,
			},
		),
	)

	addLineSeries(chart.SetXAxis(abscissa), serieses)

	file, err := os.Create(fileName)
	require.NoError(t, err)

	err = chart.Render(file)
	require.NoError(t, err)
}

func createBarGraph(
	t *testing.T,
	title string,
	subtitle string,
	fileName string,
	serieses map[uint][]chartsopts.BarData,
	abscissa []int,
) {
	if len(serieses) == 0 {
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

	addBarSeries(chart.SetXAxis(abscissa), serieses)

	file, err := os.Create(fileName)
	require.NoError(t, err)

	err = chart.Render(file)
	require.NoError(t, err)
}

func createGraphs(
	t *testing.T,
	subtitleBase string,
	filePrefix string,
	handlersQuantity uint,
	unbufferedInput bool,
	measures []measure,
	overTimeResolution time.Duration,
	overTimeUnit time.Duration,
	writeToFeedbackInterval time.Duration,
) {
	received := filterByKind(measures, measureKindReceived)

	dqot, dqotX := convertToLineEcharts(
		calcDataQuantity(received, overTimeResolution),
		overTimeUnit,
	)

	ipot, ipotX := convertToLineEcharts(
		calcInProcessing(measures, overTimeResolution),
		overTimeUnit,
	)

	wtfl, wtflX := convertToBarEcharts(
		calcWriteToFeedbackLatency(measures, writeToFeedbackInterval),
	)

	subtitle := fmt.Sprintf(
		subtitleBase+
			", "+
			"handlers quantity: %d, "+
			"buffered: %t, "+
			"time: %s",
		handlersQuantity,
		!unbufferedInput,
		time.Now().Format(time.RFC3339),
	)

	baseName := "graph_" + filePrefix + "_" +
		strconv.Itoa(int(handlersQuantity)) +
		"_buffered_" +
		strconv.FormatBool(!unbufferedInput)

	createLineGraph(
		t,
		"Data retrieval graph",
		subtitle,
		baseName+"_data_retrieval.html",
		dqot,
		dqotX,
	)

	createLineGraph(
		t,
		"In processing graph",
		subtitle,
		baseName+"_in_processing.html",
		ipot,
		ipotX,
	)

	createBarGraph(
		t,
		"Write to feedback latency",
		subtitle,
		baseName+"_write_feedback_latency.html",
		wtfl,
		wtflX,
	)
}

func testGraphFairEven(t *testing.T, factor uint, unbufferedInput bool) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	measurerOpts := measurerOpts{
		HandlersQuantity: 6 * factor,
		UnbufferedInput:  unbufferedInput,
	}

	msr := newMeasurer(measurerOpts)

	msr.AddWrite(1, 4000*factor)

	msr.AddWrite(2, 500*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 500*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 4*time.Second)
	msr.AddWrite(2, 1000*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 2000*factor)

	msr.AddWrite(3, 500*factor)
	msr.AddWaitDevastation(3)
	msr.AddDelay(3, 5*time.Second)
	msr.AddWrite(3, 4000*factor)

	msr.SetProcessDelay(1, 10*time.Millisecond)
	msr.SetProcessDelay(2, 10*time.Millisecond)
	msr.SetProcessDelay(3, 10*time.Millisecond)

	opts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         msr.GetFeedback(),
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
		Output:           msr.GetOutput(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	defer discipline.Stop()

	measures := msr.Play(discipline, false)

	createGraphs(
		t,
		"Fair divider, even time processing",
		"fair_even",
		measurerOpts.HandlersQuantity,
		unbufferedInput,
		measures,
		100*time.Millisecond,
		1*time.Second,
		100*time.Nanosecond,
	)
}

func TestGraphFairEven(t *testing.T) {
	testGraphFairEven(t, 1, true)
	testGraphFairEven(t, 10, true)
	testGraphFairEven(t, 1, false)
	testGraphFairEven(t, 10, false)
}

func testGraphFairUneven(t *testing.T, factor uint, unbufferedInput bool) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	measurerOpts := measurerOpts{
		HandlersQuantity: 6 * factor,
		UnbufferedInput:  unbufferedInput,
	}

	msr := newMeasurer(measurerOpts)

	msr.AddWrite(1, 450*factor)

	msr.AddWrite(2, 100*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 100*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 4*time.Second)
	msr.AddWrite(2, 200*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 400*factor)

	msr.AddWrite(3, 500*factor)
	msr.AddWaitDevastation(3)
	msr.AddDelay(3, 6*time.Second)
	msr.AddWrite(3, 3000*factor)

	msr.SetProcessDelay(1, 100*time.Millisecond)
	msr.SetProcessDelay(2, 50*time.Millisecond)
	msr.SetProcessDelay(3, 10*time.Millisecond)

	opts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         msr.GetFeedback(),
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
		Output:           msr.GetOutput(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	defer discipline.Stop()

	measures := msr.Play(discipline, false)

	createGraphs(
		t,
		"Fair divider, uneven time processing",
		"fair_uneven",
		measurerOpts.HandlersQuantity,
		unbufferedInput,
		measures,
		100*time.Millisecond,
		1*time.Second,
		100*time.Nanosecond,
	)
}

func TestGraphFairUneven(t *testing.T) {
	testGraphFairUneven(t, 1, true)
	testGraphFairUneven(t, 10, true)
	testGraphFairUneven(t, 1, false)
	testGraphFairUneven(t, 10, false)
}

func testGraphRateEven(t *testing.T, factor uint, unbufferedInput bool) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	measurerOpts := measurerOpts{
		HandlersQuantity: 6 * factor,
		UnbufferedInput:  unbufferedInput,
	}

	msr := newMeasurer(measurerOpts)

	msr.AddWrite(1, 4100*factor)

	msr.AddWrite(2, 1500*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 750*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 4*time.Second)
	msr.AddWrite(2, 700*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 3*time.Second)
	msr.AddWrite(2, 1200*factor)

	msr.AddWrite(3, 1000*factor)
	msr.AddWaitDevastation(3)
	msr.AddDelay(3, 8*time.Second)
	msr.AddWrite(3, 3700*factor)

	msr.SetProcessDelay(1, 10*time.Millisecond)
	msr.SetProcessDelay(2, 10*time.Millisecond)
	msr.SetProcessDelay(3, 10*time.Millisecond)

	opts := Opts[uint]{
		Divider:          RateDivider,
		Feedback:         msr.GetFeedback(),
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
		Output:           msr.GetOutput(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	defer discipline.Stop()

	measures := msr.Play(discipline, false)

	createGraphs(
		t,
		"Rate divider, even time processing",
		"rate_even",
		measurerOpts.HandlersQuantity,
		unbufferedInput,
		measures,
		100*time.Millisecond,
		1*time.Second,
		100*time.Nanosecond,
	)
}

func TestGraphRateEven(t *testing.T) {
	testGraphRateEven(t, 1, true)
	testGraphRateEven(t, 10, true)
	testGraphRateEven(t, 1, false)
	testGraphRateEven(t, 10, false)
}

func testGraphRateUneven(t *testing.T, factor uint, unbufferedInput bool) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	measurerOpts := measurerOpts{
		HandlersQuantity: 6 * factor,
		UnbufferedInput:  unbufferedInput,
	}

	msr := newMeasurer(measurerOpts)

	msr.AddWrite(1, 430*factor)

	msr.AddWrite(2, 250*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 100*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 4*time.Second)
	msr.AddWrite(2, 150*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 300*factor)

	msr.AddWrite(3, 1000*factor)
	msr.AddWaitDevastation(3)
	msr.AddDelay(3, 8*time.Second)
	msr.AddWrite(3, 3500*factor)

	msr.SetProcessDelay(1, 100*time.Millisecond)
	msr.SetProcessDelay(2, 50*time.Millisecond)
	msr.SetProcessDelay(3, 10*time.Millisecond)

	opts := Opts[uint]{
		Divider:          RateDivider,
		Feedback:         msr.GetFeedback(),
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
		Output:           msr.GetOutput(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	defer discipline.Stop()

	measures := msr.Play(discipline, false)

	createGraphs(
		t,
		"Rate divider, uneven time processing",
		"rate_uneven",
		measurerOpts.HandlersQuantity,
		unbufferedInput,
		measures,
		100*time.Millisecond,
		1*time.Second,
		100*time.Nanosecond,
	)
}

func TestGraphRateUneven(t *testing.T) {
	testGraphRateUneven(t, 1, true)
	testGraphRateUneven(t, 10, true)
	testGraphRateUneven(t, 1, false)
	testGraphRateUneven(t, 10, false)
}

func testGraphUnmanagedEven(t *testing.T, factor uint, unbufferedInput bool) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	measurerOpts := measurerOpts{
		HandlersQuantity: 6 * factor,
		NoFeedback:       true,
		UnbufferedInput:  unbufferedInput,
	}

	msr := newMeasurer(measurerOpts)

	msr.AddWrite(1, 4000*factor)

	msr.AddWrite(2, 500*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 500*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 4*time.Second)
	msr.AddWrite(2, 1000*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 2000*factor)

	msr.AddWrite(3, 500*factor)
	msr.AddWaitDevastation(3)
	msr.AddDelay(3, 5*time.Second)
	msr.AddWrite(3, 4000*factor)

	msr.SetProcessDelay(1, 10*time.Millisecond)
	msr.SetProcessDelay(2, 10*time.Millisecond)
	msr.SetProcessDelay(3, 10*time.Millisecond)

	unmanagedOpts := unmanagedOpts[uint]{
		Inputs: msr.GetInputs(),
		Output: msr.GetOutput(),
	}

	unmanaged, err := newUnmanaged(unmanagedOpts)
	require.NoError(t, err)

	defer unmanaged.Stop()

	measures := msr.Play(unmanaged, false)

	createGraphs(
		t,
		"Unmanaged, even time processing",
		"unmanaged_even",
		measurerOpts.HandlersQuantity,
		unbufferedInput,
		measures,
		100*time.Millisecond,
		1*time.Second,
		100*time.Nanosecond,
	)
}

func TestGraphUnmanagedEven(t *testing.T) {
	testGraphUnmanagedEven(t, 1, true)
	testGraphUnmanagedEven(t, 10, true)
	testGraphUnmanagedEven(t, 1, false)
	testGraphUnmanagedEven(t, 10, false)
}

func testGraphUnmanagedUneven(t *testing.T, factor uint, unbufferedInput bool) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	measurerOpts := measurerOpts{
		HandlersQuantity: 6 * factor,
		NoFeedback:       true,
		UnbufferedInput:  unbufferedInput,
	}

	msr := newMeasurer(measurerOpts)

	msr.AddWrite(1, 500*factor)

	msr.AddWrite(2, 100*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 100*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 200*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 400*factor)

	msr.AddWrite(3, 100*factor)
	msr.AddWaitDevastation(3)
	msr.AddDelay(3, 6*time.Second)
	msr.AddWrite(3, 1350*factor)

	msr.SetProcessDelay(1, 100*time.Millisecond)
	msr.SetProcessDelay(2, 50*time.Millisecond)
	msr.SetProcessDelay(3, 10*time.Millisecond)

	unmanagedOpts := unmanagedOpts[uint]{
		Inputs: msr.GetInputs(),
		Output: msr.GetOutput(),
	}

	unmanaged, err := newUnmanaged(unmanagedOpts)
	require.NoError(t, err)

	defer unmanaged.Stop()

	measures := msr.Play(unmanaged, false)

	createGraphs(
		t,
		"Unmanaged, uneven time processing",
		"unmanaged_uneven",
		measurerOpts.HandlersQuantity,
		unbufferedInput,
		measures,
		100*time.Millisecond,
		1*time.Second,
		100*time.Nanosecond,
	)
}

func TestGraphUnmanagedUneven(t *testing.T) {
	testGraphUnmanagedUneven(t, 1, true)
	testGraphUnmanagedUneven(t, 10, true)
	testGraphUnmanagedUneven(t, 1, false)
	testGraphUnmanagedUneven(t, 10, false)
}

func testGraphFairEvenDividingError(t *testing.T, handlersQuantity uint) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	measurerOpts := measurerOpts{
		HandlersQuantity: handlersQuantity,
	}

	msr := newMeasurer(measurerOpts)

	msr.AddWrite(1, 4000)

	msr.AddWrite(2, 500)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 500)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 4*time.Second)
	msr.AddWrite(2, 1000)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 2000)

	msr.AddWrite(3, 500)
	msr.AddWaitDevastation(3)
	msr.AddDelay(3, 5*time.Second)
	msr.AddWrite(3, 4000)

	msr.AddWrite(4, 500)
	msr.AddWaitDevastation(3)
	msr.AddDelay(4, 5*time.Second)
	msr.AddWrite(4, 4000)

	msr.SetProcessDelay(1, 10*time.Millisecond)
	msr.SetProcessDelay(2, 10*time.Millisecond)
	msr.SetProcessDelay(3, 10*time.Millisecond)
	msr.SetProcessDelay(4, 10*time.Millisecond)

	opts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         msr.GetFeedback(),
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
		Output:           msr.GetOutput(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	defer discipline.Stop()

	measures := msr.Play(discipline, false)

	createGraphs(
		t,
		"Fair divider, even time processing, significant dividing error",
		"fair_even_dividing_error",
		measurerOpts.HandlersQuantity,
		measurerOpts.UnbufferedInput,
		measures,
		100*time.Millisecond,
		1*time.Second,
		100*time.Nanosecond,
	)
}

func TestGraphFairEvenDividingError(t *testing.T) {
	testGraphFairEvenDividingError(t, 6)
	testGraphFairEvenDividingError(t, 7)
	testGraphFairEvenDividingError(t, 8)
	testGraphFairEvenDividingError(t, 9)
	testGraphFairEvenDividingError(t, 10)
	testGraphFairEvenDividingError(t, 11)
	testGraphFairEvenDividingError(t, 12)
}

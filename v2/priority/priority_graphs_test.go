package priority

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/consts"
	"github.com/akramarenkov/cqos/v2/priority/divider"
	"github.com/akramarenkov/cqos/v2/priority/internal/common"
	"github.com/akramarenkov/cqos/v2/priority/internal/measurer"
	"github.com/akramarenkov/cqos/v2/priority/internal/research"
	"github.com/akramarenkov/cqos/v2/priority/internal/unmanaged"

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
	abscissa []uint,
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
	abscissa []uint,
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
	dqot map[uint][]chartsopts.LineData,
	dqotX []uint,
	ipot map[uint][]chartsopts.LineData,
	ipotX []uint,
	wtfl map[uint][]chartsopts.BarData,
	wtflX []uint,
) {
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

func testDisciplineFairEvenProcessingTime(
	t *testing.T,
	factor uint,
	unbufferedInput bool,
) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	measurerOpts := measurer.Opts{
		HandlersQuantity: 6 * factor,
		UnbufferedInput:  unbufferedInput,
	}

	msr := measurer.New(measurerOpts)

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
		Divider:          divider.Fair,
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measures := msr.Play(discipline)

	received := research.FilterByKind(measures, measurer.MeasureKindReceived)

	dqot, dqotX := research.ConvertToLineEcharts(
		research.CalcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := research.ConvertToLineEcharts(
		research.CalcInProcessing(measures, 100*time.Millisecond),
		1*time.Second,
	)

	wtfl, wtflX := research.ConvertToBarEcharts(
		research.CalcWriteToFeedbackLatency(measures, 100*time.Nanosecond),
	)

	createGraphs(
		t,
		"Fair divider, even time processing",
		"fair_even",
		measurerOpts.HandlersQuantity,
		unbufferedInput,
		dqot,
		dqotX,
		ipot,
		ipotX,
		wtfl,
		wtflX,
	)
}

func TestDisciplineFairEvenProcessingTime(t *testing.T) {
	testDisciplineFairEvenProcessingTime(t, 1, true)
	testDisciplineFairEvenProcessingTime(t, 10, true)
	testDisciplineFairEvenProcessingTime(t, 1, false)
	testDisciplineFairEvenProcessingTime(t, 10, false)
}

func testDisciplineFairUnevenProcessingTime(
	t *testing.T,
	factor uint,
	unbufferedInput bool,
) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	measurerOpts := measurer.Opts{
		HandlersQuantity: 6 * factor,
		UnbufferedInput:  unbufferedInput,
	}

	msr := measurer.New(measurerOpts)

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
		Divider:          divider.Fair,
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measures := msr.Play(discipline)

	received := research.FilterByKind(measures, measurer.MeasureKindReceived)

	dqot, dqotX := research.ConvertToLineEcharts(
		research.CalcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := research.ConvertToLineEcharts(
		research.CalcInProcessing(measures, 100*time.Millisecond),
		1*time.Second,
	)

	wtfl, wtflX := research.ConvertToBarEcharts(
		research.CalcWriteToFeedbackLatency(measures, 100*time.Nanosecond),
	)

	createGraphs(
		t,
		"Fair divider, uneven time processing",
		"fair_uneven",
		measurerOpts.HandlersQuantity,
		unbufferedInput,
		dqot,
		dqotX,
		ipot,
		ipotX,
		wtfl,
		wtflX,
	)
}

func TestDisciplineFairUnevenProcessingTime(t *testing.T) {
	testDisciplineFairUnevenProcessingTime(t, 1, true)
	testDisciplineFairUnevenProcessingTime(t, 10, true)
	testDisciplineFairUnevenProcessingTime(t, 1, false)
	testDisciplineFairUnevenProcessingTime(t, 10, false)
}

func testDisciplineRateEvenProcessingTime(
	t *testing.T,
	factor uint,
	unbufferedInput bool,
) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	measurerOpts := measurer.Opts{
		HandlersQuantity: 6 * factor,
		UnbufferedInput:  unbufferedInput,
	}

	msr := measurer.New(measurerOpts)

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
		Divider:          divider.Rate,
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measures := msr.Play(discipline)

	received := research.FilterByKind(measures, measurer.MeasureKindReceived)

	dqot, dqotX := research.ConvertToLineEcharts(
		research.CalcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := research.ConvertToLineEcharts(
		research.CalcInProcessing(measures, 100*time.Millisecond),
		1*time.Second,
	)

	wtfl, wtflX := research.ConvertToBarEcharts(
		research.CalcWriteToFeedbackLatency(measures, 100*time.Nanosecond),
	)

	createGraphs(
		t,
		"Rate divider, even time processing",
		"rate_even",
		measurerOpts.HandlersQuantity,
		unbufferedInput,
		dqot,
		dqotX,
		ipot,
		ipotX,
		wtfl,
		wtflX,
	)
}

func TestDisciplineRateEvenProcessingTime(t *testing.T) {
	testDisciplineRateEvenProcessingTime(t, 1, true)
	testDisciplineRateEvenProcessingTime(t, 10, true)
	testDisciplineRateEvenProcessingTime(t, 1, false)
	testDisciplineRateEvenProcessingTime(t, 10, false)
}

func testDisciplineRateUnevenProcessingTime(
	t *testing.T,
	factor uint,
	unbufferedInput bool,
) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	measurerOpts := measurer.Opts{
		HandlersQuantity: 6 * factor,
		UnbufferedInput:  unbufferedInput,
	}

	msr := measurer.New(measurerOpts)

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
		Divider:          divider.Rate,
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measures := msr.Play(discipline)

	received := research.FilterByKind(measures, measurer.MeasureKindReceived)

	dqot, dqotX := research.ConvertToLineEcharts(
		research.CalcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := research.ConvertToLineEcharts(
		research.CalcInProcessing(measures, 100*time.Millisecond),
		1*time.Second,
	)

	wtfl, wtflX := research.ConvertToBarEcharts(
		research.CalcWriteToFeedbackLatency(measures, 100*time.Nanosecond),
	)

	createGraphs(
		t,
		"Rate divider, uneven time processing",
		"rate_uneven",
		measurerOpts.HandlersQuantity,
		unbufferedInput,
		dqot,
		dqotX,
		ipot,
		ipotX,
		wtfl,
		wtflX,
	)
}

func TestDisciplineRateUnevenProcessingTime(t *testing.T) {
	testDisciplineRateUnevenProcessingTime(t, 1, true)
	testDisciplineRateUnevenProcessingTime(t, 10, true)
	testDisciplineRateUnevenProcessingTime(t, 1, false)
	testDisciplineRateUnevenProcessingTime(t, 10, false)
}

func testUnmanagedEven(t *testing.T, factor uint, unbufferedInput bool) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	measurerOpts := measurer.Opts{
		HandlersQuantity: 6 * factor,
		UnbufferedInput:  unbufferedInput,
	}

	msr := measurer.New(measurerOpts)

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

	unmanagedOpts := unmanaged.Opts[uint]{
		Inputs: msr.GetInputs(),
	}

	unmanaged, err := unmanaged.New(unmanagedOpts)
	require.NoError(t, err)

	measures := msr.Play(unmanaged)

	received := research.FilterByKind(measures, measurer.MeasureKindReceived)

	dqot, dqotX := research.ConvertToLineEcharts(
		research.CalcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := research.ConvertToLineEcharts(
		research.CalcInProcessing(measures, 100*time.Millisecond),
		1*time.Second,
	)

	createGraphs(
		t,
		"Unmanaged, even time processing",
		"unmanaged_even",
		measurerOpts.HandlersQuantity,
		unbufferedInput,
		dqot,
		dqotX,
		ipot,
		ipotX,
		nil,
		nil,
	)
}

func TestUnmanagedEven(t *testing.T) {
	testUnmanagedEven(t, 1, true)
	testUnmanagedEven(t, 10, true)
	testUnmanagedEven(t, 1, false)
	testUnmanagedEven(t, 10, false)
}

func testUnmanagedUneven(t *testing.T, factor uint, unbufferedInput bool) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	measurerOpts := measurer.Opts{
		HandlersQuantity: 6 * factor,
		UnbufferedInput:  unbufferedInput,
	}

	msr := measurer.New(measurerOpts)

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

	unmanagedOpts := unmanaged.Opts[uint]{
		Inputs: msr.GetInputs(),
	}

	unmanaged, err := unmanaged.New(unmanagedOpts)
	require.NoError(t, err)

	measures := msr.Play(unmanaged)

	received := research.FilterByKind(measures, measurer.MeasureKindReceived)

	dqot, dqotX := research.ConvertToLineEcharts(
		research.CalcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := research.ConvertToLineEcharts(
		research.CalcInProcessing(measures, 100*time.Millisecond),
		1*time.Second,
	)

	createGraphs(
		t,
		"Unmanaged, uneven time processing",
		"unmanaged_uneven",
		measurerOpts.HandlersQuantity,
		unbufferedInput,
		dqot,
		dqotX,
		ipot,
		ipotX,
		nil,
		nil,
	)
}

func TestUnmanagedUneven(t *testing.T) {
	testUnmanagedUneven(t, 1, true)
	testUnmanagedUneven(t, 10, true)
	testUnmanagedUneven(t, 1, false)
	testUnmanagedUneven(t, 10, false)
}

func testDisciplineFairEvenProcessingTimeDividingError(t *testing.T, handlersQuantity uint) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	measurerOpts := measurer.Opts{
		HandlersQuantity: handlersQuantity,
	}

	msr := measurer.New(measurerOpts)

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
		Divider:          divider.Fair,
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measures := msr.Play(discipline)

	received := research.FilterByKind(measures, measurer.MeasureKindReceived)

	dqot, dqotX := research.ConvertToLineEcharts(
		research.CalcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	createGraphs(
		t,
		"Fair divider, even time processing, significant dividing error",
		"fair_even_dividing_error",
		measurerOpts.HandlersQuantity,
		measurerOpts.UnbufferedInput,
		dqot,
		dqotX,
		nil,
		nil,
		nil,
		nil,
	)
}

func TestDisciplineFairEvenProcessingTimeDividingError(t *testing.T) {
	testDisciplineFairEvenProcessingTimeDividingError(t, 6)
	testDisciplineFairEvenProcessingTimeDividingError(t, 7)
	testDisciplineFairEvenProcessingTimeDividingError(t, 8)
	testDisciplineFairEvenProcessingTimeDividingError(t, 9)
	testDisciplineFairEvenProcessingTimeDividingError(t, 10)
	testDisciplineFairEvenProcessingTimeDividingError(t, 11)
	testDisciplineFairEvenProcessingTimeDividingError(t, 12)
}

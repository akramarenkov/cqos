package priority

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/akramarenkov/cqos/internal/consts"

	"github.com/go-echarts/go-echarts/v2/charts"
	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
	"github.com/stretchr/testify/require"
)

func testDisciplineRateEvenProcessingTime(t *testing.T, factor uint, inputBuffered bool) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	handlersQuantity := uint(6) * factor

	measurerOpts := measurerOpts{
		HandlersQuantity: handlersQuantity,
		UnbufferedInput:  !inputBuffered,
	}

	measurer := newMeasurer(measurerOpts)
	defer measurer.Finalize()

	measurer.AddWrite(1, 4100*factor)

	measurer.AddWrite(2, 1500*factor)
	measurer.AddWaitDevastation(2)
	measurer.AddDelay(2, 2*time.Second)
	measurer.AddWrite(2, 750*factor)
	measurer.AddWaitDevastation(2)
	measurer.AddDelay(2, 4*time.Second)
	measurer.AddWrite(2, 700*factor)
	measurer.AddWaitDevastation(2)
	measurer.AddDelay(2, 3*time.Second)
	measurer.AddWrite(2, 1200*factor)

	measurer.AddWrite(3, 1000*factor)
	measurer.AddWaitDevastation(3)
	measurer.AddDelay(3, 8*time.Second)
	measurer.AddWrite(3, 3700*factor)

	measurer.SetProcessDelay(1, 10*time.Millisecond)
	measurer.SetProcessDelay(2, 10*time.Millisecond)
	measurer.SetProcessDelay(3, 10*time.Millisecond)

	disciplineOpts := Opts[uint]{
		Divider:          RateDivider,
		Feedback:         measurer.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           measurer.GetInputs(),
		Output:           measurer.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	measures := measurer.Play(discipline)

	received := filterByKind(measures, measureKindReceived)

	dqot, dqotX := convertToLineEcharts(
		calcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := convertToLineEcharts(
		calcInProcessing(measures, 100*time.Millisecond),
		1*time.Second,
	)

	wtfl, wtflX := convertToBarEcharts(
		calcWriteToFeedbackLatency(measures, 100*time.Nanosecond),
	)

	dqotChart := charts.NewLine()
	ipotChart := charts.NewLine()
	wtflChart := charts.NewBar()

	subtitle := fmt.Sprintf(
		"Rate divider, even time processing, "+
			"handlers quantity: %d, buffered: %t, time: %s",
		handlersQuantity,
		inputBuffered,
		time.Now().Format(time.RFC3339),
	)

	dqotChart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "Data retrieval graph",
				Subtitle: subtitle,
			},
		),
	)

	ipotChart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "In processing graph",
				Subtitle: subtitle,
			},
		),
	)

	wtflChart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "Write to feedback latency",
				Subtitle: subtitle,
			},
		),
	)

	dqotChart.SetXAxis(dqotX).
		AddSeries("3", dqot[3]).
		AddSeries("2", dqot[2]).
		AddSeries("1", dqot[1])

	ipotChart.SetXAxis(ipotX).
		AddSeries("3", ipot[3]).
		AddSeries("2", ipot[2]).
		AddSeries("1", ipot[1])

	wtflChart.SetXAxis(wtflX).
		AddSeries("3", wtfl[3]).
		AddSeries("2", wtfl[2]).
		AddSeries("1", wtfl[1])

	baseName := "graph_rate_even_" + strconv.Itoa(int(handlersQuantity)) +
		"_buffered_" + strconv.FormatBool(inputBuffered)

	dqotFile, err := os.Create(baseName + "_data_retrieval.html")
	require.NoError(t, err)

	err = dqotChart.Render(dqotFile)
	require.NoError(t, err)

	ipotFile, err := os.Create(baseName + "_in_processing.html")
	require.NoError(t, err)

	err = ipotChart.Render(ipotFile)
	require.NoError(t, err)

	wtflFile, err := os.Create(baseName + "_write_feedback_latency.html")
	require.NoError(t, err)

	err = wtflChart.Render(wtflFile)
	require.NoError(t, err)
}

func TestDisciplineRateEvenProcessingTime(t *testing.T) {
	testDisciplineRateEvenProcessingTime(t, 1, true)
	testDisciplineRateEvenProcessingTime(t, 10, true)
	testDisciplineRateEvenProcessingTime(t, 1, false)
	testDisciplineRateEvenProcessingTime(t, 10, false)
}

func testDisciplineRateUnevenProcessingTime(t *testing.T, factor uint, inputBuffered bool) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	handlersQuantity := uint(6) * factor

	measurerOpts := measurerOpts{
		HandlersQuantity: handlersQuantity,
		UnbufferedInput:  !inputBuffered,
	}

	measurer := newMeasurer(measurerOpts)
	defer measurer.Finalize()

	measurer.AddWrite(1, 430*factor)

	measurer.AddWrite(2, 250*factor)
	measurer.AddWaitDevastation(2)
	measurer.AddDelay(2, 2*time.Second)
	measurer.AddWrite(2, 100*factor)
	measurer.AddWaitDevastation(2)
	measurer.AddDelay(2, 4*time.Second)
	measurer.AddWrite(2, 150*factor)
	measurer.AddWaitDevastation(2)
	measurer.AddDelay(2, 2*time.Second)
	measurer.AddWrite(2, 300*factor)

	measurer.AddWrite(3, 1000*factor)
	measurer.AddWaitDevastation(3)
	measurer.AddDelay(3, 8*time.Second)
	measurer.AddWrite(3, 3500*factor)

	measurer.SetProcessDelay(1, 100*time.Millisecond)
	measurer.SetProcessDelay(2, 50*time.Millisecond)
	measurer.SetProcessDelay(3, 10*time.Millisecond)

	disciplineOpts := Opts[uint]{
		Divider:          RateDivider,
		Feedback:         measurer.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           measurer.GetInputs(),
		Output:           measurer.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	measures := measurer.Play(discipline)

	received := filterByKind(measures, measureKindReceived)

	dqot, dqotX := convertToLineEcharts(
		calcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := convertToLineEcharts(
		calcInProcessing(measures, 100*time.Millisecond),
		1*time.Second,
	)

	wtfl, wtflX := convertToBarEcharts(
		calcWriteToFeedbackLatency(measures, 100*time.Nanosecond),
	)

	dqotChart := charts.NewLine()
	ipotChart := charts.NewLine()
	wtflChart := charts.NewBar()

	subtitle := fmt.Sprintf(
		"Rate divider, uneven time processing, "+
			"handlers quantity: %d, buffered: %t, time: %s",
		handlersQuantity,
		inputBuffered,
		time.Now().Format(time.RFC3339),
	)

	dqotChart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "Data retrieval graph",
				Subtitle: subtitle,
			},
		),
	)

	ipotChart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "In processing graph",
				Subtitle: subtitle,
			},
		),
	)

	wtflChart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "Write to feedback latency",
				Subtitle: subtitle,
			},
		),
	)

	dqotChart.SetXAxis(dqotX).
		AddSeries("3", dqot[3]).
		AddSeries("2", dqot[2]).
		AddSeries("1", dqot[1])

	ipotChart.SetXAxis(ipotX).
		AddSeries("3", ipot[3]).
		AddSeries("2", ipot[2]).
		AddSeries("1", ipot[1])

	wtflChart.SetXAxis(wtflX).
		AddSeries("3", wtfl[3]).
		AddSeries("2", wtfl[2]).
		AddSeries("1", wtfl[1])

	baseName := "graph_rate_uneven_" + strconv.Itoa(int(handlersQuantity)) +
		"_buffered_" + strconv.FormatBool(inputBuffered)

	dqotFile, err := os.Create(baseName + "_data_retrieval.html")
	require.NoError(t, err)

	err = dqotChart.Render(dqotFile)
	require.NoError(t, err)

	ipotFile, err := os.Create(baseName + "_in_processing.html")
	require.NoError(t, err)

	err = ipotChart.Render(ipotFile)
	require.NoError(t, err)

	wtflFile, err := os.Create(baseName + "_write_feedback_latency.html")
	require.NoError(t, err)

	err = wtflChart.Render(wtflFile)
	require.NoError(t, err)
}

func TestDisciplineRateUnevenProcessingTime(t *testing.T) {
	testDisciplineRateUnevenProcessingTime(t, 1, true)
	testDisciplineRateUnevenProcessingTime(t, 10, true)
	testDisciplineRateUnevenProcessingTime(t, 1, false)
	testDisciplineRateUnevenProcessingTime(t, 10, false)
}

func testDisciplineFairEvenProcessingTime(t *testing.T, factor uint, inputBuffered bool) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	handlersQuantity := uint(6) * factor

	measurerOpts := measurerOpts{
		HandlersQuantity: handlersQuantity,
		UnbufferedInput:  !inputBuffered,
	}

	measurer := newMeasurer(measurerOpts)
	defer measurer.Finalize()

	measurer.AddWrite(1, 4000*factor)

	measurer.AddWrite(2, 500*factor)
	measurer.AddWaitDevastation(2)
	measurer.AddDelay(2, 2*time.Second)
	measurer.AddWrite(2, 500*factor)
	measurer.AddWaitDevastation(2)
	measurer.AddDelay(2, 4*time.Second)
	measurer.AddWrite(2, 1000*factor)
	measurer.AddWaitDevastation(2)
	measurer.AddDelay(2, 2*time.Second)
	measurer.AddWrite(2, 2000*factor)

	measurer.AddWrite(3, 500*factor)
	measurer.AddWaitDevastation(3)
	measurer.AddDelay(3, 5*time.Second)
	measurer.AddWrite(3, 4000*factor)

	measurer.SetProcessDelay(1, 10*time.Millisecond)
	measurer.SetProcessDelay(2, 10*time.Millisecond)
	measurer.SetProcessDelay(3, 10*time.Millisecond)

	disciplineOpts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         measurer.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           measurer.GetInputs(),
		Output:           measurer.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	measures := measurer.Play(discipline)

	received := filterByKind(measures, measureKindReceived)

	dqot, dqotX := convertToLineEcharts(
		calcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := convertToLineEcharts(
		calcInProcessing(measures, 100*time.Millisecond),
		1*time.Second,
	)

	wtfl, wtflX := convertToBarEcharts(
		calcWriteToFeedbackLatency(measures, 100*time.Nanosecond),
	)

	dqotChart := charts.NewLine()
	ipotChart := charts.NewLine()
	wtflChart := charts.NewBar()

	subtitle := fmt.Sprintf(
		"Fair divider, even time processing, "+
			"handlers quantity: %d, buffered: %t, time: %s",
		handlersQuantity,
		inputBuffered,
		time.Now().Format(time.RFC3339),
	)

	dqotChart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "Data retrieval graph",
				Subtitle: subtitle,
			},
		),
	)

	ipotChart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "In processing graph",
				Subtitle: subtitle,
			},
		),
	)

	wtflChart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "Write to feedback latency",
				Subtitle: subtitle,
			},
		),
	)

	dqotChart.SetXAxis(dqotX).
		AddSeries("3", dqot[3]).
		AddSeries("2", dqot[2]).
		AddSeries("1", dqot[1])

	ipotChart.SetXAxis(ipotX).
		AddSeries("3", ipot[3]).
		AddSeries("2", ipot[2]).
		AddSeries("1", ipot[1])

	wtflChart.SetXAxis(wtflX).
		AddSeries("3", wtfl[3]).
		AddSeries("2", wtfl[2]).
		AddSeries("1", wtfl[1])

	baseName := "graph_fair_even_" + strconv.Itoa(int(handlersQuantity)) +
		"_buffered_" + strconv.FormatBool(inputBuffered)

	dqotFile, err := os.Create(baseName + "_data_retrieval.html")
	require.NoError(t, err)

	err = dqotChart.Render(dqotFile)
	require.NoError(t, err)

	ipotFile, err := os.Create(baseName + "_in_processing.html")
	require.NoError(t, err)

	err = ipotChart.Render(ipotFile)
	require.NoError(t, err)

	wtflFile, err := os.Create(baseName + "_write_feedback_latency.html")
	require.NoError(t, err)

	err = wtflChart.Render(wtflFile)
	require.NoError(t, err)
}

func TestDisciplineFairEvenProcessingTime(t *testing.T) {
	testDisciplineFairEvenProcessingTime(t, 1, true)
	testDisciplineFairEvenProcessingTime(t, 10, true)
	testDisciplineFairEvenProcessingTime(t, 1, false)
	testDisciplineFairEvenProcessingTime(t, 10, false)
}

func testDisciplineFairUnevenProcessingTime(t *testing.T, factor uint, inputBuffered bool) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	handlersQuantity := uint(6) * factor

	measurerOpts := measurerOpts{
		HandlersQuantity: handlersQuantity,
		UnbufferedInput:  !inputBuffered,
	}

	measurer := newMeasurer(measurerOpts)
	defer measurer.Finalize()

	measurer.AddWrite(1, 450*factor)

	measurer.AddWrite(2, 100*factor)
	measurer.AddWaitDevastation(2)
	measurer.AddDelay(2, 2*time.Second)
	measurer.AddWrite(2, 100*factor)
	measurer.AddWaitDevastation(2)
	measurer.AddDelay(2, 4*time.Second)
	measurer.AddWrite(2, 200*factor)
	measurer.AddWaitDevastation(2)
	measurer.AddDelay(2, 2*time.Second)
	measurer.AddWrite(2, 400*factor)

	measurer.AddWrite(3, 500*factor)
	measurer.AddWaitDevastation(3)
	measurer.AddDelay(3, 6*time.Second)
	measurer.AddWrite(3, 3000*factor)

	measurer.SetProcessDelay(1, 100*time.Millisecond)
	measurer.SetProcessDelay(2, 50*time.Millisecond)
	measurer.SetProcessDelay(3, 10*time.Millisecond)

	disciplineOpts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         measurer.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           measurer.GetInputs(),
		Output:           measurer.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	measures := measurer.Play(discipline)

	received := filterByKind(measures, measureKindReceived)

	dqot, dqotX := convertToLineEcharts(
		calcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := convertToLineEcharts(
		calcInProcessing(measures, 100*time.Millisecond),
		1*time.Second,
	)

	wtfl, wtflX := convertToBarEcharts(
		calcWriteToFeedbackLatency(measures, 100*time.Nanosecond),
	)

	dqotChart := charts.NewLine()
	ipotChart := charts.NewLine()
	wtflChart := charts.NewBar()

	subtitle := fmt.Sprintf(
		"Fair divider, uneven time processing, "+
			"handlers quantity: %d, buffered: %t, time: %s",
		handlersQuantity,
		inputBuffered,
		time.Now().Format(time.RFC3339),
	)

	dqotChart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "Data retrieval graph",
				Subtitle: subtitle,
			},
		),
	)

	ipotChart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "In processing graph",
				Subtitle: subtitle,
			},
		),
	)

	wtflChart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "Write to feedback latency",
				Subtitle: subtitle,
			},
		),
	)

	dqotChart.SetXAxis(dqotX).
		AddSeries("3", dqot[3]).
		AddSeries("2", dqot[2]).
		AddSeries("1", dqot[1])

	ipotChart.SetXAxis(ipotX).
		AddSeries("3", ipot[3]).
		AddSeries("2", ipot[2]).
		AddSeries("1", ipot[1])

	wtflChart.SetXAxis(wtflX).
		AddSeries("3", wtfl[3]).
		AddSeries("2", wtfl[2]).
		AddSeries("1", wtfl[1])

	baseName := "graph_fair_uneven_" + strconv.Itoa(int(handlersQuantity)) +
		"_buffered_" + strconv.FormatBool(inputBuffered)

	dqotFile, err := os.Create(baseName + "_data_retrieval.html")
	require.NoError(t, err)

	err = dqotChart.Render(dqotFile)
	require.NoError(t, err)

	ipotFile, err := os.Create(baseName + "_in_processing.html")
	require.NoError(t, err)

	err = ipotChart.Render(ipotFile)
	require.NoError(t, err)

	wtflFile, err := os.Create(baseName + "_write_feedback_latency.html")
	require.NoError(t, err)

	err = wtflChart.Render(wtflFile)
	require.NoError(t, err)
}

func TestDisciplineFairUnevenProcessingTime(t *testing.T) {
	testDisciplineFairUnevenProcessingTime(t, 1, true)
	testDisciplineFairUnevenProcessingTime(t, 10, true)
	testDisciplineFairUnevenProcessingTime(t, 1, false)
	testDisciplineFairUnevenProcessingTime(t, 10, false)
}

func testUnmanagedEven(t *testing.T, factor uint, inputBuffered bool) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	handlersQuantity := uint(6) * factor

	measurerOpts := measurerOpts{
		HandlersQuantity: handlersQuantity,
		NoFeedback:       true,
		UnbufferedInput:  !inputBuffered,
	}

	measurer := newMeasurer(measurerOpts)
	defer measurer.Finalize()

	measurer.AddWrite(1, 4000*factor)

	measurer.AddWrite(2, 500*factor)
	measurer.AddWaitDevastation(2)
	measurer.AddDelay(2, 2*time.Second)
	measurer.AddWrite(2, 500*factor)
	measurer.AddWaitDevastation(2)
	measurer.AddDelay(2, 4*time.Second)
	measurer.AddWrite(2, 1000*factor)
	measurer.AddWaitDevastation(2)
	measurer.AddDelay(2, 2*time.Second)
	measurer.AddWrite(2, 2000*factor)

	measurer.AddWrite(3, 500*factor)
	measurer.AddWaitDevastation(3)
	measurer.AddDelay(3, 5*time.Second)
	measurer.AddWrite(3, 4000*factor)

	measurer.SetProcessDelay(1, 10*time.Millisecond)
	measurer.SetProcessDelay(2, 10*time.Millisecond)
	measurer.SetProcessDelay(3, 10*time.Millisecond)

	unmanagedOpts := unmanagedOpts[uint]{
		Inputs: measurer.GetInputs(),
		Output: measurer.GetOutput(),
	}

	unmanaged, err := newUnmanaged(unmanagedOpts)
	require.NoError(t, err)

	defer unmanaged.Stop()

	measures := measurer.Play(unmanaged)

	received := filterByKind(measures, measureKindReceived)

	dqot, dqotX := convertToLineEcharts(
		calcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := convertToLineEcharts(
		calcInProcessing(measures, 100*time.Millisecond),
		1*time.Second,
	)

	dqotChart := charts.NewLine()
	ipotChart := charts.NewLine()

	subtitle := fmt.Sprintf(
		"Unmanaged, even time processing, "+
			"handlers quantity: %d, buffered: %t, time: %s",
		handlersQuantity,
		inputBuffered,
		time.Now().Format(time.RFC3339),
	)

	dqotChart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "Data retrieval graph",
				Subtitle: subtitle,
			},
		),
	)

	ipotChart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "In processing graph",
				Subtitle: subtitle,
			},
		),
	)

	dqotChart.SetXAxis(dqotX).
		AddSeries("3", dqot[3]).
		AddSeries("2", dqot[2]).
		AddSeries("1", dqot[1])

	ipotChart.SetXAxis(ipotX).
		AddSeries("3", ipot[3]).
		AddSeries("2", ipot[2]).
		AddSeries("1", ipot[1])

	baseName := "graph_unmanaged_even_" + strconv.Itoa(int(handlersQuantity)) +
		"_buffered_" + strconv.FormatBool(inputBuffered)

	dqotFile, err := os.Create(baseName + "_data_retrieval.html")
	require.NoError(t, err)

	err = dqotChart.Render(dqotFile)
	require.NoError(t, err)

	ipotFile, err := os.Create(baseName + "_in_processing.html")
	require.NoError(t, err)

	err = ipotChart.Render(ipotFile)
	require.NoError(t, err)
}

func TestUnmanagedEven(t *testing.T) {
	testUnmanagedEven(t, 1, true)
	testUnmanagedEven(t, 10, true)
	testUnmanagedEven(t, 1, false)
	testUnmanagedEven(t, 10, false)
}

func testUnmanagedUneven(t *testing.T, factor uint, inputBuffered bool) {
	if os.Getenv(consts.EnableGraphsEnv) == "" {
		t.SkipNow()
	}

	handlersQuantity := uint(6) * factor

	measurerOpts := measurerOpts{
		HandlersQuantity: handlersQuantity,
		NoFeedback:       true,
		UnbufferedInput:  !inputBuffered,
	}

	measurer := newMeasurer(measurerOpts)
	defer measurer.Finalize()

	measurer.AddWrite(1, 500*factor)

	measurer.AddWrite(2, 100*factor)
	measurer.AddWaitDevastation(2)
	measurer.AddDelay(2, 2*time.Second)
	measurer.AddWrite(2, 100*factor)
	measurer.AddWaitDevastation(2)
	measurer.AddDelay(2, 2*time.Second)
	measurer.AddWrite(2, 200*factor)
	measurer.AddWaitDevastation(2)
	measurer.AddDelay(2, 2*time.Second)
	measurer.AddWrite(2, 400*factor)

	measurer.AddWrite(3, 100*factor)
	measurer.AddWaitDevastation(3)
	measurer.AddDelay(3, 6*time.Second)
	measurer.AddWrite(3, 1350*factor)

	measurer.SetProcessDelay(1, 100*time.Millisecond)
	measurer.SetProcessDelay(2, 50*time.Millisecond)
	measurer.SetProcessDelay(3, 10*time.Millisecond)

	unmanagedOpts := unmanagedOpts[uint]{
		Inputs: measurer.GetInputs(),
		Output: measurer.GetOutput(),
	}

	unmanaged, err := newUnmanaged(unmanagedOpts)
	require.NoError(t, err)

	defer unmanaged.Stop()

	measures := measurer.Play(unmanaged)

	received := filterByKind(measures, measureKindReceived)

	dqot, dqotX := convertToLineEcharts(
		calcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := convertToLineEcharts(
		calcInProcessing(measures, 100*time.Millisecond),
		1*time.Second,
	)

	dqotChart := charts.NewLine()
	ipotChart := charts.NewLine()

	subtitle := fmt.Sprintf(
		"Unmanaged, uneven time processing, "+
			"handlers quantity: %d, buffered: %t, time: %s",
		handlersQuantity,
		inputBuffered,
		time.Now().Format(time.RFC3339),
	)

	dqotChart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "Data retrieval graph",
				Subtitle: subtitle,
			},
		),
	)

	ipotChart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "In processing graph",
				Subtitle: subtitle,
			},
		),
	)

	dqotChart.SetXAxis(dqotX).
		AddSeries("3", dqot[3]).
		AddSeries("2", dqot[2]).
		AddSeries("1", dqot[1])

	ipotChart.SetXAxis(ipotX).
		AddSeries("3", ipot[3]).
		AddSeries("2", ipot[2]).
		AddSeries("1", ipot[1])

	baseName := "graph_unmanaged_uneven_" + strconv.Itoa(int(handlersQuantity)) +
		"_buffered_" + strconv.FormatBool(inputBuffered)

	dqotFile, err := os.Create(baseName + "_data_retrieval.html")
	require.NoError(t, err)

	err = dqotChart.Render(dqotFile)
	require.NoError(t, err)

	ipotFile, err := os.Create(baseName + "_in_processing.html")
	require.NoError(t, err)

	err = ipotChart.Render(ipotFile)
	require.NoError(t, err)
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

	measurerOpts := measurerOpts{
		HandlersQuantity: handlersQuantity,
	}

	measurer := newMeasurer(measurerOpts)
	defer measurer.Finalize()

	measurer.AddWrite(1, 4000)

	measurer.AddWrite(2, 500)
	measurer.AddWaitDevastation(2)
	measurer.AddDelay(2, 2*time.Second)
	measurer.AddWrite(2, 500)
	measurer.AddWaitDevastation(2)
	measurer.AddDelay(2, 4*time.Second)
	measurer.AddWrite(2, 1000)
	measurer.AddWaitDevastation(2)
	measurer.AddDelay(2, 2*time.Second)
	measurer.AddWrite(2, 2000)

	measurer.AddWrite(3, 500)
	measurer.AddWaitDevastation(3)
	measurer.AddDelay(3, 5*time.Second)
	measurer.AddWrite(3, 4000)

	measurer.AddWrite(4, 500)
	measurer.AddWaitDevastation(3)
	measurer.AddDelay(4, 5*time.Second)
	measurer.AddWrite(4, 4000)

	measurer.SetProcessDelay(1, 10*time.Millisecond)
	measurer.SetProcessDelay(2, 10*time.Millisecond)
	measurer.SetProcessDelay(3, 10*time.Millisecond)
	measurer.SetProcessDelay(4, 10*time.Millisecond)

	disciplineOpts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         measurer.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           measurer.GetInputs(),
		Output:           measurer.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	measures := measurer.Play(discipline)

	received := filterByKind(measures, measureKindReceived)

	dqot, dqotX := convertToLineEcharts(
		calcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	dqotChart := charts.NewLine()

	subtitle := fmt.Sprintf(
		"Fair divider, even time processing, "+
			"significant dividing error, "+
			"handlers quantity: %d, buffered: %t, time: %s",
		handlersQuantity,
		!measurerOpts.UnbufferedInput,
		time.Now().Format(time.RFC3339),
	)

	dqotChart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "Data retrieval graph",
				Subtitle: subtitle,
			},
		),
	)

	dqotChart.SetXAxis(dqotX).
		AddSeries("4", dqot[4]).
		AddSeries("3", dqot[3]).
		AddSeries("2", dqot[2]).
		AddSeries("1", dqot[1])

	baseName := "graph_fair_even_" + strconv.Itoa(int(handlersQuantity)) +
		"_buffered_" + strconv.FormatBool(!measurerOpts.UnbufferedInput) + "_dividing_error"

	dqotFile, err := os.Create(baseName + "_data_retrieval.html")
	require.NoError(t, err)

	err = dqotChart.Render(dqotFile)
	require.NoError(t, err)
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

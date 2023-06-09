package priority

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
	"github.com/stretchr/testify/require"
)

func testDisciplineRateEvenProcessingTime(t *testing.T, factor uint, inputBuffered bool) {
	if os.Getenv("CQOS_ENABLE_GRAPHS") == "" {
		t.SkipNow()
	}

	handlersQuantity := uint(6) * factor

	gaugerOpts := gaugerOpts{
		HandlersQuantity: handlersQuantity,
		NoInputBuffer:    !inputBuffered,
	}

	gauger := newGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 4100*factor)

	gauger.AddWrite(2, 1500*factor)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 750*factor)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 4*time.Second)
	gauger.AddWrite(2, 700*factor)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 3*time.Second)
	gauger.AddWrite(2, 1200*factor)

	gauger.AddWrite(3, 1000*factor)
	gauger.AddWaitDevastation(3)
	gauger.AddDelay(3, 8*time.Second)
	gauger.AddWrite(3, 3700*factor)

	gauger.SetProcessDelay(1, 10*time.Millisecond)
	gauger.SetProcessDelay(2, 10*time.Millisecond)
	gauger.SetProcessDelay(3, 10*time.Millisecond)

	disciplineOpts := Opts[uint]{
		Divider:          RateDivider,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	gauges := gauger.Play(context.Background())

	received := filterByKind(gauges, gaugeKindReceived)

	dqot, dqotX := convertQuantityOverTimeToLineEcharts(
		calcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := convertQuantityOverTimeToLineEcharts(
		calcInProcessing(gauges, 100*time.Millisecond),
		1*time.Second,
	)

	wtfl, wtflX := convertQuantityOverTimeToBarEcharts(
		calcWriteToFeedbackLatency(gauges, 100*time.Nanosecond),
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
	if os.Getenv("CQOS_ENABLE_GRAPHS") == "" {
		t.SkipNow()
	}

	handlersQuantity := uint(6) * factor

	gaugerOpts := gaugerOpts{
		HandlersQuantity: handlersQuantity,
		NoInputBuffer:    !inputBuffered,
	}

	gauger := newGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 430*factor)

	gauger.AddWrite(2, 250*factor)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 100*factor)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 4*time.Second)
	gauger.AddWrite(2, 150*factor)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 300*factor)

	gauger.AddWrite(3, 1000*factor)
	gauger.AddWaitDevastation(3)
	gauger.AddDelay(3, 8*time.Second)
	gauger.AddWrite(3, 3500*factor)

	gauger.SetProcessDelay(1, 100*time.Millisecond)
	gauger.SetProcessDelay(2, 50*time.Millisecond)
	gauger.SetProcessDelay(3, 10*time.Millisecond)

	disciplineOpts := Opts[uint]{
		Divider:          RateDivider,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	gauges := gauger.Play(context.Background())

	received := filterByKind(gauges, gaugeKindReceived)

	dqot, dqotX := convertQuantityOverTimeToLineEcharts(
		calcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := convertQuantityOverTimeToLineEcharts(
		calcInProcessing(gauges, 100*time.Millisecond),
		1*time.Second,
	)

	wtfl, wtflX := convertQuantityOverTimeToBarEcharts(
		calcWriteToFeedbackLatency(gauges, 100*time.Nanosecond),
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
	if os.Getenv("CQOS_ENABLE_GRAPHS") == "" {
		t.SkipNow()
	}

	handlersQuantity := uint(6) * factor

	gaugerOpts := gaugerOpts{
		HandlersQuantity: handlersQuantity,
		NoInputBuffer:    !inputBuffered,
	}

	gauger := newGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 4000*factor)

	gauger.AddWrite(2, 500*factor)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 500*factor)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 4*time.Second)
	gauger.AddWrite(2, 1000*factor)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 2000*factor)

	gauger.AddWrite(3, 500*factor)
	gauger.AddWaitDevastation(3)
	gauger.AddDelay(3, 5*time.Second)
	gauger.AddWrite(3, 4000*factor)

	gauger.SetProcessDelay(1, 10*time.Millisecond)
	gauger.SetProcessDelay(2, 10*time.Millisecond)
	gauger.SetProcessDelay(3, 10*time.Millisecond)

	disciplineOpts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	gauges := gauger.Play(context.Background())

	received := filterByKind(gauges, gaugeKindReceived)

	dqot, dqotX := convertQuantityOverTimeToLineEcharts(
		calcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := convertQuantityOverTimeToLineEcharts(
		calcInProcessing(gauges, 100*time.Millisecond),
		1*time.Second,
	)

	wtfl, wtflX := convertQuantityOverTimeToBarEcharts(
		calcWriteToFeedbackLatency(gauges, 100*time.Nanosecond),
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
	if os.Getenv("CQOS_ENABLE_GRAPHS") == "" {
		t.SkipNow()
	}

	handlersQuantity := uint(6) * factor

	gaugerOpts := gaugerOpts{
		HandlersQuantity: handlersQuantity,
		NoInputBuffer:    !inputBuffered,
	}

	gauger := newGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 450*factor)

	gauger.AddWrite(2, 100*factor)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 100*factor)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 4*time.Second)
	gauger.AddWrite(2, 200*factor)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 400*factor)

	gauger.AddWrite(3, 500*factor)
	gauger.AddWaitDevastation(3)
	gauger.AddDelay(3, 6*time.Second)
	gauger.AddWrite(3, 3000*factor)

	gauger.SetProcessDelay(1, 100*time.Millisecond)
	gauger.SetProcessDelay(2, 50*time.Millisecond)
	gauger.SetProcessDelay(3, 10*time.Millisecond)

	disciplineOpts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	gauges := gauger.Play(context.Background())

	received := filterByKind(gauges, gaugeKindReceived)

	dqot, dqotX := convertQuantityOverTimeToLineEcharts(
		calcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := convertQuantityOverTimeToLineEcharts(
		calcInProcessing(gauges, 100*time.Millisecond),
		1*time.Second,
	)

	wtfl, wtflX := convertQuantityOverTimeToBarEcharts(
		calcWriteToFeedbackLatency(gauges, 100*time.Nanosecond),
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
	if os.Getenv("CQOS_ENABLE_GRAPHS") == "" {
		t.SkipNow()
	}

	handlersQuantity := uint(6) * factor

	gaugerOpts := gaugerOpts{
		HandlersQuantity: handlersQuantity,
		NoFeedback:       true,
		NoInputBuffer:    !inputBuffered,
	}

	gauger := newGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 4000*factor)

	gauger.AddWrite(2, 500*factor)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 500*factor)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 4*time.Second)
	gauger.AddWrite(2, 1000*factor)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 2000*factor)

	gauger.AddWrite(3, 500*factor)
	gauger.AddWaitDevastation(3)
	gauger.AddDelay(3, 5*time.Second)
	gauger.AddWrite(3, 4000*factor)

	gauger.SetProcessDelay(1, 10*time.Millisecond)
	gauger.SetProcessDelay(2, 10*time.Millisecond)
	gauger.SetProcessDelay(3, 10*time.Millisecond)

	unmanagedOpts := unmanagedOpts[uint]{
		Inputs: gauger.GetInputs(),
		Output: gauger.GetOutput(),
	}

	unmanaged, err := newUnmanaged(unmanagedOpts)
	require.NoError(t, err)

	defer unmanaged.Stop()

	gauges := gauger.Play(context.Background())

	received := filterByKind(gauges, gaugeKindReceived)

	dqot, dqotX := convertQuantityOverTimeToLineEcharts(
		calcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := convertQuantityOverTimeToLineEcharts(
		calcInProcessing(gauges, 100*time.Millisecond),
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
	if os.Getenv("CQOS_ENABLE_GRAPHS") == "" {
		t.SkipNow()
	}

	handlersQuantity := uint(6) * factor

	gaugerOpts := gaugerOpts{
		HandlersQuantity: handlersQuantity,
		NoFeedback:       true,
		NoInputBuffer:    !inputBuffered,
	}

	gauger := newGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 500*factor)

	gauger.AddWrite(2, 100*factor)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 100*factor)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 200*factor)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 400*factor)

	gauger.AddWrite(3, 100*factor)
	gauger.AddWaitDevastation(3)
	gauger.AddDelay(3, 6*time.Second)
	gauger.AddWrite(3, 1350*factor)

	gauger.SetProcessDelay(1, 100*time.Millisecond)
	gauger.SetProcessDelay(2, 50*time.Millisecond)
	gauger.SetProcessDelay(3, 10*time.Millisecond)

	unmanagedOpts := unmanagedOpts[uint]{
		Inputs: gauger.GetInputs(),
		Output: gauger.GetOutput(),
	}

	unmanaged, err := newUnmanaged(unmanagedOpts)
	require.NoError(t, err)

	defer unmanaged.Stop()

	gauges := gauger.Play(context.Background())

	received := filterByKind(gauges, gaugeKindReceived)

	dqot, dqotX := convertQuantityOverTimeToLineEcharts(
		calcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := convertQuantityOverTimeToLineEcharts(
		calcInProcessing(gauges, 100*time.Millisecond),
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

func BenchmarkDisciplineFair(b *testing.B) {
	handlersQuantity := uint(600)

	gaugerOpts := gaugerOpts{
		DisableGauges:    true,
		HandlersQuantity: handlersQuantity,
	}

	gauger := newGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 5000000)
	gauger.AddWrite(2, 5000000)
	gauger.AddWrite(3, 5000000)

	disciplineOpts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(b, err)

	defer discipline.Stop()

	_ = gauger.Play(context.Background())
}

func BenchmarkDisciplineRate(b *testing.B) {
	handlersQuantity := uint(600)

	gaugerOpts := gaugerOpts{
		DisableGauges:    true,
		HandlersQuantity: handlersQuantity,
	}

	gauger := newGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 5000000)
	gauger.AddWrite(2, 5000000)
	gauger.AddWrite(3, 5000000)

	disciplineOpts := Opts[uint]{
		Divider:          RateDivider,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(b, err)

	defer discipline.Stop()

	_ = gauger.Play(context.Background())
}

func BenchmarkDisciplineFairUnbuffered(b *testing.B) {
	handlersQuantity := uint(600)

	gaugerOpts := gaugerOpts{
		DisableGauges:    true,
		HandlersQuantity: handlersQuantity,
		NoInputBuffer:    true,
	}

	gauger := newGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 5000000)
	gauger.AddWrite(2, 5000000)
	gauger.AddWrite(3, 5000000)

	disciplineOpts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(b, err)

	defer discipline.Stop()

	_ = gauger.Play(context.Background())
}

func BenchmarkDisciplineRateUnbuffered(b *testing.B) {
	handlersQuantity := uint(600)

	gaugerOpts := gaugerOpts{
		DisableGauges:    true,
		HandlersQuantity: handlersQuantity,
		NoInputBuffer:    true,
	}

	gauger := newGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 5000000)
	gauger.AddWrite(2, 5000000)
	gauger.AddWrite(3, 5000000)

	disciplineOpts := Opts[uint]{
		Divider:          RateDivider,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(b, err)

	defer discipline.Stop()

	_ = gauger.Play(context.Background())
}

func TestDisciplineRate(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := gaugerOpts{
		HandlersQuantity: handlersQuantity,
	}

	gauger := newGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 100000)
	gauger.AddWrite(2, 100000)
	gauger.AddWrite(3, 100000)

	disciplineOpts := Opts[uint]{
		Divider:          RateDivider,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	gauges := gauger.Play(context.Background())

	require.Equal(t, int(gauger.CalcExpectedGuagesQuantity()), len(filterByKind(gauges, gaugeKindReceived)))
}

func TestDisciplineFair(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := gaugerOpts{
		HandlersQuantity: handlersQuantity,
	}

	gauger := newGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 100000)
	gauger.AddWrite(2, 100000)
	gauger.AddWrite(3, 100000)

	disciplineOpts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	gauges := gauger.Play(context.Background())

	require.Equal(t, int(gauger.CalcExpectedGuagesQuantity()), len(filterByKind(gauges, gaugeKindReceived)))
}

func TestDisciplineRateUnbuffered(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := gaugerOpts{
		HandlersQuantity: handlersQuantity,
		NoInputBuffer:    true,
	}

	gauger := newGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 100000)
	gauger.AddWrite(2, 100000)
	gauger.AddWrite(3, 100000)

	disciplineOpts := Opts[uint]{
		Divider:          RateDivider,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	gauges := gauger.Play(context.Background())

	require.Equal(t, int(gauger.CalcExpectedGuagesQuantity()), len(filterByKind(gauges, gaugeKindReceived)))
}

func TestDisciplineFairUnbuffered(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := gaugerOpts{
		HandlersQuantity: handlersQuantity,
		NoInputBuffer:    true,
	}

	gauger := newGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 100000)
	gauger.AddWrite(2, 100000)
	gauger.AddWrite(3, 100000)

	disciplineOpts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	gauges := gauger.Play(context.Background())

	require.Equal(t, int(gauger.CalcExpectedGuagesQuantity()), len(filterByKind(gauges, gaugeKindReceived)))
}

func TestDisciplineOptsValidation(t *testing.T) {
	handlersQuantity := uint(6)

	disciplineOpts := Opts[uint]{
		Feedback:         make(chan uint),
		HandlersQuantity: handlersQuantity,
		Output:           make(chan Prioritized[uint]),
	}

	_, err := New(disciplineOpts)
	require.Error(t, err)

	disciplineOpts = Opts[uint]{
		Divider:          FairDivider,
		HandlersQuantity: handlersQuantity,
		Output:           make(chan Prioritized[uint]),
	}

	_, err = New(disciplineOpts)
	require.Error(t, err)

	disciplineOpts = Opts[uint]{
		Divider:          FairDivider,
		Feedback:         make(chan uint),
		HandlersQuantity: handlersQuantity,
	}

	_, err = New(disciplineOpts)
	require.Error(t, err)
}

func TestDisciplineAddRemoveInput(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := gaugerOpts{
		HandlersQuantity: handlersQuantity,
	}

	gauger := newGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 1000000)
	gauger.AddWrite(2, 1000000)
	gauger.AddWrite(3, 1000000)

	inputs := gauger.GetInputs()

	disciplineOpts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	waiter := make(chan bool)

	go func() {
		defer close(waiter)

		discipline.AddInput(inputs[2], 2)
		discipline.AddInput(inputs[2], 2)
		discipline.AddInput(inputs[1], 1)
		discipline.AddInput(inputs[1], 1)
		discipline.AddInput(inputs[3], 3)
		discipline.AddInput(inputs[3], 3)

		four := make(chan uint, 10)
		close(four)

		discipline.AddInput(four, 4)

		four = make(chan uint)
		close(four)

		discipline.AddInput(four, 4)

		time.Sleep(1 * time.Second)

		discipline.RemoveInput(3)
		discipline.RemoveInput(3)
		discipline.RemoveInput(2)
		discipline.RemoveInput(2)
		discipline.RemoveInput(1)
		discipline.RemoveInput(1)
		discipline.RemoveInput(4)
		discipline.RemoveInput(4)

		discipline.AddInput(inputs[2], 6)
		discipline.AddInput(inputs[1], 5)
		discipline.AddInput(inputs[3], 7)
	}()

	gauges := gauger.Play(context.Background())

	<-waiter

	require.Equal(t, int(gauger.CalcExpectedGuagesQuantity()), len(filterByKind(gauges, gaugeKindReceived)))
}

func TestDisciplineBadDivider(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := gaugerOpts{
		HandlersQuantity: handlersQuantity,
	}

	gauger := newGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 100000)
	gauger.AddWrite(2, 100000)
	gauger.AddWrite(3, 100000)

	divider := func(priorities []uint, dividend uint, distribution map[uint]uint) map[uint]uint {
		out := FairDivider(priorities, dividend, distribution)

		for priority, quantity := range out {
			out[priority] = 2 * quantity
		}

		return out
	}

	disciplineOpts := Opts[uint]{
		Divider:          divider,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if <-discipline.Err() != nil {
			cancel()
		}
	}()

	gauges := gauger.Play(ctx)

	require.NotEqual(t, int(gauger.CalcExpectedGuagesQuantity()), len(filterByKind(gauges, gaugeKindReceived)))
}

func TestDisciplineStop(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := gaugerOpts{
		HandlersQuantity: handlersQuantity,
	}

	gauger := newGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 100000)
	gauger.AddWrite(2, 100000)
	gauger.AddWrite(3, 100000)

	gauger.SetProcessDelay(1, 10*time.Microsecond)
	gauger.SetProcessDelay(2, 10*time.Microsecond)
	gauger.SetProcessDelay(3, 10*time.Microsecond)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	disciplineOpts := Opts[uint]{
		Ctx:              ctx,
		Divider:          RateDivider,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()
	defer discipline.Stop()

	gauges := gauger.Play(ctx)

	require.NotEqual(t, int(gauger.CalcExpectedGuagesQuantity()), len(filterByKind(gauges, gaugeKindReceived)))
}

func TestDisciplineGracefulStop(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := gaugerOpts{
		HandlersQuantity: handlersQuantity,
	}

	gauger := newGauger(gaugerOpts)

	gauger.AddWrite(1, 100000)
	gauger.AddWrite(2, 100000)
	gauger.AddWrite(3, 100000)

	disciplineOpts := Opts[uint]{
		Divider:          RateDivider,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()

		discipline.GracefulStop()
	}()

	gauges := gauger.Play(ctx)
	gauger.Finalize()

	require.Equal(t, int(gauger.CalcExpectedGuagesQuantity()), len(filterByKind(gauges, gaugeKindReceived)))
}

func TestDisciplineRateOverQuantity(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := gaugerOpts{
		HandlersQuantity: 2 * handlersQuantity,
	}

	gauger := newGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 100000)
	gauger.AddWrite(2, 100000)
	gauger.AddWrite(3, 100000)

	disciplineOpts := Opts[uint]{
		Divider:          RateDivider,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	gauges := gauger.Play(context.Background())

	quantities := calcInProcessing(gauges, 100*time.Millisecond)

	for priority := range quantities {
		for id := range quantities[priority] {
			require.LessOrEqual(t, quantities[priority][id].Quantity, handlersQuantity)
		}
	}
}

func TestDisciplineFairOverQuantity(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := gaugerOpts{
		HandlersQuantity: 2 * handlersQuantity,
	}

	gauger := newGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 1000000)
	gauger.AddWrite(2, 100000)
	gauger.AddWrite(3, 10000)

	disciplineOpts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	gauges := gauger.Play(context.Background())

	quantities := calcInProcessing(gauges, 100*time.Millisecond)

	for priority := range quantities {
		for id := range quantities[priority] {
			require.LessOrEqual(t, quantities[priority][id].Quantity, handlersQuantity)
		}
	}
}

func TestDisciplineRateFatalDividingError(t *testing.T) {
	handlersQuantity := uint(5)

	gaugerOpts := gaugerOpts{
		HandlersQuantity: handlersQuantity,
	}

	gauger := newGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 100000)
	gauger.AddWrite(2, 100000)
	gauger.AddWrite(3, 100000)

	disciplineOpts := Opts[uint]{
		Divider:          RateDivider,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	gauges := gauger.Play(context.Background())

	require.Equal(t, int(gauger.CalcExpectedGuagesQuantity()), len(filterByKind(gauges, gaugeKindReceived)))
}

func TestDisciplineFairFatalDividingError(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := gaugerOpts{
		HandlersQuantity: handlersQuantity,
	}

	gauger := newGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 100000)
	gauger.AddWrite(2, 100000)
	gauger.AddWrite(3, 100000)

	disciplineOpts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	gauges := gauger.Play(context.Background())

	require.Equal(t, int(gauger.CalcExpectedGuagesQuantity()), len(filterByKind(gauges, gaugeKindReceived)))
}

func testDisciplineFairEvenProcessingTimeDividingError(t *testing.T, handlersQuantity uint) {
	if os.Getenv("CQOS_ENABLE_GRAPHS") == "" {
		t.SkipNow()
	}

	gaugerOpts := gaugerOpts{
		HandlersQuantity: handlersQuantity,
	}

	gauger := newGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 4000)

	gauger.AddWrite(2, 500)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 500)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 4*time.Second)
	gauger.AddWrite(2, 1000)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 2000)

	gauger.AddWrite(3, 500)
	gauger.AddWaitDevastation(3)
	gauger.AddDelay(3, 5*time.Second)
	gauger.AddWrite(3, 4000)

	gauger.AddWrite(4, 500)
	gauger.AddWaitDevastation(3)
	gauger.AddDelay(4, 5*time.Second)
	gauger.AddWrite(4, 4000)

	gauger.SetProcessDelay(1, 10*time.Millisecond)
	gauger.SetProcessDelay(2, 10*time.Millisecond)
	gauger.SetProcessDelay(3, 10*time.Millisecond)
	gauger.SetProcessDelay(4, 10*time.Millisecond)

	disciplineOpts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	gauges := gauger.Play(context.Background())

	received := filterByKind(gauges, gaugeKindReceived)

	dqot, dqotX := convertQuantityOverTimeToLineEcharts(
		calcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	dqotChart := charts.NewLine()

	subtitle := fmt.Sprintf(
		"Fair divider, even time processing, "+
			"significant dividing error, "+
			"handlers quantity: %d, buffered: %t, time: %s",
		handlersQuantity,
		!gaugerOpts.NoInputBuffer,
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
		"_buffered_" + strconv.FormatBool(!gaugerOpts.NoInputBuffer) + "_dividing_error"

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

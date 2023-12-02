package priority

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/priority/divider"
	"github.com/akramarenkov/cqos/v2/priority/internal/gauger"
	"github.com/akramarenkov/cqos/v2/priority/internal/research"
	"github.com/akramarenkov/cqos/v2/priority/internal/unmanaged"

	"github.com/go-echarts/go-echarts/v2/charts"
	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
	"github.com/stretchr/testify/require"
)

func TestDisciplineOptsValidation(t *testing.T) {
	handlersQuantity := uint(6)

	disciplineOpts := Opts[uint]{
		HandlersQuantity: handlersQuantity,
	}

	_, err := New(disciplineOpts)
	require.Error(t, err)

	disciplineOpts = Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: handlersQuantity,
	}

	_, err = New(disciplineOpts)
	require.Error(t, err)

	disciplineOpts = Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: handlersQuantity,
		Inputs: map[uint]<-chan uint{
			1: make(chan uint),
		},
	}

	_, err = New(disciplineOpts)
	require.NoError(t, err)
}

func testDisciplineRateEvenProcessingTime(t *testing.T, factor uint, inputBuffered bool) {
	if os.Getenv("CQOS_ENABLE_GRAPHS") == "" {
		t.SkipNow()
	}

	handlersQuantity := uint(6) * factor

	gaugerOpts := gauger.Opts{
		HandlersQuantity: handlersQuantity,
		NoInputBuffer:    !inputBuffered,
	}

	ggr := gauger.New(gaugerOpts)

	ggr.AddWrite(1, 4100*factor)

	ggr.AddWrite(2, 1500*factor)
	ggr.AddWaitDevastation(2)
	ggr.AddDelay(2, 2*time.Second)
	ggr.AddWrite(2, 750*factor)
	ggr.AddWaitDevastation(2)
	ggr.AddDelay(2, 4*time.Second)
	ggr.AddWrite(2, 700*factor)
	ggr.AddWaitDevastation(2)
	ggr.AddDelay(2, 3*time.Second)
	ggr.AddWrite(2, 1200*factor)

	ggr.AddWrite(3, 1000*factor)
	ggr.AddWaitDevastation(3)
	ggr.AddDelay(3, 8*time.Second)
	ggr.AddWrite(3, 3700*factor)

	ggr.SetProcessDelay(1, 10*time.Millisecond)
	ggr.SetProcessDelay(2, 10*time.Millisecond)
	ggr.SetProcessDelay(3, 10*time.Millisecond)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Rate,
		HandlersQuantity: handlersQuantity,
		Inputs:           ggr.GetInputs(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	ggr.SetDiscipline(discipline)

	gauges := ggr.Play(context.Background())

	received := research.FilterByKind(gauges, gauger.GaugeKindReceived)

	dqot, dqotX := research.ConvertToLineEcharts(
		research.CalcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := research.ConvertToLineEcharts(
		research.CalcInProcessing(gauges, 100*time.Millisecond),
		1*time.Second,
	)

	wtfl, wtflX := research.ConvertToBarEcharts(
		research.CalcWriteToFeedbackLatency(gauges, 100*time.Nanosecond),
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

	gaugerOpts := gauger.Opts{
		HandlersQuantity: handlersQuantity,
		NoInputBuffer:    !inputBuffered,
	}

	ggr := gauger.New(gaugerOpts)

	ggr.AddWrite(1, 430*factor)

	ggr.AddWrite(2, 250*factor)
	ggr.AddWaitDevastation(2)
	ggr.AddDelay(2, 2*time.Second)
	ggr.AddWrite(2, 100*factor)
	ggr.AddWaitDevastation(2)
	ggr.AddDelay(2, 4*time.Second)
	ggr.AddWrite(2, 150*factor)
	ggr.AddWaitDevastation(2)
	ggr.AddDelay(2, 2*time.Second)
	ggr.AddWrite(2, 300*factor)

	ggr.AddWrite(3, 1000*factor)
	ggr.AddWaitDevastation(3)
	ggr.AddDelay(3, 8*time.Second)
	ggr.AddWrite(3, 3500*factor)

	ggr.SetProcessDelay(1, 100*time.Millisecond)
	ggr.SetProcessDelay(2, 50*time.Millisecond)
	ggr.SetProcessDelay(3, 10*time.Millisecond)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Rate,
		HandlersQuantity: handlersQuantity,
		Inputs:           ggr.GetInputs(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	ggr.SetDiscipline(discipline)

	gauges := ggr.Play(context.Background())

	received := research.FilterByKind(gauges, gauger.GaugeKindReceived)

	dqot, dqotX := research.ConvertToLineEcharts(
		research.CalcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := research.ConvertToLineEcharts(
		research.CalcInProcessing(gauges, 100*time.Millisecond),
		1*time.Second,
	)

	wtfl, wtflX := research.ConvertToBarEcharts(
		research.CalcWriteToFeedbackLatency(gauges, 100*time.Nanosecond),
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

	gaugerOpts := gauger.Opts{
		HandlersQuantity: handlersQuantity,
		NoInputBuffer:    !inputBuffered,
	}

	ggr := gauger.New(gaugerOpts)

	ggr.AddWrite(1, 4000*factor)

	ggr.AddWrite(2, 500*factor)
	ggr.AddWaitDevastation(2)
	ggr.AddDelay(2, 2*time.Second)
	ggr.AddWrite(2, 500*factor)
	ggr.AddWaitDevastation(2)
	ggr.AddDelay(2, 4*time.Second)
	ggr.AddWrite(2, 1000*factor)
	ggr.AddWaitDevastation(2)
	ggr.AddDelay(2, 2*time.Second)
	ggr.AddWrite(2, 2000*factor)

	ggr.AddWrite(3, 500*factor)
	ggr.AddWaitDevastation(3)
	ggr.AddDelay(3, 5*time.Second)
	ggr.AddWrite(3, 4000*factor)

	ggr.SetProcessDelay(1, 10*time.Millisecond)
	ggr.SetProcessDelay(2, 10*time.Millisecond)
	ggr.SetProcessDelay(3, 10*time.Millisecond)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: handlersQuantity,
		Inputs:           ggr.GetInputs(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	ggr.SetDiscipline(discipline)

	gauges := ggr.Play(context.Background())

	received := research.FilterByKind(gauges, gauger.GaugeKindReceived)

	dqot, dqotX := research.ConvertToLineEcharts(
		research.CalcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := research.ConvertToLineEcharts(
		research.CalcInProcessing(gauges, 100*time.Millisecond),
		1*time.Second,
	)

	wtfl, wtflX := research.ConvertToBarEcharts(
		research.CalcWriteToFeedbackLatency(gauges, 100*time.Nanosecond),
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

	gaugerOpts := gauger.Opts{
		HandlersQuantity: handlersQuantity,
		NoInputBuffer:    !inputBuffered,
	}

	ggr := gauger.New(gaugerOpts)

	ggr.AddWrite(1, 450*factor)

	ggr.AddWrite(2, 100*factor)
	ggr.AddWaitDevastation(2)
	ggr.AddDelay(2, 2*time.Second)
	ggr.AddWrite(2, 100*factor)
	ggr.AddWaitDevastation(2)
	ggr.AddDelay(2, 4*time.Second)
	ggr.AddWrite(2, 200*factor)
	ggr.AddWaitDevastation(2)
	ggr.AddDelay(2, 2*time.Second)
	ggr.AddWrite(2, 400*factor)

	ggr.AddWrite(3, 500*factor)
	ggr.AddWaitDevastation(3)
	ggr.AddDelay(3, 6*time.Second)
	ggr.AddWrite(3, 3000*factor)

	ggr.SetProcessDelay(1, 100*time.Millisecond)
	ggr.SetProcessDelay(2, 50*time.Millisecond)
	ggr.SetProcessDelay(3, 10*time.Millisecond)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: handlersQuantity,
		Inputs:           ggr.GetInputs(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	ggr.SetDiscipline(discipline)

	gauges := ggr.Play(context.Background())

	received := research.FilterByKind(gauges, gauger.GaugeKindReceived)

	dqot, dqotX := research.ConvertToLineEcharts(
		research.CalcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := research.ConvertToLineEcharts(
		research.CalcInProcessing(gauges, 100*time.Millisecond),
		1*time.Second,
	)

	wtfl, wtflX := research.ConvertToBarEcharts(
		research.CalcWriteToFeedbackLatency(gauges, 100*time.Nanosecond),
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

	gaugerOpts := gauger.Opts{
		HandlersQuantity: handlersQuantity,
		NoFeedback:       true,
		NoInputBuffer:    !inputBuffered,
	}

	ggr := gauger.New(gaugerOpts)

	ggr.AddWrite(1, 4000*factor)

	ggr.AddWrite(2, 500*factor)
	ggr.AddWaitDevastation(2)
	ggr.AddDelay(2, 2*time.Second)
	ggr.AddWrite(2, 500*factor)
	ggr.AddWaitDevastation(2)
	ggr.AddDelay(2, 4*time.Second)
	ggr.AddWrite(2, 1000*factor)
	ggr.AddWaitDevastation(2)
	ggr.AddDelay(2, 2*time.Second)
	ggr.AddWrite(2, 2000*factor)

	ggr.AddWrite(3, 500*factor)
	ggr.AddWaitDevastation(3)
	ggr.AddDelay(3, 5*time.Second)
	ggr.AddWrite(3, 4000*factor)

	ggr.SetProcessDelay(1, 10*time.Millisecond)
	ggr.SetProcessDelay(2, 10*time.Millisecond)
	ggr.SetProcessDelay(3, 10*time.Millisecond)

	unmanagedOpts := unmanaged.Opts[uint]{
		Inputs: ggr.GetInputs(),
	}

	unmanaged, err := unmanaged.New(unmanagedOpts)
	require.NoError(t, err)

	ggr.SetDiscipline(unmanaged)

	gauges := ggr.Play(context.Background())

	received := research.FilterByKind(gauges, gauger.GaugeKindReceived)

	dqot, dqotX := research.ConvertToLineEcharts(
		research.CalcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := research.ConvertToLineEcharts(
		research.CalcInProcessing(gauges, 100*time.Millisecond),
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

	gaugerOpts := gauger.Opts{
		HandlersQuantity: handlersQuantity,
		NoFeedback:       true,
		NoInputBuffer:    !inputBuffered,
	}

	ggr := gauger.New(gaugerOpts)

	ggr.AddWrite(1, 500*factor)

	ggr.AddWrite(2, 100*factor)
	ggr.AddWaitDevastation(2)
	ggr.AddDelay(2, 2*time.Second)
	ggr.AddWrite(2, 100*factor)
	ggr.AddWaitDevastation(2)
	ggr.AddDelay(2, 2*time.Second)
	ggr.AddWrite(2, 200*factor)
	ggr.AddWaitDevastation(2)
	ggr.AddDelay(2, 2*time.Second)
	ggr.AddWrite(2, 400*factor)

	ggr.AddWrite(3, 100*factor)
	ggr.AddWaitDevastation(3)
	ggr.AddDelay(3, 6*time.Second)
	ggr.AddWrite(3, 1350*factor)

	ggr.SetProcessDelay(1, 100*time.Millisecond)
	ggr.SetProcessDelay(2, 50*time.Millisecond)
	ggr.SetProcessDelay(3, 10*time.Millisecond)

	unmanagedOpts := unmanaged.Opts[uint]{
		Inputs: ggr.GetInputs(),
	}

	unmanaged, err := unmanaged.New(unmanagedOpts)
	require.NoError(t, err)

	ggr.SetDiscipline(unmanaged)

	gauges := ggr.Play(context.Background())

	received := research.FilterByKind(gauges, gauger.GaugeKindReceived)

	dqot, dqotX := research.ConvertToLineEcharts(
		research.CalcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := research.ConvertToLineEcharts(
		research.CalcInProcessing(gauges, 100*time.Millisecond),
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

	gaugerOpts := gauger.Opts{
		DisableGauges:    true,
		HandlersQuantity: handlersQuantity,
	}

	gauger := gauger.New(gaugerOpts)

	gauger.AddWrite(1, 5000000)
	gauger.AddWrite(2, 5000000)
	gauger.AddWrite(3, 5000000)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(b, err)

	gauger.SetDiscipline(discipline)

	_ = gauger.Play(context.Background())
}

func BenchmarkDisciplineRate(b *testing.B) {
	handlersQuantity := uint(600)

	gaugerOpts := gauger.Opts{
		DisableGauges:    true,
		HandlersQuantity: handlersQuantity,
	}

	gauger := gauger.New(gaugerOpts)

	gauger.AddWrite(1, 5000000)
	gauger.AddWrite(2, 5000000)
	gauger.AddWrite(3, 5000000)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Rate,
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(b, err)

	gauger.SetDiscipline(discipline)

	_ = gauger.Play(context.Background())
}

func BenchmarkDisciplineFairUnbuffered(b *testing.B) {
	handlersQuantity := uint(600)

	gaugerOpts := gauger.Opts{
		DisableGauges:    true,
		HandlersQuantity: handlersQuantity,
		NoInputBuffer:    true,
	}

	gauger := gauger.New(gaugerOpts)

	gauger.AddWrite(1, 5000000)
	gauger.AddWrite(2, 5000000)
	gauger.AddWrite(3, 5000000)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(b, err)

	gauger.SetDiscipline(discipline)

	_ = gauger.Play(context.Background())
}

func BenchmarkDisciplineRateUnbuffered(b *testing.B) {
	handlersQuantity := uint(600)

	gaugerOpts := gauger.Opts{
		DisableGauges:    true,
		HandlersQuantity: handlersQuantity,
		NoInputBuffer:    true,
	}

	gauger := gauger.New(gaugerOpts)

	gauger.AddWrite(1, 5000000)
	gauger.AddWrite(2, 5000000)
	gauger.AddWrite(3, 5000000)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Rate,
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(b, err)

	gauger.SetDiscipline(discipline)

	_ = gauger.Play(context.Background())
}

func TestDisciplineRate(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := gauger.Opts{
		HandlersQuantity: handlersQuantity,
	}

	ggr := gauger.New(gaugerOpts)

	ggr.AddWrite(1, 100000)
	ggr.AddWrite(2, 100000)
	ggr.AddWrite(3, 100000)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Rate,
		HandlersQuantity: handlersQuantity,
		Inputs:           ggr.GetInputs(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	ggr.SetDiscipline(discipline)

	gauges := ggr.Play(context.Background())

	require.Equal(t, int(ggr.CalcExpectedGuagesQuantity()), len(research.FilterByKind(gauges, gauger.GaugeKindReceived)))
}

func TestDisciplineFair(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := gauger.Opts{
		HandlersQuantity: handlersQuantity,
	}

	ggr := gauger.New(gaugerOpts)

	ggr.AddWrite(1, 100000)
	ggr.AddWrite(2, 100000)
	ggr.AddWrite(3, 100000)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: handlersQuantity,
		Inputs:           ggr.GetInputs(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	ggr.SetDiscipline(discipline)

	gauges := ggr.Play(context.Background())

	require.Equal(t, int(ggr.CalcExpectedGuagesQuantity()), len(research.FilterByKind(gauges, gauger.GaugeKindReceived)))
}

func TestDisciplineRateUnbuffered(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := gauger.Opts{
		HandlersQuantity: handlersQuantity,
		NoInputBuffer:    true,
	}

	ggr := gauger.New(gaugerOpts)

	ggr.AddWrite(1, 100000)
	ggr.AddWrite(2, 100000)
	ggr.AddWrite(3, 100000)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Rate,
		HandlersQuantity: handlersQuantity,
		Inputs:           ggr.GetInputs(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	ggr.SetDiscipline(discipline)

	gauges := ggr.Play(context.Background())

	require.Equal(t, int(ggr.CalcExpectedGuagesQuantity()), len(research.FilterByKind(gauges, gauger.GaugeKindReceived)))
}

func TestDisciplineFairUnbuffered(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := gauger.Opts{
		HandlersQuantity: handlersQuantity,
		NoInputBuffer:    true,
	}

	ggr := gauger.New(gaugerOpts)

	ggr.AddWrite(1, 100000)
	ggr.AddWrite(2, 100000)
	ggr.AddWrite(3, 100000)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: handlersQuantity,
		Inputs:           ggr.GetInputs(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	ggr.SetDiscipline(discipline)

	gauges := ggr.Play(context.Background())

	require.Equal(t, int(ggr.CalcExpectedGuagesQuantity()), len(research.FilterByKind(gauges, gauger.GaugeKindReceived)))
}

func TestDisciplineBadDivider(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := gauger.Opts{
		HandlersQuantity: handlersQuantity,
	}

	ggr := gauger.New(gaugerOpts)

	ggr.AddWrite(1, 100000)
	ggr.AddWrite(2, 100000)
	ggr.AddWrite(3, 100000)

	dividerCalled := 0

	divider := func(priorities []uint, dividend uint, distribution map[uint]uint) {
		divider.Fair(priorities, dividend, distribution)

		if dividerCalled < 1 {
			dividerCalled++
			return
		}

		for priority := range distribution {
			distribution[priority] *= 2
		}
	}

	disciplineOpts := Opts[uint]{
		Divider:          divider,
		HandlersQuantity: handlersQuantity,
		Inputs:           ggr.GetInputs(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	ggr.SetDiscipline(discipline)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if <-discipline.Err() != nil {
			cancel()
		}
	}()

	gauges := ggr.Play(ctx)

	require.NotEqual(
		t,
		int(ggr.CalcExpectedGuagesQuantity()),
		len(research.FilterByKind(gauges, gauger.GaugeKindReceived)),
	)
}

func TestDisciplineBadDividerInNew(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := gauger.Opts{
		HandlersQuantity: handlersQuantity,
	}

	ggr := gauger.New(gaugerOpts)

	ggr.AddWrite(1, 1)
	ggr.AddWrite(2, 1)
	ggr.AddWrite(3, 1)

	divider := func(priorities []uint, dividend uint, distribution map[uint]uint) {
		divider.Fair(priorities, dividend, distribution)

		for priority := range distribution {
			distribution[priority] *= 2
		}
	}

	disciplineOpts := Opts[uint]{
		Divider:          divider,
		HandlersQuantity: handlersQuantity,
		Inputs:           ggr.GetInputs(),
	}

	_, err := New(disciplineOpts)
	require.Error(t, err)
}

func TestDisciplineRateOverQuantity(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := gauger.Opts{
		HandlersQuantity: 2 * handlersQuantity,
	}

	gauger := gauger.New(gaugerOpts)

	gauger.AddWrite(1, 100000)
	gauger.AddWrite(2, 100000)
	gauger.AddWrite(3, 100000)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Rate,
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	gauger.SetDiscipline(discipline)

	gauges := gauger.Play(context.Background())

	quantities := research.CalcInProcessing(gauges, 100*time.Millisecond)

	for priority := range quantities {
		for id := range quantities[priority] {
			require.LessOrEqual(t, quantities[priority][id].Quantity, handlersQuantity)
		}
	}
}

func TestDisciplineFairOverQuantity(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := gauger.Opts{
		HandlersQuantity: 2 * handlersQuantity,
	}

	gauger := gauger.New(gaugerOpts)

	gauger.AddWrite(1, 1000000)
	gauger.AddWrite(2, 100000)
	gauger.AddWrite(3, 10000)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	gauger.SetDiscipline(discipline)

	gauges := gauger.Play(context.Background())

	quantities := research.CalcInProcessing(gauges, 100*time.Millisecond)

	for priority := range quantities {
		for id := range quantities[priority] {
			require.LessOrEqual(t, quantities[priority][id].Quantity, handlersQuantity)
		}
	}
}

func TestDisciplineRateFatalDividingError(t *testing.T) {
	handlersQuantity := uint(5)

	gaugerOpts := gauger.Opts{
		HandlersQuantity: handlersQuantity,
	}

	ggr := gauger.New(gaugerOpts)

	ggr.AddWrite(1, 100000)
	ggr.AddWrite(2, 100000)
	ggr.AddWrite(3, 100000)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Rate,
		HandlersQuantity: handlersQuantity,
		Inputs:           ggr.GetInputs(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	ggr.SetDiscipline(discipline)

	gauges := ggr.Play(context.Background())

	require.Equal(t, int(ggr.CalcExpectedGuagesQuantity()), len(research.FilterByKind(gauges, gauger.GaugeKindReceived)))
}

func TestDisciplineFairFatalDividingError(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := gauger.Opts{
		HandlersQuantity: handlersQuantity,
	}

	ggr := gauger.New(gaugerOpts)

	ggr.AddWrite(1, 100000)
	ggr.AddWrite(2, 100000)
	ggr.AddWrite(3, 100000)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: handlersQuantity,
		Inputs:           ggr.GetInputs(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	ggr.SetDiscipline(discipline)

	gauges := ggr.Play(context.Background())

	require.Equal(t, int(ggr.CalcExpectedGuagesQuantity()), len(research.FilterByKind(gauges, gauger.GaugeKindReceived)))
}

func testDisciplineFairEvenProcessingTimeDividingError(t *testing.T, handlersQuantity uint) {
	if os.Getenv("CQOS_ENABLE_GRAPHS") == "" {
		t.SkipNow()
	}

	gaugerOpts := gauger.Opts{
		HandlersQuantity: handlersQuantity,
	}

	ggr := gauger.New(gaugerOpts)

	ggr.AddWrite(1, 4000)

	ggr.AddWrite(2, 500)
	ggr.AddWaitDevastation(2)
	ggr.AddDelay(2, 2*time.Second)
	ggr.AddWrite(2, 500)
	ggr.AddWaitDevastation(2)
	ggr.AddDelay(2, 4*time.Second)
	ggr.AddWrite(2, 1000)
	ggr.AddWaitDevastation(2)
	ggr.AddDelay(2, 2*time.Second)
	ggr.AddWrite(2, 2000)

	ggr.AddWrite(3, 500)
	ggr.AddWaitDevastation(3)
	ggr.AddDelay(3, 5*time.Second)
	ggr.AddWrite(3, 4000)

	ggr.AddWrite(4, 500)
	ggr.AddWaitDevastation(3)
	ggr.AddDelay(4, 5*time.Second)
	ggr.AddWrite(4, 4000)

	ggr.SetProcessDelay(1, 10*time.Millisecond)
	ggr.SetProcessDelay(2, 10*time.Millisecond)
	ggr.SetProcessDelay(3, 10*time.Millisecond)
	ggr.SetProcessDelay(4, 10*time.Millisecond)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: handlersQuantity,
		Inputs:           ggr.GetInputs(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	ggr.SetDiscipline(discipline)

	gauges := ggr.Play(context.Background())

	received := research.FilterByKind(gauges, gauger.GaugeKindReceived)

	dqot, dqotX := research.ConvertToLineEcharts(
		research.CalcDataQuantity(received, 100*time.Millisecond),
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

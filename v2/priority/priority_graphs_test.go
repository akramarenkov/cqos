package priority

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/priority/divider"
	"github.com/akramarenkov/cqos/v2/priority/internal/measurer"
	"github.com/akramarenkov/cqos/v2/priority/internal/research"
	"github.com/akramarenkov/cqos/v2/priority/internal/unmanaged"

	"github.com/go-echarts/go-echarts/v2/charts"
	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
	"github.com/stretchr/testify/require"
)

func testDisciplineRateEvenProcessingTime(
	t *testing.T,
	factor uint,
	unbufferedInput bool,
) {
	if os.Getenv("CQOS_ENABLE_GRAPHS") == "" {
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

	dqotChart := charts.NewLine()
	ipotChart := charts.NewLine()
	wtflChart := charts.NewBar()

	subtitle := fmt.Sprintf(
		"Rate divider, even time processing, "+
			"handlers quantity: %d, "+
			"buffered: %t, "+
			"time: %s",
		measurerOpts.HandlersQuantity,
		!unbufferedInput,
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

	baseName := "graph_rate_even_" +
		strconv.Itoa(int(measurerOpts.HandlersQuantity)) +
		"_buffered_" +
		strconv.FormatBool(!unbufferedInput)

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

func testDisciplineRateUnevenProcessingTime(
	t *testing.T,
	factor uint,
	unbufferedInput bool,
) {
	if os.Getenv("CQOS_ENABLE_GRAPHS") == "" {
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

	dqotChart := charts.NewLine()
	ipotChart := charts.NewLine()
	wtflChart := charts.NewBar()

	subtitle := fmt.Sprintf(
		"Rate divider, uneven time processing, "+
			"handlers quantity: %d, "+
			"buffered: %t, "+
			"time: %s",
		measurerOpts.HandlersQuantity,
		!unbufferedInput,
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

	baseName := "graph_rate_uneven_" +
		strconv.Itoa(int(measurerOpts.HandlersQuantity)) +
		"_buffered_" +
		strconv.FormatBool(!unbufferedInput)

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

func testDisciplineFairEvenProcessingTime(
	t *testing.T,
	factor uint,
	unbufferedInput bool,
) {
	if os.Getenv("CQOS_ENABLE_GRAPHS") == "" {
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

	dqotChart := charts.NewLine()
	ipotChart := charts.NewLine()
	wtflChart := charts.NewBar()

	subtitle := fmt.Sprintf(
		"Fair divider, even time processing, "+
			"handlers quantity: %d, "+
			"buffered: %t, "+
			"time: %s",
		measurerOpts.HandlersQuantity,
		!unbufferedInput,
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

	baseName := "graph_fair_even_" +
		strconv.Itoa(int(measurerOpts.HandlersQuantity)) +
		"_buffered_" +
		strconv.FormatBool(!unbufferedInput)

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

func testDisciplineFairUnevenProcessingTime(
	t *testing.T,
	factor uint,
	unbufferedInput bool,
) {
	if os.Getenv("CQOS_ENABLE_GRAPHS") == "" {
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

	dqotChart := charts.NewLine()
	ipotChart := charts.NewLine()
	wtflChart := charts.NewBar()

	subtitle := fmt.Sprintf(
		"Fair divider, uneven time processing, "+
			"handlers quantity: %d, "+
			"buffered: %t, "+
			"time: %s",
		measurerOpts.HandlersQuantity,
		!unbufferedInput,
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

	baseName := "graph_fair_uneven_" +
		strconv.Itoa(int(measurerOpts.HandlersQuantity)) +
		"_buffered_" +
		strconv.FormatBool(!unbufferedInput)

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

func testUnmanagedEven(t *testing.T, factor uint, unbufferedInput bool) {
	if os.Getenv("CQOS_ENABLE_GRAPHS") == "" {
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

	dqotChart := charts.NewLine()
	ipotChart := charts.NewLine()

	subtitle := fmt.Sprintf(
		"Unmanaged, even time processing, "+
			"handlers quantity: %d, "+
			"buffered: %t, "+
			"time: %s",
		measurerOpts.HandlersQuantity,
		!unbufferedInput,
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

	baseName := "graph_unmanaged_even_" +
		strconv.Itoa(int(measurerOpts.HandlersQuantity)) +
		"_buffered_" +
		strconv.FormatBool(!unbufferedInput)

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

func testUnmanagedUneven(t *testing.T, factor uint, unbufferedInput bool) {
	if os.Getenv("CQOS_ENABLE_GRAPHS") == "" {
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

	dqotChart := charts.NewLine()
	ipotChart := charts.NewLine()

	subtitle := fmt.Sprintf(
		"Unmanaged, uneven time processing, "+
			"handlers quantity: %d, "+
			"buffered: %t, "+
			"time: %s",
		measurerOpts.HandlersQuantity,
		!unbufferedInput,
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

	baseName := "graph_unmanaged_uneven_" +
		strconv.Itoa(int(measurerOpts.HandlersQuantity)) +
		"_buffered_" +
		strconv.FormatBool(!unbufferedInput)

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
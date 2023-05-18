package priority

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/akramarenkov/cqos/priority/divider"
	"github.com/akramarenkov/cqos/priority/test"

	"github.com/go-echarts/go-echarts/v2/charts"
	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
	"github.com/stretchr/testify/require"
)

func TestDisciplineRateEvenProcessingTime(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := test.GaugerOpts{
		HandlersQuantity: handlersQuantity,
	}

	gauger := test.NewGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 4100)

	gauger.AddWrite(2, 1500)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 750)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 4*time.Second)
	gauger.AddWrite(2, 700)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 3*time.Second)
	gauger.AddWrite(2, 1200)

	gauger.AddWrite(3, 1000)
	gauger.AddWaitDevastation(3)
	gauger.AddDelay(3, 8*time.Second)
	gauger.AddWrite(3, 3700)

	gauger.SetProcessDelay(1, 10*time.Millisecond)
	gauger.SetProcessDelay(2, 10*time.Millisecond)
	gauger.SetProcessDelay(3, 10*time.Millisecond)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Rate,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	gauges := gauger.Play()

	received := test.FilterByKind(gauges, test.GaugeKindReceived)

	dqot, dqotX := test.ConvertQuantityOverTimeToLineEcharts(
		test.CalcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := test.ConvertQuantityOverTimeToLineEcharts(
		test.CalcInProcessing(gauges, 100*time.Millisecond),
		1*time.Second,
	)

	wtfl, wtflX := test.ConvertQuantityOverTimeToBarEcharts(
		test.CalcWriteToFeedbackLatency(gauges, 100*time.Nanosecond),
	)

	dqotChart := charts.NewLine()
	ipotChart := charts.NewLine()
	wtflChart := charts.NewBar()

	subtitle := fmt.Sprintf(
		"Rate divider, even time processing, handlers quantity: %d, time: %s",
		handlersQuantity,
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

	dqotFile, err := os.Create("graph_rate_even_data_retrieval.html")
	require.NoError(t, err)

	err = dqotChart.Render(dqotFile)
	require.NoError(t, err)

	ipotFile, err := os.Create("graph_rate_even_in_processing.html")
	require.NoError(t, err)

	err = ipotChart.Render(ipotFile)
	require.NoError(t, err)

	wtflFile, err := os.Create("graph_rate_even_write_feedback_latency.html")
	require.NoError(t, err)

	err = wtflChart.Render(wtflFile)
	require.NoError(t, err)
}

func TestDisciplineRateUnevenProcessingTime(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := test.GaugerOpts{
		HandlersQuantity: handlersQuantity,
	}

	gauger := test.NewGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 430)

	gauger.AddWrite(2, 250)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 100)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 4*time.Second)
	gauger.AddWrite(2, 150)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 300)

	gauger.AddWrite(3, 1000)
	gauger.AddWaitDevastation(3)
	gauger.AddDelay(3, 8*time.Second)
	gauger.AddWrite(3, 3500)

	gauger.SetProcessDelay(1, 100*time.Millisecond)
	gauger.SetProcessDelay(2, 50*time.Millisecond)
	gauger.SetProcessDelay(3, 10*time.Millisecond)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Rate,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	gauges := gauger.Play()

	received := test.FilterByKind(gauges, test.GaugeKindReceived)

	dqot, dqotX := test.ConvertQuantityOverTimeToLineEcharts(
		test.CalcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := test.ConvertQuantityOverTimeToLineEcharts(
		test.CalcInProcessing(gauges, 100*time.Millisecond),
		1*time.Second,
	)

	wtfl, wtflX := test.ConvertQuantityOverTimeToBarEcharts(
		test.CalcWriteToFeedbackLatency(gauges, 100*time.Nanosecond),
	)

	dqotChart := charts.NewLine()
	ipotChart := charts.NewLine()
	wtflChart := charts.NewBar()

	subtitle := fmt.Sprintf(
		"Rate divider, uneven time processing, handlers quantity: %d, time: %s",
		handlersQuantity,
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

	dqotFile, err := os.Create("graph_rate_uneven_data_retrieval.html")
	require.NoError(t, err)

	err = dqotChart.Render(dqotFile)
	require.NoError(t, err)

	ipotFile, err := os.Create("graph_rate_uneven_in_processing.html")
	require.NoError(t, err)

	err = ipotChart.Render(ipotFile)
	require.NoError(t, err)

	wtflFile, err := os.Create("graph_rate_uneven_write_feedback_latency.html")
	require.NoError(t, err)

	err = wtflChart.Render(wtflFile)
	require.NoError(t, err)
}

func TestDisciplineFairEvenProcessingTime(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := test.GaugerOpts{
		HandlersQuantity: handlersQuantity,
	}

	gauger := test.NewGauger(gaugerOpts)
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

	gauger.SetProcessDelay(1, 10*time.Millisecond)
	gauger.SetProcessDelay(2, 10*time.Millisecond)
	gauger.SetProcessDelay(3, 10*time.Millisecond)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Fair,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	gauges := gauger.Play()

	received := test.FilterByKind(gauges, test.GaugeKindReceived)

	dqot, dqotX := test.ConvertQuantityOverTimeToLineEcharts(
		test.CalcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := test.ConvertQuantityOverTimeToLineEcharts(
		test.CalcInProcessing(gauges, 100*time.Millisecond),
		1*time.Second,
	)

	wtfl, wtflX := test.ConvertQuantityOverTimeToBarEcharts(
		test.CalcWriteToFeedbackLatency(gauges, 100*time.Nanosecond),
	)

	dqotChart := charts.NewLine()
	ipotChart := charts.NewLine()
	wtflChart := charts.NewBar()

	subtitle := fmt.Sprintf(
		"Fair divider, even time processing, handlers quantity: %d, time: %s",
		handlersQuantity,
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

	dqotFile, err := os.Create("graph_fair_even_data_retrieval.html")
	require.NoError(t, err)

	err = dqotChart.Render(dqotFile)
	require.NoError(t, err)

	ipotFile, err := os.Create("graph_fair_even_in_processing.html")
	require.NoError(t, err)

	err = ipotChart.Render(ipotFile)
	require.NoError(t, err)

	wtflFile, err := os.Create("graph_fair_even_write_feedback_latency.html")
	require.NoError(t, err)

	err = wtflChart.Render(wtflFile)
	require.NoError(t, err)
}

func TestDisciplineFairUnevenProcessingTime(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := test.GaugerOpts{
		HandlersQuantity: handlersQuantity,
	}

	gauger := test.NewGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 450)

	gauger.AddWrite(2, 100)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 100)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 4*time.Second)
	gauger.AddWrite(2, 200)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 400)

	gauger.AddWrite(3, 500)
	gauger.AddWaitDevastation(3)
	gauger.AddDelay(3, 6*time.Second)
	gauger.AddWrite(3, 3000)

	gauger.SetProcessDelay(1, 100*time.Millisecond)
	gauger.SetProcessDelay(2, 50*time.Millisecond)
	gauger.SetProcessDelay(3, 10*time.Millisecond)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Fair,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	gauges := gauger.Play()

	received := test.FilterByKind(gauges, test.GaugeKindReceived)

	dqot, dqotX := test.ConvertQuantityOverTimeToLineEcharts(
		test.CalcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := test.ConvertQuantityOverTimeToLineEcharts(
		test.CalcInProcessing(gauges, 100*time.Millisecond),
		1*time.Second,
	)

	wtfl, wtflX := test.ConvertQuantityOverTimeToBarEcharts(
		test.CalcWriteToFeedbackLatency(gauges, 100*time.Nanosecond),
	)

	dqotChart := charts.NewLine()
	ipotChart := charts.NewLine()
	wtflChart := charts.NewBar()

	subtitle := fmt.Sprintf(
		"Fair divider, uneven time processing, handlers quantity: %d, time: %s",
		handlersQuantity,
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

	dqotFile, err := os.Create("graph_fair_uneven_data_retrieval.html")
	require.NoError(t, err)

	err = dqotChart.Render(dqotFile)
	require.NoError(t, err)

	ipotFile, err := os.Create("graph_fair_uneven_in_processing.html")
	require.NoError(t, err)

	err = ipotChart.Render(ipotFile)
	require.NoError(t, err)

	wtflFile, err := os.Create("graph_fair_uneven_write_feedback_latency.html")
	require.NoError(t, err)

	err = wtflChart.Render(wtflFile)
	require.NoError(t, err)
}

func TestUnmanagedEven(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := test.GaugerOpts{
		HandlersQuantity: handlersQuantity,
		NoFeedback:       true,
	}

	gauger := test.NewGauger(gaugerOpts)
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

	gauger.SetProcessDelay(1, 10*time.Millisecond)
	gauger.SetProcessDelay(2, 10*time.Millisecond)
	gauger.SetProcessDelay(3, 10*time.Millisecond)

	unmanagedOpts := test.UnmanagedOpts[uint]{
		Inputs: gauger.GetInputs(),
		Output: gauger.GetOutput(),
	}

	unmanaged, err := test.NewUnmanaged(unmanagedOpts)
	require.NoError(t, err)

	defer unmanaged.Stop()

	gauges := gauger.Play()

	received := test.FilterByKind(gauges, test.GaugeKindReceived)

	dqot, dqotX := test.ConvertQuantityOverTimeToLineEcharts(
		test.CalcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := test.ConvertQuantityOverTimeToLineEcharts(
		test.CalcInProcessing(gauges, 100*time.Millisecond),
		1*time.Second,
	)

	dqotChart := charts.NewLine()
	ipotChart := charts.NewLine()

	subtitle := fmt.Sprintf(
		"Unmanaged, even time processing, handlers quantity: %d, time: %s",
		handlersQuantity,
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

	dqotFile, err := os.Create("graph_unmanaged_even_data_retrieval.html")
	require.NoError(t, err)

	err = dqotChart.Render(dqotFile)
	require.NoError(t, err)

	ipotFile, err := os.Create("graph_unmanaged_even_in_processing.html")
	require.NoError(t, err)

	err = ipotChart.Render(ipotFile)
	require.NoError(t, err)
}

func TestUnmanagedUneven(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := test.GaugerOpts{
		HandlersQuantity: handlersQuantity,
		NoFeedback:       true,
	}

	gauger := test.NewGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 500)

	gauger.AddWrite(2, 100)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 100)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 200)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 400)

	gauger.AddWrite(3, 100)
	gauger.AddWaitDevastation(3)
	gauger.AddDelay(3, 6*time.Second)
	gauger.AddWrite(3, 1350)

	gauger.SetProcessDelay(1, 100*time.Millisecond)
	gauger.SetProcessDelay(2, 50*time.Millisecond)
	gauger.SetProcessDelay(3, 10*time.Millisecond)

	unmanagedOpts := test.UnmanagedOpts[uint]{
		Inputs: gauger.GetInputs(),
		Output: gauger.GetOutput(),
	}

	unmanaged, err := test.NewUnmanaged(unmanagedOpts)
	require.NoError(t, err)

	defer unmanaged.Stop()

	gauges := gauger.Play()

	received := test.FilterByKind(gauges, test.GaugeKindReceived)

	dqot, dqotX := test.ConvertQuantityOverTimeToLineEcharts(
		test.CalcDataQuantity(received, 100*time.Millisecond),
		1*time.Second,
	)

	ipot, ipotX := test.ConvertQuantityOverTimeToLineEcharts(
		test.CalcInProcessing(gauges, 100*time.Millisecond),
		1*time.Second,
	)

	dqotChart := charts.NewLine()
	ipotChart := charts.NewLine()

	subtitle := fmt.Sprintf(
		"Unmanaged, uneven time processing, handlers quantity: %d, time: %s",
		handlersQuantity,
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

	dqotFile, err := os.Create("graph_unmanaged_uneven_data_retrieval.html")
	require.NoError(t, err)

	err = dqotChart.Render(dqotFile)
	require.NoError(t, err)

	ipotFile, err := os.Create("graph_unmanaged_uneven_in_processing.html")
	require.NoError(t, err)

	err = ipotChart.Render(ipotFile)
	require.NoError(t, err)
}

func BenchmarkDisciplineFair(b *testing.B) {
	handlersQuantity := uint(600)

	gaugerOpts := test.GaugerOpts{
		DisableGauges:    true,
		HandlersQuantity: handlersQuantity,
	}

	gauger := test.NewGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 5000000)
	gauger.AddWrite(2, 5000000)
	gauger.AddWrite(3, 5000000)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Fair,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(b, err)

	defer discipline.Stop()

	_ = gauger.Play()
}

func TestOverQuantity(t *testing.T) {
	handlersQuantity := uint(6)

	gaugerOpts := test.GaugerOpts{
		HandlersQuantity: 2 * handlersQuantity,
	}

	gauger := test.NewGauger(gaugerOpts)
	defer gauger.Finalize()

	gauger.AddWrite(1, 500000)
	gauger.AddWaitDevastation(1)
	gauger.AddDelay(1, 1*time.Second)
	gauger.AddWrite(1, 500000)
	gauger.AddWaitDevastation(1)
	gauger.AddDelay(1, 1*time.Second)
	gauger.AddWrite(1, 500000)

	gauger.AddWrite(2, 500000)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 1*time.Second)
	gauger.AddWrite(2, 500000)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 1*time.Second)
	gauger.AddWrite(2, 500000)

	gauger.AddWrite(3, 500000)
	gauger.AddWaitDevastation(3)
	gauger.AddDelay(3, 1*time.Second)
	gauger.AddWrite(3, 500000)
	gauger.AddWaitDevastation(3)
	gauger.AddDelay(3, 1*time.Second)
	gauger.AddWrite(3, 500000)

	disciplineOpts := Opts[uint]{
		Divider:          divider.Rate,
		Feedback:         gauger.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           gauger.GetInputs(),
		Output:           gauger.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	gauges := gauger.Play()

	quantities := test.CalcInProcessing(gauges, 100*time.Millisecond)

	for priority := range quantities {
		for id := range quantities[priority] {
			require.LessOrEqual(t, quantities[priority][id].Quantity, handlersQuantity)
		}
	}
}

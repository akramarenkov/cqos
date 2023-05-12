package priority

import (
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
	gauger.AddDelay(2, 3*time.Second)
	gauger.AddWrite(2, 700)
	gauger.AddWaitDevastation(2)
	gauger.AddDelay(2, 2*time.Second)
	gauger.AddWrite(2, 1000)

	gauger.AddWrite(3, 1000)
	gauger.AddWaitDevastation(3)
	gauger.AddDelay(3, 9*time.Second)
	gauger.AddWrite(3, 3500)

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

	dqot, dqotX := test.CalcDataQuantityOverTime(received, 100*time.Millisecond, 1*time.Second)
	wfl, wflX := test.CalcWriteToFeedbackLatency(gauges, 100*time.Nanosecond)

	dqotChart := charts.NewLine()
	wflChart := charts.NewBar()

	dqotChart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "Data retrieval graph",
				Subtitle: time.Now().Format(time.RFC3339),
			},
		),
	)

	wflChart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "Write to feedback latency",
				Subtitle: time.Now().Format(time.RFC3339),
			},
		),
	)

	dqotChart.SetXAxis(dqotX).
		AddSeries("3", dqot[3]).
		AddSeries("2", dqot[2]).
		AddSeries("1", dqot[1])

	dqotFile, err := os.Create("graph_rate_even_data_retrieval.html")
	require.NoError(t, err)

	dqotChart.Render(dqotFile)

	wflChart.SetXAxis(wflX).
		AddSeries("3", wfl[3]).
		AddSeries("2", wfl[2]).
		AddSeries("1", wfl[1])

	wflFile, err := os.Create("graph_rate_even_write_feedback_latency.html")
	require.NoError(t, err)

	wflChart.Render(wflFile)
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

	dqot, dqotX := test.CalcDataQuantityOverTime(received, 100*time.Millisecond, 1*time.Second)
	wfl, wflX := test.CalcWriteToFeedbackLatency(gauges, 100*time.Nanosecond)

	dqotChart := charts.NewLine()
	wflChart := charts.NewBar()

	dqotChart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "Data retrieval graph",
				Subtitle: time.Now().Format(time.RFC3339),
			},
		),
	)

	wflChart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "Write to feedback latency",
				Subtitle: time.Now().Format(time.RFC3339),
			},
		),
	)

	dqotChart.SetXAxis(dqotX).
		AddSeries("3", dqot[3]).
		AddSeries("2", dqot[2]).
		AddSeries("1", dqot[1])

	dqotFile, err := os.Create("graph_fair_even_data_retrieval.html")
	require.NoError(t, err)

	dqotChart.Render(dqotFile)

	wflChart.SetXAxis(wflX).
		AddSeries("3", wfl[3]).
		AddSeries("2", wfl[2]).
		AddSeries("1", wfl[1])

	wflFile, err := os.Create("graph_fair_even_write_feedback_latency.html")
	require.NoError(t, err)

	wflChart.Render(wflFile)
}

func TestUnmanaged(t *testing.T) {
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

	discipline, err := test.NewUnmanaged(unmanagedOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	gauges := gauger.Play()

	received := test.FilterByKind(gauges, test.GaugeKindReceived)

	dqot, dqotX := test.CalcDataQuantityOverTime(received, 100*time.Millisecond, 1*time.Second)

	dqotChart := charts.NewLine()

	dqotChart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    "Data retrieval graph",
				Subtitle: time.Now().Format(time.RFC3339),
			},
		),
	)

	dqotChart.SetXAxis(dqotX).
		AddSeries("3", dqot[3]).
		AddSeries("2", dqot[2]).
		AddSeries("1", dqot[1])

	dqotFile, err := os.Create("graph_unmanaged_data_retrieval.html")
	require.NoError(t, err)

	dqotChart.Render(dqotFile)
}

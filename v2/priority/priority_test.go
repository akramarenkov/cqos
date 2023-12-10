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

	"github.com/go-echarts/go-echarts/v2/charts"
	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
	"github.com/stretchr/testify/require"
)

func TestOptsValidation(t *testing.T) {
	opts := Opts[uint]{
		HandlersQuantity: 6,
	}

	_, err := New(opts)
	require.Error(t, err)

	opts = Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: 6,
	}

	_, err = New(opts)
	require.Error(t, err)

	opts = Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: 6,
		Inputs: map[uint]<-chan uint{
			1: make(chan uint),
		},
	}

	_, err = New(opts)
	require.NoError(t, err)
}

func BenchmarkDisciplineFair(b *testing.B) {
	measurerOpts := measurer.Opts{
		DisableMeasures:  true,
		HandlersQuantity: 600,
	}

	measurer := measurer.New(measurerOpts)

	measurer.AddWrite(1, 5000000)
	measurer.AddWrite(2, 5000000)
	measurer.AddWrite(3, 5000000)

	opts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           measurer.GetInputs(),
	}

	discipline, err := New(opts)
	require.NoError(b, err)

	_ = measurer.Play(discipline)
}

func BenchmarkDisciplineRate(b *testing.B) {
	measurerOpts := measurer.Opts{
		DisableMeasures:  true,
		HandlersQuantity: 600,
	}

	measurer := measurer.New(measurerOpts)

	measurer.AddWrite(1, 5000000)
	measurer.AddWrite(2, 5000000)
	measurer.AddWrite(3, 5000000)

	opts := Opts[uint]{
		Divider:          divider.Rate,
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           measurer.GetInputs(),
	}

	discipline, err := New(opts)
	require.NoError(b, err)

	_ = measurer.Play(discipline)
}

func BenchmarkDisciplineFairUnbuffered(b *testing.B) {
	measurerOpts := measurer.Opts{
		DisableMeasures:  true,
		HandlersQuantity: 600,
		UnbufferedInput:  true,
	}

	measurer := measurer.New(measurerOpts)

	measurer.AddWrite(1, 5000000)
	measurer.AddWrite(2, 5000000)
	measurer.AddWrite(3, 5000000)

	opts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           measurer.GetInputs(),
	}

	discipline, err := New(opts)
	require.NoError(b, err)

	_ = measurer.Play(discipline)
}

func BenchmarkDisciplineRateUnbuffered(b *testing.B) {
	measurerOpts := measurer.Opts{
		DisableMeasures:  true,
		HandlersQuantity: 600,
		UnbufferedInput:  true,
	}

	measurer := measurer.New(measurerOpts)

	measurer.AddWrite(1, 5000000)
	measurer.AddWrite(2, 5000000)
	measurer.AddWrite(3, 5000000)

	opts := Opts[uint]{
		Divider:          divider.Rate,
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           measurer.GetInputs(),
	}

	discipline, err := New(opts)
	require.NoError(b, err)

	_ = measurer.Play(discipline)
}

func TestDisciplineRate(t *testing.T) {
	measurerOpts := measurer.Opts{
		HandlersQuantity: 6,
	}

	msr := measurer.New(measurerOpts)

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

	opts := Opts[uint]{
		Divider:          divider.Rate,
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measures := msr.Play(discipline)

	require.Equal(t, int(msr.GetExpectedMeasuresQuantity()), len(measures))
}

func TestDisciplineFair(t *testing.T) {
	measurerOpts := measurer.Opts{
		HandlersQuantity: 6,
	}

	msr := measurer.New(measurerOpts)

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

	opts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measures := msr.Play(discipline)

	require.Equal(t, int(msr.GetExpectedMeasuresQuantity()), len(measures))
}

func TestDisciplineRateUnbuffered(t *testing.T) {
	measurerOpts := measurer.Opts{
		HandlersQuantity: 6,
		UnbufferedInput:  true,
	}

	msr := measurer.New(measurerOpts)

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

	opts := Opts[uint]{
		Divider:          divider.Rate,
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measures := msr.Play(discipline)

	require.Equal(t, int(msr.GetExpectedMeasuresQuantity()), len(measures))
}

func TestDisciplineFairUnbuffered(t *testing.T) {
	measurerOpts := measurer.Opts{
		HandlersQuantity: 6,
		UnbufferedInput:  true,
	}

	msr := measurer.New(measurerOpts)

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

	opts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measures := msr.Play(discipline)

	require.Equal(t, int(msr.GetExpectedMeasuresQuantity()), len(measures))
}

func TestDisciplineBadDivider(t *testing.T) {
	measurerOpts := measurer.Opts{
		HandlersQuantity: 6,
	}

	msr := measurer.New(measurerOpts)

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

	dividerCallsQuantity := 0

	divider := func(priorities []uint, dividend uint, distribution map[uint]uint) {
		divider.Fair(priorities, dividend, distribution)

		dividerCallsQuantity++

		if dividerCallsQuantity == 1 {
			return
		}

		for priority := range distribution {
			distribution[priority] *= 2
		}
	}

	opts := Opts[uint]{
		Divider:          divider,
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measures := msr.Play(discipline)

	require.NotEqual(t, int(msr.GetExpectedMeasuresQuantity()), len(measures))
}

func TestDisciplineBadDividerInRecalc(t *testing.T) {
	measurerOpts := measurer.Opts{
		HandlersQuantity: 6,
	}

	msr := measurer.New(measurerOpts)

	msr.AddWrite(1, 0)
	msr.AddWrite(2, 0)
	msr.AddWrite(3, 100000)

	dividerCallsQuantity := 0

	divider := func(priorities []uint, dividend uint, distribution map[uint]uint) {
		divider.Fair(priorities, dividend, distribution)

		dividerCallsQuantity++

		if dividerCallsQuantity == 1 || dividerCallsQuantity%2 == 0 {
			return
		}

		for priority := range distribution {
			distribution[priority] *= 2
		}
	}

	opts := Opts[uint]{
		Divider:          divider,
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measures := msr.Play(discipline)

	require.NotEqual(t, int(msr.GetExpectedMeasuresQuantity()), len(measures))
}

func TestDisciplineBadDividerInNew(t *testing.T) {
	measurerOpts := measurer.Opts{
		HandlersQuantity: 6,
	}

	msr := measurer.New(measurerOpts)

	msr.AddWrite(1, 1)
	msr.AddWrite(2, 1)
	msr.AddWrite(3, 1)

	divider := func(priorities []uint, dividend uint, distribution map[uint]uint) {
		divider.Fair(priorities, dividend, distribution)

		for priority := range distribution {
			distribution[priority] *= 2
		}
	}

	opts := Opts[uint]{
		Divider:          divider,
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
	}

	_, err := New(opts)
	require.Error(t, err)
}

func TestDisciplineRateOverQuantity(t *testing.T) {
	handlersQuantity := uint(6)

	measurerOpts := measurer.Opts{
		HandlersQuantity: 2 * handlersQuantity,
	}

	measurer := measurer.New(measurerOpts)

	measurer.AddWrite(1, 100000)
	measurer.AddWrite(2, 100000)
	measurer.AddWrite(3, 100000)

	opts := Opts[uint]{
		Divider:          divider.Rate,
		HandlersQuantity: handlersQuantity,
		Inputs:           measurer.GetInputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measures := measurer.Play(discipline)

	quantities := research.CalcInProcessing(measures, 100*time.Millisecond)

	for priority := range quantities {
		for id := range quantities[priority] {
			require.LessOrEqual(t, quantities[priority][id].Quantity, handlersQuantity)
		}
	}
}

func TestDisciplineFairOverQuantity(t *testing.T) {
	handlersQuantity := uint(6)

	measurerOpts := measurer.Opts{
		HandlersQuantity: 2 * handlersQuantity,
	}

	measurer := measurer.New(measurerOpts)

	measurer.AddWrite(1, 1000000)
	measurer.AddWrite(2, 100000)
	measurer.AddWrite(3, 10000)

	opts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: handlersQuantity,
		Inputs:           measurer.GetInputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measures := measurer.Play(discipline)

	quantities := research.CalcInProcessing(measures, 100*time.Millisecond)

	for priority := range quantities {
		for id := range quantities[priority] {
			require.LessOrEqual(t, quantities[priority][id].Quantity, handlersQuantity)
		}
	}
}

func TestDisciplineRateFatalDividingError(t *testing.T) {
	measurerOpts := measurer.Opts{
		HandlersQuantity: 5,
	}

	msr := measurer.New(measurerOpts)

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

	opts := Opts[uint]{
		Divider:          divider.Rate,
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measures := msr.Play(discipline)

	require.Equal(t, int(msr.GetExpectedMeasuresQuantity()), len(measures))
}

func TestDisciplineFairFatalDividingError(t *testing.T) {
	measurerOpts := measurer.Opts{
		HandlersQuantity: 2,
	}

	msr := measurer.New(measurerOpts)

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

	opts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measures := msr.Play(discipline)

	require.Equal(t, int(msr.GetExpectedMeasuresQuantity()), len(measures))
}

func testDisciplineFairEvenProcessingTimeDividingError(t *testing.T, handlersQuantity uint) {
	if os.Getenv("CQOS_ENABLE_GRAPHS") == "" {
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

	dqotChart := charts.NewLine()

	subtitle := fmt.Sprintf(
		"Fair divider, even time processing, "+
			"significant dividing error, "+
			"handlers quantity: %d, "+
			"buffered: %t, "+
			"time: %s",
		measurerOpts.HandlersQuantity,
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

	baseName := "graph_fair_even_" +
		strconv.Itoa(int(measurerOpts.HandlersQuantity)) +
		"_buffered_" +
		strconv.FormatBool(!measurerOpts.UnbufferedInput) + "_dividing_error"

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

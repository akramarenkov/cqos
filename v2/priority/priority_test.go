package priority

import (
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/priority/divider"
	"github.com/akramarenkov/cqos/v2/priority/internal/measurer"
	"github.com/akramarenkov/cqos/v2/priority/internal/research"

	"github.com/stretchr/testify/require"
)

func TestOptsValidation(t *testing.T) {
	opts := Opts[uint]{}

	_, err := New(opts)
	require.Error(t, err)

	opts = Opts[uint]{
		Divider: divider.Fair,
	}

	_, err = New(opts)
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

func TestDisciplineFairTooSmallHandlersQuantity(t *testing.T) {
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

	_, err := New(opts)
	require.Error(t, err)
}

func TestDisciplineRateTooSmallHandlersQuantity(t *testing.T) {
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

	_, err := New(opts)
	require.Error(t, err)
}

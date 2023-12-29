package priority

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOptsValidation(t *testing.T) {
	handlersQuantity := uint(6)

	opts := Opts[uint]{
		Feedback:         make(chan uint),
		HandlersQuantity: handlersQuantity,
		Output:           make(chan Prioritized[uint]),
	}

	_, err := New(opts)
	require.Error(t, err)

	opts = Opts[uint]{
		Divider:          FairDivider,
		HandlersQuantity: handlersQuantity,
		Output:           make(chan Prioritized[uint]),
	}

	_, err = New(opts)
	require.Error(t, err)

	opts = Opts[uint]{
		Divider:          FairDivider,
		Feedback:         make(chan uint),
		HandlersQuantity: handlersQuantity,
	}

	_, err = New(opts)
	require.Error(t, err)
}

func BenchmarkDisciplineFair(b *testing.B) {
	measurerOpts := measurerOpts{
		DisableMeasures:  true,
		HandlersQuantity: 600,
	}

	msr := newMeasurer(measurerOpts)
	defer msr.Finalize()

	msr.AddWrite(1, 5000000)
	msr.AddWrite(2, 5000000)
	msr.AddWrite(3, 5000000)

	opts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         msr.GetFeedback(),
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
		Output:           msr.GetOutput(),
	}

	discipline, err := New(opts)
	require.NoError(b, err)

	defer discipline.Stop()

	_ = msr.Play(discipline)
}

func BenchmarkDisciplineRate(b *testing.B) {
	measurerOpts := measurerOpts{
		DisableMeasures:  true,
		HandlersQuantity: 600,
	}

	msr := newMeasurer(measurerOpts)
	defer msr.Finalize()

	msr.AddWrite(1, 5000000)
	msr.AddWrite(2, 5000000)
	msr.AddWrite(3, 5000000)

	opts := Opts[uint]{
		Divider:          RateDivider,
		Feedback:         msr.GetFeedback(),
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
		Output:           msr.GetOutput(),
	}

	discipline, err := New(opts)
	require.NoError(b, err)

	defer discipline.Stop()

	_ = msr.Play(discipline)
}

func BenchmarkDisciplineFairUnbuffered(b *testing.B) {
	measurerOpts := measurerOpts{
		DisableMeasures:  true,
		HandlersQuantity: 600,
		UnbufferedInput:  true,
	}

	msr := newMeasurer(measurerOpts)
	defer msr.Finalize()

	msr.AddWrite(1, 5000000)
	msr.AddWrite(2, 5000000)
	msr.AddWrite(3, 5000000)

	opts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         msr.GetFeedback(),
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
		Output:           msr.GetOutput(),
	}

	discipline, err := New(opts)
	require.NoError(b, err)

	defer discipline.Stop()

	_ = msr.Play(discipline)
}

func BenchmarkDisciplineRateUnbuffered(b *testing.B) {
	measurerOpts := measurerOpts{
		DisableMeasures:  true,
		HandlersQuantity: 600,
		UnbufferedInput:  true,
	}

	msr := newMeasurer(measurerOpts)
	defer msr.Finalize()

	msr.AddWrite(1, 5000000)
	msr.AddWrite(2, 5000000)
	msr.AddWrite(3, 5000000)

	opts := Opts[uint]{
		Divider:          RateDivider,
		Feedback:         msr.GetFeedback(),
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
		Output:           msr.GetOutput(),
	}

	discipline, err := New(opts)
	require.NoError(b, err)

	defer discipline.Stop()

	_ = msr.Play(discipline)
}

func TestDisciplineFair(t *testing.T) {
	measurerOpts := measurerOpts{
		HandlersQuantity: 6,
	}

	msr := newMeasurer(measurerOpts)
	defer msr.Finalize()

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

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

	measures := msr.Play(discipline)

	require.Len(t, measures, int(msr.GetExpectedMeasuresQuantity()))
}

func TestDisciplineRate(t *testing.T) {
	measurerOpts := measurerOpts{
		HandlersQuantity: 6,
	}

	msr := newMeasurer(measurerOpts)
	defer msr.Finalize()

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

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

	measures := msr.Play(discipline)

	require.Len(t, measures, int(msr.GetExpectedMeasuresQuantity()))
}

func TestDisciplineFairUnbuffered(t *testing.T) {
	measurerOpts := measurerOpts{
		HandlersQuantity: 6,
		UnbufferedInput:  true,
	}

	msr := newMeasurer(measurerOpts)
	defer msr.Finalize()

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

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

	measures := msr.Play(discipline)

	require.Len(t, measures, int(msr.GetExpectedMeasuresQuantity()))
}

func TestDisciplineRateUnbuffered(t *testing.T) {
	measurerOpts := measurerOpts{
		HandlersQuantity: 6,
		UnbufferedInput:  true,
	}

	msr := newMeasurer(measurerOpts)
	defer msr.Finalize()

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

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

	measures := msr.Play(discipline)

	require.Len(t, measures, int(msr.GetExpectedMeasuresQuantity()))
}

func TestDisciplineAddRemoveInput(t *testing.T) {
	measurerOpts := measurerOpts{
		HandlersQuantity: 6,
	}

	msr := newMeasurer(measurerOpts)
	defer msr.Finalize()

	msr.AddWrite(1, 1000000)
	msr.AddWrite(2, 1000000)
	msr.AddWrite(3, 1000000)

	inputs := msr.GetInputs()

	opts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         msr.GetFeedback(),
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Output:           msr.GetOutput(),
	}

	discipline, err := New(opts)
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

	measures := msr.Play(discipline)

	<-waiter

	require.Len(t, measures, int(msr.GetExpectedMeasuresQuantity()))
}

func TestDisciplineBadDivider(t *testing.T) {
	measurerOpts := measurerOpts{
		HandlersQuantity: 6,
	}

	msr := newMeasurer(measurerOpts)
	defer msr.Finalize()

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

	divider := func(priorities []uint, dividend uint, distribution map[uint]uint) map[uint]uint {
		out := FairDivider(priorities, dividend, distribution)

		for priority := range out {
			out[priority] *= 2
		}

		return out
	}

	opts := Opts[uint]{
		Divider:          divider,
		Feedback:         msr.GetFeedback(),
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
		Output:           msr.GetOutput(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	defer discipline.Stop()

	measures := msr.Play(discipline)

	require.NotEqual(t, int(msr.GetExpectedMeasuresQuantity()), len(measures))
}

func TestDisciplineStop(t *testing.T) {
	measurerOpts := measurerOpts{
		HandlersQuantity: 6,
	}

	msr := newMeasurer(measurerOpts)
	defer msr.Finalize()

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

	msr.SetProcessDelay(1, 10*time.Microsecond)
	msr.SetProcessDelay(2, 10*time.Microsecond)
	msr.SetProcessDelay(3, 10*time.Microsecond)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	opts := Opts[uint]{
		Ctx:              ctx,
		Divider:          RateDivider,
		Feedback:         msr.GetFeedback(),
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
		Output:           msr.GetOutput(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	defer discipline.Stop()
	defer discipline.Stop()

	measures := msr.Play(discipline)

	require.NotEqual(t, int(msr.GetExpectedMeasuresQuantity()), len(measures))
}

func TestDisciplineGracefulStop(t *testing.T) {
	measurerOpts := measurerOpts{
		HandlersQuantity: 6,
	}

	msr := newMeasurer(measurerOpts)

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

	opts := Opts[uint]{
		Divider:          RateDivider,
		Feedback:         msr.GetFeedback(),
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
		Output:           msr.GetOutput(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	go func() {
		discipline.GracefulStop()
	}()

	measures := msr.Play(discipline)
	msr.Finalize()

	require.Len(t, measures, int(msr.GetExpectedMeasuresQuantity()))
}

func TestDisciplineFairOverQuantity(t *testing.T) {
	handlersQuantity := uint(6)

	measurerOpts := measurerOpts{
		HandlersQuantity: 2 * handlersQuantity,
	}

	msr := newMeasurer(measurerOpts)
	defer msr.Finalize()

	msr.AddWrite(1, 1000000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 10000)

	opts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         msr.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           msr.GetInputs(),
		Output:           msr.GetOutput(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	defer discipline.Stop()

	measures := msr.Play(discipline)

	quantities := calcInProcessing(measures, 100*time.Millisecond)

	for priority := range quantities {
		for id := range quantities[priority] {
			require.LessOrEqual(t, quantities[priority][id].Quantity, handlersQuantity)
		}
	}
}

func TestDisciplineRateOverQuantity(t *testing.T) {
	handlersQuantity := uint(6)

	measurerOpts := measurerOpts{
		HandlersQuantity: 2 * handlersQuantity,
	}

	msr := newMeasurer(measurerOpts)
	defer msr.Finalize()

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

	opts := Opts[uint]{
		Divider:          RateDivider,
		Feedback:         msr.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           msr.GetInputs(),
		Output:           msr.GetOutput(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	defer discipline.Stop()

	measures := msr.Play(discipline)

	quantities := calcInProcessing(measures, 100*time.Millisecond)

	for priority := range quantities {
		for id := range quantities[priority] {
			require.LessOrEqual(t, quantities[priority][id].Quantity, handlersQuantity)
		}
	}
}

func TestDisciplineFairTooSmallHandlersQuantity(t *testing.T) {
	measurerOpts := measurerOpts{
		HandlersQuantity: 6,
	}

	msr := newMeasurer(measurerOpts)
	defer msr.Finalize()

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

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

	measures := msr.Play(discipline)

	require.Len(t, measures, int(msr.GetExpectedMeasuresQuantity()))
}

func TestDisciplineRateTooSmallHandlersQuantity(t *testing.T) {
	measurerOpts := measurerOpts{
		HandlersQuantity: 5,
	}

	msr := newMeasurer(measurerOpts)
	defer msr.Finalize()

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

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

	measures := msr.Play(discipline)

	require.Len(t, measures, int(msr.GetExpectedMeasuresQuantity()))
}

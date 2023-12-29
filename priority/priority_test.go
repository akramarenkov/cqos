package priority

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func BenchmarkDisciplineFair(b *testing.B) {
	handlersQuantity := uint(600)

	measurerOpts := measurerOpts{
		DisableMeasures:  true,
		HandlersQuantity: handlersQuantity,
	}

	measurer := newMeasurer(measurerOpts)
	defer measurer.Finalize()

	measurer.AddWrite(1, 5000000)
	measurer.AddWrite(2, 5000000)
	measurer.AddWrite(3, 5000000)

	disciplineOpts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         measurer.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           measurer.GetInputs(),
		Output:           measurer.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(b, err)

	defer discipline.Stop()

	_ = measurer.Play(discipline)
}

func BenchmarkDisciplineRate(b *testing.B) {
	handlersQuantity := uint(600)

	measurerOpts := measurerOpts{
		DisableMeasures:  true,
		HandlersQuantity: handlersQuantity,
	}

	measurer := newMeasurer(measurerOpts)
	defer measurer.Finalize()

	measurer.AddWrite(1, 5000000)
	measurer.AddWrite(2, 5000000)
	measurer.AddWrite(3, 5000000)

	disciplineOpts := Opts[uint]{
		Divider:          RateDivider,
		Feedback:         measurer.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           measurer.GetInputs(),
		Output:           measurer.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(b, err)

	defer discipline.Stop()

	_ = measurer.Play(discipline)
}

func BenchmarkDisciplineFairUnbuffered(b *testing.B) {
	handlersQuantity := uint(600)

	measurerOpts := measurerOpts{
		DisableMeasures:  true,
		HandlersQuantity: handlersQuantity,
		UnbufferedInput:  true,
	}

	measurer := newMeasurer(measurerOpts)
	defer measurer.Finalize()

	measurer.AddWrite(1, 5000000)
	measurer.AddWrite(2, 5000000)
	measurer.AddWrite(3, 5000000)

	disciplineOpts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         measurer.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           measurer.GetInputs(),
		Output:           measurer.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(b, err)

	defer discipline.Stop()

	_ = measurer.Play(discipline)
}

func BenchmarkDisciplineRateUnbuffered(b *testing.B) {
	handlersQuantity := uint(600)

	measurerOpts := measurerOpts{
		DisableMeasures:  true,
		HandlersQuantity: handlersQuantity,
		UnbufferedInput:  true,
	}

	measurer := newMeasurer(measurerOpts)
	defer measurer.Finalize()

	measurer.AddWrite(1, 5000000)
	measurer.AddWrite(2, 5000000)
	measurer.AddWrite(3, 5000000)

	disciplineOpts := Opts[uint]{
		Divider:          RateDivider,
		Feedback:         measurer.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           measurer.GetInputs(),
		Output:           measurer.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(b, err)

	defer discipline.Stop()

	_ = measurer.Play(discipline)
}

func TestDisciplineRate(t *testing.T) {
	handlersQuantity := uint(6)

	measurerOpts := measurerOpts{
		HandlersQuantity: handlersQuantity,
	}

	measurer := newMeasurer(measurerOpts)
	defer measurer.Finalize()

	measurer.AddWrite(1, 100000)
	measurer.AddWrite(2, 100000)
	measurer.AddWrite(3, 100000)

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

	require.Equal(t, int(measurer.GetExpectedItemsQuantity()), len(filterByKind(measures, measureKindReceived)))
}

func TestDisciplineFair(t *testing.T) {
	handlersQuantity := uint(6)

	measurerOpts := measurerOpts{
		HandlersQuantity: handlersQuantity,
	}

	measurer := newMeasurer(measurerOpts)
	defer measurer.Finalize()

	measurer.AddWrite(1, 100000)
	measurer.AddWrite(2, 100000)
	measurer.AddWrite(3, 100000)

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

	require.Equal(t, int(measurer.GetExpectedItemsQuantity()), len(filterByKind(measures, measureKindReceived)))
}

func TestDisciplineRateUnbuffered(t *testing.T) {
	handlersQuantity := uint(6)

	measurerOpts := measurerOpts{
		HandlersQuantity: handlersQuantity,
		UnbufferedInput:  true,
	}

	measurer := newMeasurer(measurerOpts)
	defer measurer.Finalize()

	measurer.AddWrite(1, 100000)
	measurer.AddWrite(2, 100000)
	measurer.AddWrite(3, 100000)

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

	require.Equal(t, int(measurer.GetExpectedItemsQuantity()), len(filterByKind(measures, measureKindReceived)))
}

func TestDisciplineFairUnbuffered(t *testing.T) {
	handlersQuantity := uint(6)

	measurerOpts := measurerOpts{
		HandlersQuantity: handlersQuantity,
		UnbufferedInput:  true,
	}

	measurer := newMeasurer(measurerOpts)
	defer measurer.Finalize()

	measurer.AddWrite(1, 100000)
	measurer.AddWrite(2, 100000)
	measurer.AddWrite(3, 100000)

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

	require.Equal(t, int(measurer.GetExpectedItemsQuantity()), len(filterByKind(measures, measureKindReceived)))
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

	measurerOpts := measurerOpts{
		HandlersQuantity: handlersQuantity,
	}

	measurer := newMeasurer(measurerOpts)
	defer measurer.Finalize()

	measurer.AddWrite(1, 1000000)
	measurer.AddWrite(2, 1000000)
	measurer.AddWrite(3, 1000000)

	inputs := measurer.GetInputs()

	disciplineOpts := Opts[uint]{
		Divider:          FairDivider,
		Feedback:         measurer.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Output:           measurer.GetOutput(),
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

	measures := measurer.Play(discipline)

	<-waiter

	require.Equal(t, int(measurer.GetExpectedItemsQuantity()), len(filterByKind(measures, measureKindReceived)))
}

func TestDisciplineBadDivider(t *testing.T) {
	handlersQuantity := uint(6)

	measurerOpts := measurerOpts{
		HandlersQuantity: handlersQuantity,
	}

	measurer := newMeasurer(measurerOpts)
	defer measurer.Finalize()

	measurer.AddWrite(1, 100000)
	measurer.AddWrite(2, 100000)
	measurer.AddWrite(3, 100000)

	divider := func(priorities []uint, dividend uint, distribution map[uint]uint) map[uint]uint {
		out := FairDivider(priorities, dividend, distribution)

		for priority, quantity := range out {
			out[priority] = 2 * quantity
		}

		return out
	}

	disciplineOpts := Opts[uint]{
		Divider:          divider,
		Feedback:         measurer.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           measurer.GetInputs(),
		Output:           measurer.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()

	measures := measurer.Play(discipline)

	require.NotEqual(t, int(measurer.GetExpectedItemsQuantity()), len(filterByKind(measures, measureKindReceived)))
}

func TestDisciplineStop(t *testing.T) {
	handlersQuantity := uint(6)

	measurerOpts := measurerOpts{
		HandlersQuantity: handlersQuantity,
	}

	measurer := newMeasurer(measurerOpts)
	defer measurer.Finalize()

	measurer.AddWrite(1, 100000)
	measurer.AddWrite(2, 100000)
	measurer.AddWrite(3, 100000)

	measurer.SetProcessDelay(1, 10*time.Microsecond)
	measurer.SetProcessDelay(2, 10*time.Microsecond)
	measurer.SetProcessDelay(3, 10*time.Microsecond)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	disciplineOpts := Opts[uint]{
		Ctx:              ctx,
		Divider:          RateDivider,
		Feedback:         measurer.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           measurer.GetInputs(),
		Output:           measurer.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	defer discipline.Stop()
	defer discipline.Stop()

	measures := measurer.Play(discipline)

	require.NotEqual(t, int(measurer.GetExpectedItemsQuantity()), len(filterByKind(measures, measureKindReceived)))
}

func TestDisciplineGracefulStop(t *testing.T) {
	handlersQuantity := uint(6)

	measurerOpts := measurerOpts{
		HandlersQuantity: handlersQuantity,
	}

	measurer := newMeasurer(measurerOpts)

	measurer.AddWrite(1, 100000)
	measurer.AddWrite(2, 100000)
	measurer.AddWrite(3, 100000)

	disciplineOpts := Opts[uint]{
		Divider:          RateDivider,
		Feedback:         measurer.GetFeedback(),
		HandlersQuantity: handlersQuantity,
		Inputs:           measurer.GetInputs(),
		Output:           measurer.GetOutput(),
	}

	discipline, err := New(disciplineOpts)
	require.NoError(t, err)

	go func() {
		discipline.GracefulStop()
	}()

	measures := measurer.Play(discipline)
	measurer.Finalize()

	require.Equal(t, int(measurer.GetExpectedItemsQuantity()), len(filterByKind(measures, measureKindReceived)))
}

func TestDisciplineRateOverQuantity(t *testing.T) {
	handlersQuantity := uint(6)

	measurerOpts := measurerOpts{
		HandlersQuantity: 2 * handlersQuantity,
	}

	measurer := newMeasurer(measurerOpts)
	defer measurer.Finalize()

	measurer.AddWrite(1, 100000)
	measurer.AddWrite(2, 100000)
	measurer.AddWrite(3, 100000)

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

	quantities := calcInProcessing(measures, 100*time.Millisecond)

	for priority := range quantities {
		for id := range quantities[priority] {
			require.LessOrEqual(t, quantities[priority][id].Quantity, handlersQuantity)
		}
	}
}

func TestDisciplineFairOverQuantity(t *testing.T) {
	handlersQuantity := uint(6)

	measurerOpts := measurerOpts{
		HandlersQuantity: 2 * handlersQuantity,
	}

	measurer := newMeasurer(measurerOpts)
	defer measurer.Finalize()

	measurer.AddWrite(1, 1000000)
	measurer.AddWrite(2, 100000)
	measurer.AddWrite(3, 10000)

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

	quantities := calcInProcessing(measures, 100*time.Millisecond)

	for priority := range quantities {
		for id := range quantities[priority] {
			require.LessOrEqual(t, quantities[priority][id].Quantity, handlersQuantity)
		}
	}
}

func TestDisciplineRateFatalDividingError(t *testing.T) {
	handlersQuantity := uint(5)

	measurerOpts := measurerOpts{
		HandlersQuantity: handlersQuantity,
	}

	measurer := newMeasurer(measurerOpts)
	defer measurer.Finalize()

	measurer.AddWrite(1, 100000)
	measurer.AddWrite(2, 100000)
	measurer.AddWrite(3, 100000)

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

	require.Equal(t, int(measurer.GetExpectedItemsQuantity()), len(filterByKind(measures, measureKindReceived)))
}

func TestDisciplineFairFatalDividingError(t *testing.T) {
	handlersQuantity := uint(6)

	measurerOpts := measurerOpts{
		HandlersQuantity: handlersQuantity,
	}

	measurer := newMeasurer(measurerOpts)
	defer measurer.Finalize()

	measurer.AddWrite(1, 100000)
	measurer.AddWrite(2, 100000)
	measurer.AddWrite(3, 100000)

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

	require.Equal(t, int(measurer.GetExpectedItemsQuantity()), len(filterByKind(measures, measureKindReceived)))
}

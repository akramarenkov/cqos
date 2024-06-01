package measurer

import (
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/priority/internal/unmanaged"
	"github.com/stretchr/testify/require"
)

func TestGeneral(t *testing.T) {
	measurerOpts := Opts{
		HandlersQuantity: 6,
	}

	msr := New(measurerOpts)

	msr.AddWrite(1, 1000)

	msr.AddWrite(2, 500)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 1*time.Second)
	msr.AddWrite(2, 500)

	msr.AddWrite(3, 500)
	msr.AddWaitDevastation(3)
	msr.AddDelay(3, 1*time.Second)
	msr.AddWrite(3, 500)

	msr.SetProcessDelay(1, 1*time.Millisecond)
	msr.SetProcessDelay(2, 1*time.Millisecond)
	msr.SetProcessDelay(3, 1*time.Millisecond)

	opts := unmanaged.Opts[uint]{
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
	}

	discipline, err := unmanaged.New(opts)
	require.NoError(t, err)

	measures := msr.Play(discipline)

	require.Equal(t, 3000, int(msr.GetExpectedItemsQuantity()))
	require.Len(t, measures, int(msr.GetExpectedMeasuresQuantity()))
}

func TestWriteWithDelay(t *testing.T) {
	measurerOpts := Opts{
		HandlersQuantity: 6,
	}

	msr := New(measurerOpts)

	msr.AddWriteWithDelay(1, 1000, 1*time.Millisecond)
	msr.AddWriteWithDelay(2, 1000, 1*time.Millisecond)
	msr.AddWriteWithDelay(3, 1000, 1*time.Millisecond)

	opts := unmanaged.Opts[uint]{
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
	}

	discipline, err := unmanaged.New(opts)
	require.NoError(t, err)

	measures := msr.Play(discipline)

	require.Equal(t, 3000, int(msr.GetExpectedItemsQuantity()))
	require.Len(t, measures, int(msr.GetExpectedMeasuresQuantity()))
}

func TestDisableMeasures(t *testing.T) {
	measurerOpts := Opts{
		HandlersQuantity: 6,
		DisableMeasures:  true,
	}

	msr := New(measurerOpts)

	msr.AddWrite(1, 1000)
	msr.AddWrite(2, 1000)
	msr.AddWrite(3, 1000)

	opts := unmanaged.Opts[uint]{
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
	}

	discipline, err := unmanaged.New(opts)
	require.NoError(t, err)

	measures := msr.Play(discipline)

	require.Equal(t, 3000, int(msr.GetExpectedItemsQuantity()))
	require.Equal(t, 0, int(msr.GetExpectedMeasuresQuantity()))
	require.Len(t, measures, int(msr.GetExpectedMeasuresQuantity()))
}

func TestBufferedInput(t *testing.T) {
	measurerOpts := Opts{
		HandlersQuantity: 6,
	}

	msr := New(measurerOpts)

	msr.AddWrite(1, 1000)
	msr.AddWrite(2, 1000)
	msr.AddWrite(3, 1000)

	require.Len(t, msr.GetInputs(), 3)

	for _, channel := range msr.GetInputs() {
		require.NotEqual(t, 0, cap(channel))
	}
}

func TestUnbufferedInput(t *testing.T) {
	measurerOpts := Opts{
		HandlersQuantity: 6,
		UnbufferedInput:  true,
	}

	msr := New(measurerOpts)

	msr.AddWrite(1, 1000)
	msr.AddWrite(2, 1000)
	msr.AddWrite(3, 1000)

	require.Len(t, msr.GetInputs(), 3)

	for _, channel := range msr.GetInputs() {
		require.Equal(t, 0, cap(channel))
	}
}

func TestFail(t *testing.T) {
	measurerOpts := Opts{
		HandlersQuantity: 6,
	}

	msr := New(measurerOpts)

	msr.AddWrite(1, 1000)
	msr.AddWrite(2, 1000)
	msr.AddWrite(3, 1000)

	opts := unmanaged.Opts[uint]{
		FailAt:           500,
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
	}

	discipline, err := unmanaged.New(opts)
	require.NoError(t, err)

	_ = msr.Play(discipline)
}

func TestFailAtWaitDevastation(t *testing.T) {
	measurerOpts := Opts{
		HandlersQuantity: 6,
	}

	msr := New(measurerOpts)

	msr.AddWrite(1, 500)
	msr.AddWaitDevastation(1)
	msr.AddWrite(1, 500)

	msr.AddWrite(2, 500)
	msr.AddWaitDevastation(2)
	msr.AddWrite(2, 500)

	msr.AddWrite(3, 500)
	msr.AddWaitDevastation(3)
	msr.AddWrite(3, 500)

	opts := unmanaged.Opts[uint]{
		FailAt:           500,
		HandlersQuantity: measurerOpts.HandlersQuantity,
		Inputs:           msr.GetInputs(),
	}

	discipline, err := unmanaged.New(opts)
	require.NoError(t, err)

	_ = msr.Play(discipline)
}

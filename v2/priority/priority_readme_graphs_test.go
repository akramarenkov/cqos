package priority

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/env"
	"github.com/akramarenkov/cqos/v2/priority/divider"
	"github.com/akramarenkov/cqos/v2/priority/internal/common"
	"github.com/akramarenkov/cqos/v2/priority/internal/measurer"
	"github.com/akramarenkov/cqos/v2/priority/internal/research"
	"github.com/akramarenkov/cqos/v2/priority/internal/unmanaged"

	"github.com/stretchr/testify/require"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

func TestReadmeGraph(t *testing.T) {
	t.Run(
		"equaling",
		func(t *testing.T) {
			t.Parallel()
			testReadmeGraph(t, true)
		},
	)

	t.Run(
		"unmanaged",
		func(t *testing.T) {
			t.Parallel()
			testReadmeGraph(t, false)
		},
	)
}

func testReadmeGraph(t *testing.T, equaling bool) {
	if os.Getenv(env.EnableGraphs) == "" {
		t.SkipNow()
	}

	measurerOpts := measurer.Opts{
		HandlersQuantity: 6,
	}

	msr := measurer.New(measurerOpts)

	msr.AddWrite(1, 500)
	msr.AddWrite(2, 500)
	msr.AddWrite(3, 500)

	msr.SetProcessDelay(1, 100*time.Millisecond)
	msr.SetProcessDelay(2, 50*time.Millisecond)
	msr.SetProcessDelay(3, 10*time.Millisecond)

	overTimeResolution := 100 * time.Millisecond
	overTimeUnit := time.Second
	overTimeUnitName := "seconds"

	if equaling {
		opts := Opts[uint]{
			Divider:          divider.Fair,
			HandlersQuantity: measurerOpts.HandlersQuantity,
			Inputs:           msr.GetInputs(),
		}

		discipline, err := New(opts)
		require.NoError(t, err)

		createReadmeGraph(
			t,
			"./doc/different-processing-time-equaling.svg",
			msr.Play(discipline),
			overTimeResolution,
			overTimeUnit,
			overTimeUnitName,
		)

		return
	}

	unmanagedOpts := unmanaged.Opts[uint]{
		Inputs: msr.GetInputs(),
	}

	unmanaged, err := unmanaged.New(unmanagedOpts)
	require.NoError(t, err)

	createReadmeGraph(
		t,
		"./doc/different-processing-time-unmanagement.svg",
		msr.Play(unmanaged),
		overTimeResolution,
		overTimeUnit,
		overTimeUnitName,
	)
}

func createReadmeGraph(
	t *testing.T,
	fileName string,
	measures []measurer.Measure,
	overTimeResolution time.Duration,
	overTimeUnit time.Duration,
	overTimeUnitName string,
) {
	received := research.FilterByKind(measures, measurer.MeasureKindReceived)
	researched := research.CalcDataQuantity(received, overTimeResolution)

	serieses := make([]chart.Series, 0, len(researched))
	priorities := make([]uint, 0, len(researched))

	for priority := range researched {
		priorities = append(priorities, priority)
	}

	// To keep the legends in the same order
	common.SortPriorities(priorities)

	for _, priority := range priorities {
		xaxis := make([]float64, len(researched[priority]))
		yaxis := make([]float64, len(researched[priority]))

		for id, item := range researched[priority] {
			xaxis[id] = float64(item.RelativeTime) / float64(overTimeUnit)
			yaxis[id] = float64(item.Quantity)
		}

		series := chart.ContinuousSeries{
			Name:    fmt.Sprintf("Data with priority %d", priority),
			XValues: xaxis,
			YValues: yaxis,
			Style:   chart.Style{StrokeWidth: 4},
		}

		serieses = append(serieses, series)
	}

	graph := chart.Chart{
		Title:        "Data retrieval graph",
		ColorPalette: readmeColorPalette{},
		Background: chart.Style{
			Padding: chart.Box{
				Top:  50,
				Left: 140,
			},
			FillColor: chart.ColorTransparent,
		},
		Canvas: chart.Style{
			FillColor: chart.ColorTransparent,
		},
		XAxis: chart.XAxis{
			Name: "Time, " + overTimeUnitName,
		},
		YAxis: chart.YAxis{
			Name: "Quantity of data received by handlers, pieces",
		},
		Series: serieses,
	}

	graph.Elements = []chart.Renderable{
		chart.LegendLeft(&graph, chart.Style{FillColor: chart.ColorTransparent}),
	}

	file, err := os.Create(fileName)
	require.NoError(t, err)

	defer file.Close()

	err = graph.Render(chart.SVG, file)
	require.NoError(t, err)
}

type readmeColorPalette struct{}

func (readmeColorPalette) BackgroundColor() drawing.Color {
	return chart.DefaultColorPalette.BackgroundColor()
}

func (readmeColorPalette) BackgroundStrokeColor() drawing.Color {
	return chart.DefaultColorPalette.BackgroundStrokeColor()
}

func (readmeColorPalette) CanvasColor() drawing.Color {
	return chart.DefaultColorPalette.CanvasColor()
}

func (readmeColorPalette) CanvasStrokeColor() drawing.Color {
	return chart.DefaultColorPalette.CanvasStrokeColor()
}

func (readmeColorPalette) AxisStrokeColor() drawing.Color {
	return chart.DefaultColorPalette.AxisStrokeColor()
}

func (readmeColorPalette) TextColor() drawing.Color {
	return chart.DefaultColorPalette.TextColor()
}

func (readmeColorPalette) GetSeriesColor(index int) drawing.Color {
	colors := []drawing.Color{
		{R: 0x54, G: 0x70, B: 0xc6, A: 255},
		{R: 0x91, G: 0xcc, B: 0x75, A: 255},
		{R: 0xfa, G: 0xc8, B: 0x58, A: 255},
		{R: 0xee, G: 0x66, B: 0x66, A: 255},
		{R: 0x73, G: 0xc0, B: 0xde, A: 255},
		{R: 0x3b, G: 0xa2, B: 0x72, A: 255},
		{R: 0xfc, G: 0x84, B: 0x52, A: 255},
		{R: 0x9a, G: 0x60, B: 0xb4, A: 255},
		{R: 0xea, G: 0x7c, B: 0xcc, A: 255},
	}

	return colors[index%len(colors)]
}

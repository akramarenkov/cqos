// Internal package with Spinner implementation that used to iterate over integer
// values that returns the begin value when trying to go beyond the end value.
package spinner

type Spinner struct {
	actual int
	begin  int
	end    int
}

// Creates Spinner instance.
func New(begin int, end int) *Spinner {
	spn := &Spinner{
		actual: begin,
		begin:  begin,
		end:    end,
	}

	return spn
}

// Returns actual value of counter.
func (spn *Spinner) Actual() int {
	return spn.actual
}

// Increases the current value of the counter, if its next value exceeds the end value,
// it will be reset to the begin value.
func (spn *Spinner) Spin() {
	spn.actual = spn.spin()
}

func (spn *Spinner) spin() int {
	if spn.actual >= spn.end {
		return spn.begin
	}

	return spn.actual + 1
}

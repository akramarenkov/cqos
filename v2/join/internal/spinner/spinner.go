// Internal package used to iterate over integer values that returns the begin
// value when trying to go beyond the end value.
package spinner

type Spinner struct {
	actual int
	begin  int
	end    int
}

func New(begin int, end int) *Spinner {
	spn := &Spinner{
		actual: begin,
		begin:  begin,
		end:    end,
	}

	return spn
}

func (spn *Spinner) Actual() int {
	return spn.actual
}

func (spn *Spinner) Spin() {
	spn.actual = spn.spin()
}

func (spn *Spinner) spin() int {
	if spn.actual >= spn.end {
		return spn.begin
	}

	return spn.actual + 1
}

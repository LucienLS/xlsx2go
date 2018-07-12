package xlsx

//RangeIterator is a interface for iterating cells inside of range
type RangeIterator interface {
	//Next returns next Cell in range and corresponding indexes
	Next() (*Cell, int, int)

	//HasNext returns true if there are cells to iterate or false in other case
	HasNext() bool
}

//rangeIterator is object that holds required information for range's iterator
type rangeIterator struct {
	r *Range
	cIdx int
	rIdx int
}

var _ RangeIterator = (*rangeIterator)(nil)

func newRangeIterator(r *Range) RangeIterator {
	return &rangeIterator {
		r:  r,
		cIdx:  r.fromCol,
		rIdx:  r.fromRow,
	}
}

//Next returns next Cell in range and corresponding indexes
func (i *rangeIterator) Next() (*Cell, int, int) {
	c := i.r.sheet.Cell(i.cIdx, i.rIdx)

	i.cIdx++
	if i.cIdx > i.r.toCol {
		i.cIdx = i.r.fromCol
		i.rIdx++
	}

	return c, i.cIdx, i.rIdx
}

//HasNext returns true if there are cells to iterate or false in other case
func (i *rangeIterator) HasNext() bool {
	return i.rIdx <= i.r.toRow
}

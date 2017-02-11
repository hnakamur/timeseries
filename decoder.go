package timeseries

import "io"

type Decoder struct {
	rd           io.Reader
	pointsLeft   int
	prevTime     uint64
	prevTmeDelta int64
	prevValBits  uint64
}

func NewDecoder(rd io.Reader) *Decoder {
	return &Decoder{
		rd:         rd,
		pointsLeft: -1,
	}
}

func (d *Decoder) Decode(points []Point) (n int, err error) {
	return 0, nil
}

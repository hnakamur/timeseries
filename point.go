package timeseries

import (
	"bytes"
	"fmt"
	"io"
)

type Point struct {
	Timestamp uint32
	Value     float64
}

func Marshal(t0 uint32, points []Point) ([]byte, error) {
	var b bytes.Buffer
	enc := NewEncoder(&b)
	err := enc.EncodeHeader(t0)
	if err != nil {
		return nil, fmt.Errorf("failed to encode time series header: err=%+v", err)
	}

	for _, p := range points {
		err := enc.EncodePoint(p)
		if err != nil {
			return nil, fmt.Errorf("failed to encode time series point: err=%+v", err)
		}
	}

	err = enc.Finish()
	if err != nil {
		return nil, fmt.Errorf("failed to encode time series finish marker: err=%+v", err)
	}

	return b.Bytes(), nil
}

func Unmarshal(data []byte) (t0 uint32, points []Point, err error) {
	b := bytes.NewBuffer(data)
	dec := NewDecoder(b)

	t0, err = dec.DecodeHeader()
	if err != nil {
		return 0, nil, fmt.Errorf("failed to decode time series header: err=%+v", err)
	}

	for {
		var p Point
		p, err = dec.DecodePoint()
		if err == io.EOF {
			break
		} else if err != nil {
			return 0, nil, fmt.Errorf("failed to decode time series point: err=%+v", err)
		}
		points = append(points, p)
	}

	return t0, points, nil
}

package timeseries

import (
	"bytes"
	"fmt"
	"io"
)

// Point is a time-series data point.
type Point struct {
	// Timestamp represents the seconds since 1970-01-01 00:00:00 +0000 UTC.
	// The max value is 2106-02-07 06:28:15 +0000 UTC.
	// You can verify the max value at https://play.golang.org/p/XgHveabiUd
	Timestamp uint32

	// Value represents the data point value.
	Value float64
}

// Marshal encodes a block timestamp and data points to bytes.
func Marshal(t0 uint32, points []Point) ([]byte, error) {
	var b bytes.Buffer
	enc := NewEncoder(&b)
	err := enc.EncodeHeader(t0)
	if err != nil {
		return nil, fmt.Errorf("failed to encode time series header: err=%+v", err)
	}

	for _, p := range points {
		err = enc.EncodePoint(p)
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

// Unmarshal decodes bytes to a block timestamp and data points.
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

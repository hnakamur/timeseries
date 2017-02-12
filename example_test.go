package timeseries_test

import (
	"encoding/hex"
	"fmt"
	"time"

	"bitbucket.org/hnakamur/timeseries"
)

func ExampleMarshal() {
	t0 := uint32(time.Date(2015, 3, 24, 2, 0, 0, 0, time.UTC).Unix())
	points := []timeseries.Point{
		{
			Timestamp: uint32(time.Date(2015, 3, 24, 2, 1, 2, 0, time.UTC).Unix()),
			Value:     12.0,
		},
		{
			Timestamp: uint32(time.Date(2015, 3, 24, 2, 2, 2, 0, time.UTC).Unix()),
			Value:     12.0,
		},
		{
			Timestamp: uint32(time.Date(2015, 3, 24, 2, 3, 2, 0, time.UTC).Unix()),
			Value:     24.0,
		},
	}

	buf, err := timeseries.Marshal(t0, points)
	if err != nil {
		fmt.Printf("failed to marshal points: err=%+v\n", err)
		return
	}
	fmt.Println(hex.EncodeToString(buf))

	// Output: 5510c52000f900a0000000000002fc6b07ffffffffe0
}

func ExampleUnmarshal() {
	buf, err := hex.DecodeString("5510c52000f900a0000000000002fc6b07ffffffffe0")
	if err != nil {
		fmt.Printf("failed to decode hex string: err=%+v\n", err)
		return
	}

	t0, points, err := timeseries.Unmarshal(buf)
	if err != nil {
		fmt.Printf("failed to unmarshal time series: err=%+v\n", err)
		return
	}
	fmt.Printf("block timestamp=%v\n", time.Unix(int64(t0), 0).UTC())
	for _, p := range points {
		fmt.Printf("timestamp=%v, value=%f\n", time.Unix(int64(p.Timestamp), 0).UTC(), p.Value)
	}

	// Output:
	// block timestamp=2015-03-24 02:00:00 +0000 UTC
	// timestamp=2015-03-24 02:01:02 +0000 UTC, value=12.000000
	// timestamp=2015-03-24 02:02:02 +0000 UTC, value=12.000000
	// timestamp=2015-03-24 02:03:02 +0000 UTC, value=24.000000
}

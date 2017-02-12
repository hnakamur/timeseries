package timeseries_test

import (
	"encoding/hex"
	"reflect"
	"testing"
	"time"

	"github.com/hnakamur/timeseries"
)

func TestUnmarshal(t *testing.T) {
	testCases := []struct {
		input      string
		wantT0     uint32
		wantPoints []timeseries.Point
	}{
		{
			input:  "5510c52000f900a0000000000002fc6b07ffffffffe0",
			wantT0: uint32(time.Date(2015, 3, 24, 2, 0, 0, 0, time.UTC).Unix()),
			wantPoints: []timeseries.Point{
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
			},
		},
		{
			input:      "5510c520fffc0000000000000000",
			wantT0:     uint32(time.Date(2015, 3, 24, 2, 0, 0, 0, time.UTC).Unix()),
			wantPoints: nil,
		},
		{
			input:  "5510c52000f900a0000000000003ffffffffc0",
			wantT0: uint32(time.Date(2015, 3, 24, 2, 0, 0, 0, time.UTC).Unix()),
			wantPoints: []timeseries.Point{
				{
					Timestamp: uint32(time.Date(2015, 3, 24, 2, 1, 2, 0, time.UTC).Unix()),
					Value:     12.0,
				},
			},
		},
		{
			input:  "5510c52000f900a0000000000002fdbc1b0010022666666666667ffffffffe",
			wantT0: uint32(time.Date(2015, 3, 24, 2, 0, 0, 0, time.UTC).Unix()),
			wantPoints: []timeseries.Point{
				{
					Timestamp: uint32(time.Date(2015, 3, 24, 2, 1, 2, 0, time.UTC).Unix()),
					Value:     12.0,
				},
				{
					Timestamp: uint32(time.Date(2015, 3, 24, 2, 2, 2, 0, time.UTC).Unix()),
					Value:     12.5,
				},
				{
					Timestamp: uint32(time.Date(2015, 3, 24, 2, 3, 2, 0, time.UTC).Unix()),
					Value:     -24.2,
				},
			},
		},
	}

	for _, tc := range testCases {
		input, err := hex.DecodeString(tc.input)
		if err != nil {
			t.Fatalf("failed to decode input hex string: tc.input=%s, err=%+v\n", tc.input, err)
		}

		t0, points, err := timeseries.Unmarshal(input)
		if err != nil {
			t.Fatalf("failed to unmarshal time series: tc.input=%s, err=%+v\n", tc.input, err)
		}
		if t0 != tc.wantT0 {
			t.Errorf("tc.input=%s, gotT0=%d, wantT0=%d", tc.input, t0, tc.wantT0)
		}
		if !reflect.DeepEqual(points, tc.wantPoints) {
			t.Errorf("tc.input=%s, gotPoints=%+v, wantPoints=%+v", tc.input, points, tc.wantPoints)
		}
	}
}

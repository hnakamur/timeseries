package timeseries_test

import (
	"encoding/hex"
	"testing"
	"time"

	"bitbucket.org/hnakamur/timeseries"
)

func TestMarshal(t *testing.T) {
	testCases := []struct {
		t0     uint32
		points []timeseries.Point
		want   string
	}{
		{
			t0: uint32(time.Date(2015, 3, 24, 2, 0, 0, 0, time.UTC).Unix()),
			points: []timeseries.Point{
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
			want: "5510c52000f900a0000000000002fc6b07ffffffffe0",
		},
		{
			t0:     uint32(time.Date(2015, 3, 24, 2, 0, 0, 0, time.UTC).Unix()),
			points: []timeseries.Point{},
			want:   "5510c520fffc0000000000000000",
		},
		{
			t0: uint32(time.Date(2015, 3, 24, 2, 0, 0, 0, time.UTC).Unix()),
			points: []timeseries.Point{
				{
					Timestamp: uint32(time.Date(2015, 3, 24, 2, 1, 2, 0, time.UTC).Unix()),
					Value:     12.0,
				},
			},
			want: "5510c52000f900a0000000000003ffffffffc0",
		},
		{
			t0: uint32(time.Date(2015, 3, 24, 2, 0, 0, 0, time.UTC).Unix()),
			points: []timeseries.Point{
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
			want: "5510c52000f900a0000000000002fdbc1b0010022666666666667ffffffffe",
		},
	}

	for _, tc := range testCases {
		buf, err := timeseries.Marshal(tc.t0, tc.points)
		if err != nil {
			t.Fatalf("failed to marshal points: t0=%d, points=%+v, err=%+v\n", tc.t0, tc.points, err)
		}
		got := hex.EncodeToString(buf)

		if got != tc.want {
			t.Errorf("t0=%d, points=%+v, got=%s, want=%s", tc.t0, tc.points, got, tc.want)
		}
	}
}

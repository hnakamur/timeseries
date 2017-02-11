package timeseries

import (
	"fmt"
	"testing"
)

func TestNumOfLeadingZeros(t *testing.T) {
	val := uint64(0x8000000000000000)
	want := uint8(0)
	for want <= 64 {
		t.Run(fmt.Sprintf("%x", val), func(t *testing.T) {
			got := numOfLeadingZeros(val)
			if got != want {
				t.Errorf("got %d; want %d", got, want)
			}
		})

		val >>= 1
		want++
	}
}

func TestNumOfTrailingZeros(t *testing.T) {
	val := uint64(0x8000000000000000)
	want := uint8(63)
	for val > 0 {
		t.Run(fmt.Sprintf("%x", val), func(t *testing.T) {
			got := numOfTrailingZeros(val)
			if got != want {
				t.Errorf("got %d; want %d", got, want)
			}
		})

		val >>= 1
		want--
	}

	val = 0
	want = 64
	t.Run("0", func(t *testing.T) {
		got := numOfTrailingZeros(val)
		if got != want {
			t.Errorf("got %d; want %d", got, want)
		}
	})
}

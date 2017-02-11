package timeseries

func numOfLeadingZeros(v uint64) uint8 {
	mask := uint64(0x8000000000000000)
	c := uint8(0)
	for ; c < 64 && v&mask == 0; c++ {
		mask >>= 1
	}
	return c
}

func numOfTrailingZeros(v uint64) uint8 {
	mask := uint64(0x0000000000000001)
	c := uint8(0)
	for ; c < 64 && v&mask == 0; c++ {
		mask <<= 1
	}
	return c
}

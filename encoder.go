package timeseries

import (
	"io"
	"math"

	"github.com/dgryski/go-bitstream"
)

// This code is based on
// https://github.com/burmanm/gorilla-tsc/blob/fb984aefffb63c7b4d48c526f69db53813df2f28/src/main/java/fi/iki/yak/ts/compression/gorilla/Compressor.java

// http://www.vldb.org/pvldb/vol8/p1816-teller.pdf
// The first time stamp delta is sized at 14 bits, because that size is enough to span a bit more than 4 hours (16,384 seconds), If one chose a Gorilla block larger than 4 hours, this size would increase.
const nBitsFirstDelta = 14

// Encoder encodes time series data in similar way to Facebook Gorilla
// in-memory time series database.
type Encoder struct {
	wr              *bitstream.BitWriter
	headerTimestamp uint32
	storedTimestamp uint32
	storedDelta     uint32

	storedLeadingZeros  uint8
	storedTrailingZeros uint8
	storedValueBits     uint64
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		wr:                 bitstream.NewWriter(w),
		storedLeadingZeros: math.MaxInt8,
	}
}

func (e *Encoder) EncodeHeader(t0 uint32) error {
	err := e.wr.WriteBits(uint64(t0), 32)
	if err != nil {
		return err
	}
	e.headerTimestamp = t0
	return nil
}

func (e *Encoder) EncodePoint(p Point) error {
	if e.storedTimestamp == 0 {
		return e.writeFirst(p)
	} else {
		return e.writePoint(p)
	}
}

func (e *Encoder) Finish() error {
	if e.storedTimestamp == 0 {
		// Add finish marker with delta = 0x3FFF (nBitsFirstDelta = 14 bits), and first value = 0
		err := e.wr.WriteBits(1<<nBitsFirstDelta-1, nBitsFirstDelta)
		if err != nil {
			return err
		}
		err = e.wr.WriteBits(0, 64)
		if err != nil {
			return err
		}
	} else {
		// Add finish marker with deltaOfDelta = 0xFFFFFFFF, and value xor = 0
		err := e.wr.WriteBits(0x0F, 4)
		if err != nil {
			return err
		}
		err = e.wr.WriteBits(0xFFFFFFFF, 32)
		if err != nil {
			return err
		}
		err = e.wr.WriteBit(bitstream.Zero)
		if err != nil {
			return err
		}
	}

	return e.wr.Flush(bitstream.Zero)
}

func (e *Encoder) writeFirst(p Point) error {
	delta := p.Timestamp - e.headerTimestamp
	e.storedTimestamp = p.Timestamp
	e.storedDelta = delta
	e.storedValueBits = math.Float64bits(p.Value)

	err := e.wr.WriteBits(uint64(delta), nBitsFirstDelta)
	if err != nil {
		return err
	}

	return e.wr.WriteBits(e.storedValueBits, 64)
}

func (e *Encoder) writePoint(p Point) error {
	err := e.writeTimestampDeltaDelta(p.Timestamp)
	if err != nil {
		return err
	}

	return e.writeValueXor(p.Value)
}

func (e *Encoder) writeTimestampDeltaDelta(timestamp uint32) error {
	delta := timestamp - e.storedTimestamp
	deltaDelta := int64(delta) - int64(e.storedDelta)
	e.storedTimestamp = timestamp
	e.storedDelta = delta

	switch {
	case deltaDelta == 0:
		err := e.wr.WriteBit(bitstream.Zero)
		if err != nil {
			return err
		}
	case -63 <= deltaDelta && deltaDelta <= 64:
		err := e.wr.WriteBits(0x02, 2) // write 2 bits header '10'
		if err != nil {
			return err
		}
		err = writeInt64Bits(e.wr, deltaDelta, 7)
		if err != nil {
			return err
		}
	case -255 <= deltaDelta && deltaDelta <= 256:
		err := e.wr.WriteBits(0x06, 3) // write 3 bits header '110'
		if err != nil {
			return err
		}
		err = writeInt64Bits(e.wr, deltaDelta, 9)
		if err != nil {
			return err
		}
	case -2047 <= deltaDelta && deltaDelta <= 2048:
		err := e.wr.WriteBits(0x0E, 4) // write 4 bits header '1110'
		if err != nil {
			return err
		}
		err = writeInt64Bits(e.wr, deltaDelta, 12)
		if err != nil {
			return err
		}
	default:
		err := e.wr.WriteBits(0x0F, 4) // write 4 bits header '1111'
		if err != nil {
			return err
		}
		err = writeInt64Bits(e.wr, deltaDelta, 32)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeInt64Bits(w *bitstream.BitWriter, i int64, nbits uint) error {
	var u uint64
	if i >= 0 || nbits >= 64 {
		u = uint64(i)
	} else {
		u = uint64(1<<nbits + i)
	}
	return w.WriteBits(u, int(nbits))
}

func (e *Encoder) writeValueXor(value float64) error {
	valueBits := math.Float64bits(value)
	xor := e.storedValueBits ^ valueBits
	e.storedValueBits = valueBits

	if xor == 0 {
		return e.wr.WriteBit(bitstream.Zero)
	}

	leadingZeros := numOfLeadingZeros(xor)
	trailingZeros := numOfTrailingZeros(xor)

	err := e.wr.WriteBit(bitstream.One)
	if err != nil {
		return err
	}

	var significantBits uint8
	if leadingZeros >= e.storedLeadingZeros && trailingZeros >= e.storedTrailingZeros {
		// write existing leading
		err := e.wr.WriteBit(bitstream.Zero)
		if err != nil {
			return err
		}

		significantBits = 64 - e.storedLeadingZeros - e.storedTrailingZeros
	} else {
		e.storedLeadingZeros = leadingZeros
		e.storedTrailingZeros = trailingZeros

		// write new leading
		err := e.wr.WriteBit(bitstream.One)
		if err != nil {
			return err
		}

		err = e.wr.WriteBits(uint64(leadingZeros), 5)
		if err != nil {
			return err
		}

		significantBits = 64 - leadingZeros - trailingZeros
		err = e.wr.WriteBits(uint64(significantBits), 6)
		if err != nil {
			return err
		}
	}

	return e.wr.WriteBits(xor>>e.storedTrailingZeros, int(significantBits))
}

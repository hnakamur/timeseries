package timeseries

import (
	"errors"
	"io"
	"math"

	"github.com/dgryski/go-bitstream"
)

type Decoder struct {
	rd              *bitstream.BitReader
	headerTimestamp uint32
	storedTimestamp uint32
	storedDelta     uint32

	storedLeadingZeros  uint8
	storedTrailingZeros uint8
	storedValueBits     uint64
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		rd: bitstream.NewReader(r),
	}
}

func (d *Decoder) DecodeHeader() (t0 uint32, err error) {
	timestamp, err := d.rd.ReadBits(32)
	if err != nil {
		return 0, err
	}
	d.headerTimestamp = uint32(timestamp)
	return d.headerTimestamp, nil
}

func (d *Decoder) DecodePoint() (p Point, err error) {
	if d.storedTimestamp == 0 {
		return d.readFirst()
	} else {
		return d.readPoint()
	}
}

func (d *Decoder) readFirst() (p Point, err error) {
	delta, err := d.rd.ReadBits(nBitsFirstDelta)
	if err != nil {
		return Point{}, err
	}
	if delta == 1<<nBitsFirstDelta-1 {
		return Point{}, io.EOF
	}

	valueBits, err := d.rd.ReadBits(64)
	if err != nil {
		return Point{}, err
	}

	d.storedDelta = uint32(delta)
	d.storedTimestamp = d.headerTimestamp + d.storedDelta
	d.storedValueBits = valueBits

	return Point{
		Timestamp: d.storedTimestamp,
		Value:     math.Float64frombits(d.storedValueBits),
	}, nil
}

func (d *Decoder) readPoint() (p Point, err error) {
	t, err := d.readTmestamp()
	if err != nil {
		return Point{}, err
	}

	v, err := d.readValue()
	if err != nil {
		return Point{}, err
	}

	return Point{
		Timestamp: t,
		Value:     v,
	}, err
}

func (d *Decoder) readTmestamp() (t uint32, err error) {
	nBits, err := d.bitsToRead()
	if err != nil {
		return 0, err
	}

	var deltaDelta int64
	if nBits > 0 {
		deltaDeltaBits, err := d.rd.ReadBits(int(nBits))
		if err != nil {
			return 0, err
		}

		if nBits == 32 {
			if deltaDeltaBits == 0xFFFFFFFF {
				return 0, io.EOF
			}

			deltaDelta = int64(deltaDeltaBits)
		} else {
			// Turn unsigned uint64 back to int64
			if 1<<(nBits-1) < deltaDeltaBits {
				deltaDelta = int64(deltaDeltaBits - 1<<nBits)
			} else {
				deltaDelta = int64(deltaDeltaBits)
			}
		}
	}

	d.storedDelta += uint32(deltaDelta)
	d.storedTimestamp += d.storedDelta

	return d.storedTimestamp, nil
}

func (d *Decoder) readValue() (v float64, err error) {
	b, err := d.rd.ReadBit()
	if err != nil {
		return 0, err
	}

	if b == bitstream.One {
		b, err = d.rd.ReadBit()
		if err != nil {
			return 0, err
		}

		if b == bitstream.One {
			// New leading and trailing zeros
			storedLeadingZeros, err := d.rd.ReadBits(5)
			if err != nil {
				return 0, err
			}

			significantBits, err := d.rd.ReadBits(6)
			if err != nil {
				return 0, err
			}
			if significantBits == 0 {
				significantBits = 64
			}

			d.storedLeadingZeros = uint8(storedLeadingZeros)
			d.storedTrailingZeros = 64 - uint8(significantBits) - d.storedLeadingZeros
		}

		valueBits, err := d.rd.ReadBits(int(64 - d.storedLeadingZeros - d.storedTrailingZeros))
		if err != nil {
			return 0, err
		}

		valueBits <<= d.storedTrailingZeros
		d.storedValueBits ^= valueBits
	}

	return math.Float64frombits(d.storedValueBits), nil
}

func (d *Decoder) bitsToRead() (n uint, err error) {
	val := 0
	for i := 0; i < 4; i++ {
		val <<= 1
		b, err := d.rd.ReadBit()
		if err != nil {
			return 0, err
		}
		if b == bitstream.One {
			val |= 1
		} else {
			break
		}
	}

	switch val {
	case 0x00:
		return 0, nil
	case 0x02:
		return 7, nil
	case 0x06:
		return 9, nil
	case 0x0E:
		return 12, nil
	case 0x0F:
		return 32, nil
	default:
		return 0, errors.New("invalid bit header for bit length to read")
	}
}

package main

import (
	"encoding/hex"
	"log"
	"os"
	"time"

	"bitbucket.org/hnakamur/timeseries"
)

func main() {
	os.Exit(run())
}

func run() int {
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
		{
			Timestamp: uint32(time.Date(2015, 3, 24, 2, 4, 2, 0, time.UTC).Unix()),
			Value:     -35.5,
		},
		{
			Timestamp: uint32(time.Date(2015, 3, 24, 2, 5, 2, 0, time.UTC).Unix()),
			Value:     -33.5,
		},
		{
			Timestamp: uint32(time.Date(2015, 3, 24, 2, 6, 2, 0, time.UTC).Unix()),
			Value:     -33.5,
		},
	}

	// var b bytes.Buffer
	// enc := timeseries.NewEncoder(&b)
	// err := enc.EncodeHeader(t0)
	// if err != nil {
	// 	log.Printf("failed to encode header: err=%+v", err)
	// 	return 1
	// }
	// for i, p := range points {
	// 	err := enc.EncodePoint(p)
	// 	if err != nil {
	// 		log.Printf("failed to encode point: i=%d, p=%+v, err=%+v", i, p, err)
	// 		return 1
	// 	}
	// }
	// err = enc.Finish()
	// if err != nil {
	// 	log.Printf("failed to finish: err=%+v", err)
	// 	return 1
	// }

	// log.Printf("buf hex=%s", hex.EncodeToString(b.Bytes()))

	buf, err := timeseries.Marshal(t0, points)
	if err != nil {
		log.Printf("failed to marshal points: err=%+v", err)
		return 1
	}
	log.Printf("buf hex=%s", hex.EncodeToString(buf))

	return 0
}

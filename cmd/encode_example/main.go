package main

import (
	"bytes"
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
	b := new(bytes.Buffer)
	t0 := uint32(time.Date(2015, 3, 24, 2, 0, 0, 0, time.UTC).Unix())
	enc := timeseries.NewEncoder(b)
	err := enc.EncodeHeader(t0)
	if err != nil {
		log.Printf("failed to encode header: err=%+v", err)
		return 1
	}

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
	err = enc.Encode(points)
	if err != nil {
		log.Printf("failed to encode points: err=%+v", err)
		return 1
	}

	log.Printf("buf hex=%s", hex.EncodeToString(b.Bytes()))

	return 0
}

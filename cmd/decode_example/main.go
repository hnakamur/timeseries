package main

import (
	"bytes"
	"encoding/hex"
	"io"
	"log"
	"os"
	"time"

	"bitbucket.org/hnakamur/timeseries"
)

func main() {
	os.Exit(run())
}

func run() int {
	buf, err := hex.DecodeString("5510c52000f900a0000000000002fc6b07ffffffffe0")
	if err != nil {
		log.Printf("failed to decode hex string: err=%+v", err)
		return 1
	}

	b := bytes.NewBuffer(buf)
	dec := timeseries.NewDecoder(b)

	t0, err := dec.DecodeHeader()
	if err != nil {
		log.Printf("failed to decode header: err=%+v", err)
		return 1
	}

	log.Printf("block timestamp=%v", time.Unix(int64(t0), 0).UTC())

	for {
		p, err := dec.DecodePoint()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("failed to decode point: err=%+v", err)
			return 1
		}
		log.Printf("timestamp=%v, value=%f", time.Unix(int64(p.Timestamp), 0).UTC(), p.Value)
	}

	return 0
}

// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	usage = "Usage: %v [-B | -KB | -KiB | -MB | -MiB]\n"

	buflen = 1024 * 1024 // 1MB
)

const (
	exit_usage = iota + 1
	exit_io
)

var flags = map[string]int{
	"-B":   1,
	"-KB":  1000,
	"-KiB": 1024,
	"-MB":  1000 * 1000,
	"-MiB": 1024 * 1024,
}

func main() {
	unit, size := "KB", 1000
	if len(os.Args) == 2 {
		var ok bool
		size, ok = flags[os.Args[1]]
		if !ok {
			fmt.Fprintf(os.Stderr, usage, os.Args[0])
			os.Exit(exit_usage)
		}

		// Do this after validating that
		// the string is in the map (and
		// thus at least 1 character long)
		unit = os.Args[1][1:]
	} else if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		os.Exit(exit_usage)
	}

	r := newRateReader(os.Stdin, size, unit)
	_, err := io.Copy(os.Stdout, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "io error: %v\n", err)
		os.Exit(exit_io)
	}
}

type rateReader struct {
	r  io.Reader
	t0 time.Time
	nn int

	size int
	unit string
}

func (r *rateReader) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)
	r.nn += n
	t1 := time.Now()
	delta := t1.Sub(r.t0)
	if delta < 500*time.Millisecond {
		return n, err
	}

	r.t0 = t1
	bps := (float64(r.nn)) / delta.Seconds()
	r.nn = 0

	fmt.Fprintf(os.Stderr, "\r%8.4f %s/s", bps/float64(r.size), r.unit)
	return n, err
}

func newRateReader(r io.Reader, size int, unit string) *rateReader {
	return &rateReader{r: r,
		t0:   time.Now(),
		size: size,
		unit: unit,
	}
}

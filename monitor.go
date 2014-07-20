// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"
)

const (
	usage = "Usage: %v [-B | -KB | -KiB | -MB | -MiB] [-s | --silent]\n"

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
	defer r.Close()
	_, err := io.Copy(os.Stdout, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "io error: %v\n", err)
		os.Exit(exit_io)
	}
}

type rateReader struct {
	r  io.Reader
	t0 time.Time
	nn int64

	size int
	unit string

	err error

	exit chan struct{}
}

func (r *rateReader) Read(p []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	n, err := r.r.Read(p)
	atomic.AddInt64(&r.nn, int64(n))
	if err != nil {
		r.close()
		r.err = err
	}
	return n, err
}

func (r *rateReader) Close() error {
	r.close()
	return nil
}

func (r *rateReader) close() {
	// Include default case in case r has already
	// been closed (and thus we would have to wait
	// for print() to read from r.exit) or r has
	// already been closed twice (in which case
	// print() has already returned, in which case
	// we would block forever).
	select {
	case r.exit <- struct{}{}:
	default:
	}
}

func (r *rateReader) print() {
	for {
		select {
		case <-r.exit:
			return
		default:
			time.Sleep(500 * time.Millisecond)
			t1 := time.Now()
			nn := atomic.SwapInt64(&r.nn, 0)

			delta := t1.Sub(r.t0)
			r.t0 = t1

			bps := float64(nn) / delta.Seconds()
			fmt.Fprintf(os.Stderr, "\r%8.4f %s/s", bps/float64(r.size), r.unit)
		}
	}
}

func newRateReader(r io.Reader, size int, unit string) *rateReader {
	ret := &rateReader{r: r,
		t0:   time.Now(),
		size: size,
		unit: unit,
		exit: make(chan struct{}, 1), // Buffered so that Close() is non-blocking
	}
	go ret.print()
	return ret
}

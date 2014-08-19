// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

const (
	buflen = 1024 * 1024 // 1MB
)

var (
	unit     = flag.String("unit", "KB", "Unit for displaying rate. Options are B, KB, KiB, MB, MiB, GB, GiB.")
	progress = flag.Bool("progress", false, "Display the total amount of data copied so far.")
)

const (
	EXIT_USAGE = 2 + iota
	EXIT_IO
)

var units = map[string]int{
	"B":   1,
	"KB":  1000,
	"KiB": 1024,
	"MB":  1000 * 1000,
	"MiB": 1024 * 1024,
	"GB":  1000 * 1000 * 1000,
	"GiB": 1024 * 1024 * 1024,
}

func main() {
	flag.Parse()
	size, ok := units[*unit]
	if !ok {
		flag.Usage()
		os.Exit(EXIT_USAGE)
	}

	r := newRateReader(os.Stdin, size, *unit, *progress)
	defer r.Close()
	_, err := io.Copy(os.Stdout, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "io error: %v\n", err)
		os.Exit(EXIT_IO)
	}
}

type rateReader struct {
	r  io.Reader
	t0 time.Time
	n  int64
	nn int64

	size     int
	unit     string
	progress bool

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
			r.n += nn

			delta := t1.Sub(r.t0)
			r.t0 = t1

			bps := float64(nn) / delta.Seconds()
			fmt.Fprintf(os.Stderr, "\r%8.4f %s/s", bps/float64(r.size), r.unit)

			if r.progress {
				size, unit := 0, ""
				// If the rate is displayed in a base-10
				// unit, do the same for the progress,
				// and the same for base-2.
				if strings.Contains(r.unit, "i") {
					switch {
					case r.n < 1024:
						size, unit = 1, "B"
					case r.n < 1024*1024:
						size, unit = 1024, "KiB"
					case r.n < 1024*1024*1024:
						size, unit = 1024*1024, "MiB"
					default:
						size, unit = 1024*1024*1024, "GiB"
					}
				} else {
					switch {
					case r.n < 1000:
						size, unit = 1, "B"
					case r.n < 1000*1000:
						size, unit = 1000, "KB"
					case r.n < 1000*1000*1000:
						size, unit = 1000*1000, "MB"
					default:
						size, unit = 1000*1000*1000, "GB"
					}
				}
				fmt.Fprintf(os.Stderr, " (%.4f %s total)", float64(r.n)/float64(size), unit)
			}
		}
	}
}

func newRateReader(r io.Reader, size int, unit string, progress bool) *rateReader {
	ret := &rateReader{r: r,
		t0:       time.Now(),
		size:     size,
		unit:     unit,
		progress: progress,
		exit:     make(chan struct{}, 1), // Buffered so that Close() is non-blocking
	}
	go ret.print()
	return ret
}

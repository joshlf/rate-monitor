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

	"github.com/synful/rate"
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

	r := rate.MakeMonitorReaderFunc(os.Stdin, 0, rateFn(size, *unit, *progress))
	defer r.Close()
	_, err := io.Copy(os.Stdout, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "io error: %v\n", err)
		os.Exit(EXIT_IO)
	}
}

func rateFn(size int, unit string, progress bool) func(r rate.Rate) {
	return func(r rate.Rate) {
		fmt.Fprintf(os.Stderr, "\r%8.4f %s/s\033[K", r.Rate/float64(size), unit)

		if progress {
			size, unit := 0, ""
			// If the rate is displayed in a base-10
			// unit, do the same for the progress,
			// and the same for base-2.
			if strings.Contains(unit, "i") {
				switch {
				case r.Total < 1024:
					size, unit = 1, "B"
				case r.Total < 1024*1024:
					size, unit = 1024, "KiB"
				case r.Total < 1024*1024*1024:
					size, unit = 1024*1024, "MiB"
				default:
					size, unit = 1024*1024*1024, "GiB"
				}
			} else {
				switch {
				case r.Total < 1000:
					size, unit = 1, "B"
				case r.Total < 1000*1000:
					size, unit = 1000, "KB"
				case r.Total < 1000*1000*1000:
					size, unit = 1000*1000, "MB"
				default:
					size, unit = 1000*1000*1000, "GB"
				}
			}
			fmt.Fprintf(os.Stderr, " (%.4f %s total)\033[K", float64(r.Total)/float64(size), unit)
		}
	}
}

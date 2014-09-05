// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/synful/rate"
)

const (
	EXIT_IO = 2
)

const usage = `Usage: rate-monitor [-p] [--B | --kB | --KiB | --MB | --MiB | --GB | --GiB]

  --B           show the rate in B/s
  --kB          show the rate in kB/s
  --KiB         show the rate in KiB/s (this is the default)
  --MB          show the rate in MB/s
  --MiB         show the rate in MiB/s
  --GB          show the rate in GB/s
  --GiB         show the rate in GiB/s
  -p --progress show the total number of bytes copied in addition to the rate`

var units = map[string]int{
	"B":   1,
	"kB":  1000,
	"KiB": 1024,
	"MB":  1000 * 1000,
	"MiB": 1024 * 1024,
	"GB":  1000 * 1000 * 1000,
	"GiB": 1024 * 1024 * 1024,
}

func main() {
	args, err := docopt.Parse(usage, nil, true, "", false)
	unit, size := "KiB", units["KiB"]
	for unt, sze := range units {
		if args["--"+unt].(bool) {
			unit, size = unt, sze
			break
		}
	}

	r := rate.MakeMonitorReaderFunc(os.Stdin, 0, rateFn(size, unit, args["-p"].(bool)))
	defer r.Close()
	_, err = io.Copy(os.Stdout, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "io error: %v\n", err)
		os.Exit(EXIT_IO)
	}
}

func rateFn(size int, unit string, progress bool) func(r rate.Rate) {
	return func(r rate.Rate) {
		fmt.Fprintf(os.Stderr, "\r%8.4f %s/s\033[K", r.Rate/float64(size), unit)

		if progress {
			sizeTmp, unitTmp := 0, ""

			// If the rate is displayed in an SI unit,
			// then display the progress in SI units,
			// and likewise for IEC units.
			if strings.Contains(unit, "i") || unit == "B" {
				switch {
				case r.Total < 1024:
					sizeTmp, unitTmp = 1, "B"
				case r.Total < 1024*1024:
					sizeTmp, unitTmp = 1024, "KiB"
				case r.Total < 1024*1024*1024:
					sizeTmp, unitTmp = 1024*1024, "MiB"
				default:
					sizeTmp, unitTmp = 1024*1024*1024, "GiB"
				}
			} else {
				switch {
				case r.Total < 1000:
					sizeTmp, unitTmp = 1, "B"
				case r.Total < 1000*1000:
					sizeTmp, unitTmp = 1000, "kB"
				case r.Total < 1000*1000*1000:
					sizeTmp, unitTmp = 1000*1000, "MB"
				default:
					sizeTmp, unitTmp = 1000*1000*1000, "GB"
				}
			}
			fmt.Fprintf(os.Stderr, " (%.4f %s total)\033[K", float64(r.Total)/float64(sizeTmp), unitTmp)
		}
	}
}

<!--
Copyright 2014 The Authors. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE fil
-->

rate-monitor
============

[![Build Status](https://travis-ci.org/joshlf13/rate-monitor.svg?branch=master)](https://travis-ci.org/joshlf13/rate-monitor)

A simple rate monitoring command-line utility. Stdin is piped to stdout, and the rate at which data is being read is displayed on stderr.

Usage:
```
Usage: rate-monitor [-p] [--B | --kB | --KiB | --MB | --MiB | --GB | --GiB]

  --B           show the rate in B/s
  --kB          show the rate in kB/s
  --KiB         show the rate in KiB/s (this is the default)
  --MB          show the rate in MB/s
  --MiB         show the rate in MiB/s
  --GB          show the rate in GB/s
  --GiB         show the rate in GiB/s
  -p --progress show the total number of bytes copied in addition to the rate
```

Installation: `go get github.com/joshlf13/rate-monitor`
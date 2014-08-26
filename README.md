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
Usage of rate-monitor:
  -progress=false: Display the total amount of data copied so far.
  -unit="KB": Unit for displaying rate. Options are B, KB, KiB, MB, MiB, GB, GiB.
```

Installation: `go get github.com/joshlf13/rate-monitor`
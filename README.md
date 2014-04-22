<!--
Copyright 2014 The Authors. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE fil
-->

rate-monitor
============

[![Build Status](https://travis-ci.org/joshlf13/rate-monitor.svg?branch=master)](https://travis-ci.org/joshlf13/rate-monitor)

A simple rate monitoring command-line utility. Stdin is piped to stdout, and the rate at which data is being read is displayed on stderr.

Usage: `rate-monitor [-B | -KB | -KiB | -MB | -MiB]`
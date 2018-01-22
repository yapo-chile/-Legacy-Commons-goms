#!/bin/bash

set -e

 # Probably there's a better way to fetch our dependencies

go get github.com/axw/gocov/gocov                   # Coverage reporting tool
go get github.com/AlekSi/gocov-xml                  # Generate XML output in Cobertura format
go get github.com/jstemmer/go-junit-report          # Converts go test output to an xml repor
go get github.com/alecthomas/gometalinter           # Concurrently run Go lint tools and normalise their output

go get github.com/Yapo/goutils
go get github.com/Yapo/logger
go get gopkg.in/facebookgo/inject.v0
go get gopkg.in/facebookgo/pidfile.v0
go get gopkg.in/facebookgo/atomicfile.v0
go get gopkg.in/gorilla/mux.v1
go get gopkg.in/stretchr/testify.v1/assert
go get gopkg.in/stretchr/testify.v1/mock

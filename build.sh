#!/bin/bash
export GOBIN=/usr/local/go/bin
export GOARCH=amd64
export GOROOT=/usr/local/go
export GOOS=linux

SCRIPTDIR=`dirname $0`
SCRIPTDIR=`cd $SCRIPTDIR; pwd`
GOPATH=`cd ${SCRIPTDIR}/../../; pwd`
export GOPATH

go build -o ucloganalyzer main.go
go build -ldflags "-s -w" -o ucloganalyzer_simple main.go
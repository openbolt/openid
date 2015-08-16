#!/bin/bash

BASE=`readlink -f $(dirname $0)`
cd $BASE
go run -race bindata.go main.go

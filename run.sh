#!/bin/sh
go clean -cache
go mod tidy
go clean -modcache
air *.go

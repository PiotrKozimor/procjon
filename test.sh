#!/usr/bin/env bash

set -e
echo "" > coverage.txt

go test -race -coverprofile=profile.out -covermode=atomic -v procjon/*.go
cat profile.out >> coverage.txt
go test -race -coverprofile=profile.out -covermode=atomic -v procjonagent/*.go
cat profile.out >> coverage.txt
go test -coverprofile=profile.out -covermode=atomic -v cmd/server/*.go
cat profile.out >> coverage.txt
go test -race -coverprofile=profile.out -covermode=atomic -v cmd/elastic/*.go
cat profile.out >> coverage.txt
go test -race -coverprofile=profile.out -covermode=atomic -v cmd/systemd/*.go
cat profile.out >> coverage.txt
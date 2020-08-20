#!/usr/bin/env bash

set -ex
echo "" > coverage.txt

go test -race -coverprofile=profile.out -covermode=atomic -v procjon/service*.go
cat profile.out >> coverage.txt
go test -race -coverprofile=profile.out -covermode=atomic -v procjon/slack*.go
cat profile.out >> coverage.txt
go test -race -coverprofile=profile.out -covermode=atomic -v procjon/status*.go
cat profile.out >> coverage.txt
go test -coverprofile=profile.out -covermode=atomic -v procjon/*.go
cat profile.out >> coverage.txt
go test -race -coverprofile=profile.out -covermode=atomic -v procjonagent/*.go
cat profile.out >> coverage.txt
go test -race -coverprofile=profile.out -covermode=atomic -v cmd/elastic/*.go
cat profile.out >> coverage.txt
go test -race -coverprofile=profile.out -covermode=atomic -v cmd/systemd/*.go
cat profile.out >> coverage.txt

sudo -E go test -race -coverprofile=profile.out -covermode=atomic -v cmd/ping/*.go
cat profile.out >> coverage.txt
export SKIP_HANDLE_MONITOR=true
export SKIP_ELASTIC=true
go test -coverprofile=profile.out -covermode=atomic -v ./...
cat profile.out >> coverage.txt
export SKIP_ELASTIC=false
export SKIP_SYSTEMD=true
go test -coverprofile=profile.out -covermode=atomic -v ./...
cat profile.out >> coverage.txt

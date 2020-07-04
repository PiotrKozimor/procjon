#!/usr/bin/env bash

set -e
echo "" > coverage.txt

go test -race -coverprofile=profile.out -covermode=atomic -v procjon/*
cat profile.out >> coverage.txt
go test -race -coverprofile=profile.out -covermode=atomic -v procjonagent/*
cat profile.out >> coverage.txt
go test -coverprofile=profile.out -covermode=atomic -v cmd/server/*
cat profile.out >> coverage.txt
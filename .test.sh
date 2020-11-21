#!/usr/bin/env bash

set -ex

cd procjon
go build -race
cd ../procjonagent
go build -race
cd ..

echo "" > coverage.txt

go test -coverprofile=profile.out
cat profile.out >> coverage.txt
cd sender
go test -coverprofile=profile.out
cat profile.out >> ../coverage.txt
cd ../agent
go test -coverprofile=profile.out
cat profile.out >> ../coverage.txt


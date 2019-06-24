#!/bin/sh
git submodule update --remote
go run sbuild-gen.go

#!/bin/sh
git submodule update --remote
go run index-gen.go

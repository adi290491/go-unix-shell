#!/bin/sh
# Run the shell

cd "$(dirname "$0")"
go run app/*.go "$@"
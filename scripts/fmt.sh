#!/usr/bin/env bash
find . -type f -name \*.go | sed 's%/[^/]*$%%' | sort -u | while read d; do
  go fmt "$d"/*.go
done

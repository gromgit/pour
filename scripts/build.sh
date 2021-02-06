#!/usr/bin/env bash
basemod=$(grep ^module go.mod | cut -d\  -f2)
specs=(linux_amd64 linux_386 darwin_amd64)

build_all() {
  for spec in "${specs[@]}"; do
    if [[ $spec =~ ([^_]*)_(.*) ]]; then
      d=build/"${spec}"; mkdir -p "$d"
      export GOOS=${BASH_REMATCH[1]} GOARCH=${BASH_REMATCH[2]}
      for i in cmd/*; do
        go build -o "${d}/${i##*/}" "${basemod}/${i}"
      done
    fi
  done
}

if [[ $1 == "all" ]]; then
  build_all
else
  for i in cmd/*; do
    go install "${basemod}/${i}"
  done
fi

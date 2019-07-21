#! /bin/bash

set -e
set -x

go test -v -race
golint
go vet

executable="vimwasm-try-plugin"

rm -rf release
gox -arch 'amd64' -os 'linux darwin windows freebsd openbsd netbsd'
mkdir -p release
mv ${executable}_* release/
cd release
for bin in *; do
    if [[ "$bin" == *windows* ]]; then
        command="${executable}.exe"
    else
        command="${executable}"
    fi
    mv "$bin" "$command"
    zip "${bin}.zip" "$command"
    rm "$command"
done

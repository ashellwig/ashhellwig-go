#!/usr/bin/env bash

set -euo pipefail

cover_mode=${COVER_MODE:-coveralls}
cover_dir=$(mktemp -d /tmp/coverage.XXXXXXXXXX)
profile="${cover_dir}/cover.out"

pushd
hash goveralls 2> /dev/null || go get github.com/mattn/goveralls
popd

function generate_cover_data() {
    for d in $(go list ./...); do
        (
            local output="${cover_dir}/${d//\//-}.cover"
            go test -coverprofile="${output}" -covermode="$cover_mode" "$d"
        )
    done

    echo "mode: $cover_mode" > "$profile"
    grep -h -v "^mode:" "$cover_dir"/*.cover >> "$profile"
}

function push_to_coveralls() {
    goveralls -coverprofile="${profile}" -service=travis-ci
}

generate_cover_data
go tool cover -func "${profile}"

case "${1-}" in
    --html)
        go tool cover -html "${profile}"
        ;;
    --coveralls)
        push_to_coveralls
        ;;
esac

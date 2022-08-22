#!/usr/bin/env bash
set -e

# tests

tests() {
    go test -timeout 3000s ./...
}

gcp_build() {
    gcloud builds submit --config=devtools/gcp/build.json
}

function usage {
    echo "Commands:"
    echo "-"
    echo ""
    echo "Builds:"
    echo "  gcp:build               build the project and store on gcp artifact"
    echo "Tests:"
    echo "  tests                   run all tests"
    echo ""
    exit 1
}

main() {
    

    case "$1" in
        "-h" | "--help" | "help")
            usage
            exit 0
            ;;
        tests) tests;;
        gcp:build) gcp_build;;
        *)
            usage
            exit 1
    esac
}

main $@
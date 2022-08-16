#!/usr/bin/env bash
set -e

# tests

tests() {
    go test -timeout 300s ./...
}

function usage {
    echo "Commands:"
    echo "-"
    echo ""
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
        *)
            usage
            exit 1
    esac
}

main $@
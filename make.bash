#!/usr/bin/env bash
set -e

# tests

tests() {
    go test -timeout 3000s ./...
}

local_run() {
    docker run --env-file ./.env -p 8080:8080 -v "$(pwd)"/credentials.json:/app/credentials.json signer_app
}

gcp_cmds() {
    gcloud builds submit --config=devtools/gcp/$1.json
}

function usage {
    echo "Avareum's fund operation signing modules"
    echo ""
    echo "Commands:"
    echo "  gcp:build               build the project and store on gcp artifact"
    echo "  gcp:deploy              deploy on a new gcp cloud compute"
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
        gcp:build) gcp_cmds build;;
        gcp:deploy) gcp_cmds deploy;;
        gcp:reset) gcp_cmds reset;;
        build) docker build -t signer_app .;;
        run) local_run;;
        *)
            usage
            exit 1
    esac
}

main $@
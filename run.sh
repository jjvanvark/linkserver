#!/usr/bin/env bash

case $1 in
    dev)    clear && reflex -r '\.go$' -s -- sh -c "go run ./cmd/server/*.go" ;;
	build)  go build ./cmd/server/lor-server.go ;;
    log)    clear && tail -F default.log ;;
    *)      echo "Invalid task: $1"; exit 1 ;;
esac

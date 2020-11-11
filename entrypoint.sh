#!/bin/sh

set -e

export LC_ALL="C.UTF-8"
export LANG="C.UTF-8"

if [ -n "$EXTRA_PATH" ]; then
    export PATH="$EXTRA_PATH:$PATH"
fi

BIN="$(dirname "$(realpath "$0")")"
ROOT="$(dirname "$BIN")"

case "$1" in
    *)
        exec ./main "$@"
    ;;
esac

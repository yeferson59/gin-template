#!/bin/sh
make fmt
make lint
make test
if [ $? -ne 0 ]; then
echo "Pre-commit checks failed. Commit aborted."
exit 1
fi

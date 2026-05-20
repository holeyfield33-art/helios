#!/bin/bash
# quality_gate.sh - Fail fast on any test error or warning-class issue.

set -euo pipefail

ROOT_DIR="$(git rev-parse --show-toplevel)"
cd "$ROOT_DIR"

echo "[quality-gate] Running Go vet"
go vet ./...

echo "[quality-gate] Running Go tests"
go test ./...

echo "[quality-gate] Running Python tests with warnings treated as errors"
python -m pytest implementations/python/tests -W error -q

echo "[quality-gate] PASS"
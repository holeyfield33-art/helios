#!/bin/bash
# install_hooks.sh - Configure this repository to use managed local hooks.

set -euo pipefail

ROOT_DIR="$(git rev-parse --show-toplevel)"
cd "$ROOT_DIR"

git config core.hooksPath .githooks
echo "Configured git hooks path: .githooks"
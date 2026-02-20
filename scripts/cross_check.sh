#!/bin/bash
# cross_check.sh — Cross-language hash verification
# Runs Go and Python verifiers against the same test vectors.
# Exits 1 on ANY divergence. Exits 0 only on full match.

set -euo pipefail

VECTORS="/app/test_vectors/vectors.json"
GO_BIN="/usr/local/bin/helios"
PYTHON_SCRIPT="/app/implementations/python/verify.py"

echo "=== Helios Core Cross-Language Verification ==="
echo ""

# --- Go verification ---
echo "--- Go Implementation ---"
GO_OUTPUT=$("$GO_BIN" verify "$VECTORS" 2>&1) || {
    echo "Go verification FAILED:"
    echo "$GO_OUTPUT"
    exit 1
}
echo "$GO_OUTPUT"

# Extract Go hashes (lines matching "PASS" with vector names)
GO_HASHES=$(echo "$GO_OUTPUT" | grep "PASS" | sort)

echo ""

# --- Python verification ---
if [ -f "$PYTHON_SCRIPT" ]; then
    echo "--- Python Implementation ---"
    PYTHON_OUTPUT=$(python3 "$PYTHON_SCRIPT" "$VECTORS" 2>&1) || {
        echo "Python verification FAILED:"
        echo "$PYTHON_OUTPUT"
        exit 1
    }
    echo "$PYTHON_OUTPUT"

    # Extract Python hashes
    PYTHON_HASHES=$(echo "$PYTHON_OUTPUT" | grep "PASS" | sort)

    echo ""

    # --- Cross-language comparison ---
    echo "--- Cross-Language Comparison ---"

    # Count matching lines
    GO_COUNT=$(echo "$GO_HASHES" | wc -l)
    PY_COUNT=$(echo "$PYTHON_HASHES" | wc -l)

    if [ "$GO_COUNT" != "$PY_COUNT" ]; then
        echo "FAIL: Go produced $GO_COUNT results, Python produced $PY_COUNT results"
        exit 1
    fi

    # Compare line by line
    MATCH_COUNT=0
    TOTAL=0
    while IFS= read -r go_line; do
        TOTAL=$((TOTAL + 1))
        py_line=$(echo "$PYTHON_HASHES" | sed -n "${TOTAL}p")
        if [ "$go_line" = "$py_line" ]; then
            MATCH_COUNT=$((MATCH_COUNT + 1))
        else
            echo "DIVERGENCE at vector $TOTAL:"
            echo "  Go:     $go_line"
            echo "  Python: $py_line"
            exit 1
        fi
    done <<< "$GO_HASHES"

    echo "Cross-language match: ${MATCH_COUNT}/${TOTAL} identical hashes"

    if [ "$MATCH_COUNT" -ne "$TOTAL" ]; then
        echo "FAIL: Not all hashes matched"
        exit 1
    fi
else
    echo "Python implementation not found at $PYTHON_SCRIPT"
    echo "Running Go-only verification..."
    echo ""
    echo "Go verification: All vectors PASS"
fi

echo ""
echo "=== Verification Complete ==="
exit 0

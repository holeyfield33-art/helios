"""Test vector verification for Helios Core (Python conformance)."""

import json
import sys
from typing import Any

from conformance.hasher import content_hash
from conformance.objects import MemoryObject, Relationship
from conformance.canon import validate_ingest_value


def load_vectors(path: str) -> list:
    """Load test vectors from a JSON file."""
    with open(path) as f:
        data = json.load(f)
    return data["vectors"]


def input_to_memory_object(inp: dict) -> MemoryObject:
    """Convert a raw JSON dict to a MemoryObject.
    Validates ingest rules: RULE-002 (no floats), RULE-009 (integer range), RULE-010 (no nulls).
    """
    # Ingest validation on the value field
    validate_ingest_value(inp.get("value"), "value")

    relationships = []
    for r in inp.get("relationships", []):
        relationships.append(Relationship(key=r["key"], type=r["type"]))

    return MemoryObject(
        category=inp.get("category", ""),
        created_at=inp.get("created_at", ""),
        key=inp.get("key", ""),
        relationships=relationships,
        source=inp.get("source", ""),
        value=inp.get("value"),
    )


def verify_vectors(path: str) -> list:
    """Verify all test vectors. Returns list of (name, expected, got, pass) tuples.
    Raises SystemExit(1) if any vector fails.
    """
    vectors = load_vectors(path)
    results = []
    failures = 0

    for vec in vectors:
        vector_id = vec["vector_id"]
        vector_type = vec.get("vector_type", "positive")
        expected_outcome = vec.get("expected_outcome", "ACCEPT")

        if vector_type == "negative":
            # Negative vectors: expect rejection
            rejection_code = vec.get("rejection_code", "")
            try:
                obj = input_to_memory_object(vec["input"])
                got = content_hash(obj)
                # Should have been rejected but wasn't
                results.append((vector_id, "REJECT", f"ACCEPT: {got}", False))
                failures += 1
            except (ValueError, Exception) as e:
                error_msg = str(e)
                passed = rejection_code and rejection_code in error_msg
                results.append((vector_id, "REJECT", error_msg, passed))
                if not passed:
                    failures += 1
        else:
            # Positive vectors: expect successful hash match
            expected_hash = vec["hash"]
            obj = input_to_memory_object(vec["input"])
            got = content_hash(obj)
            passed = got == expected_hash
            results.append((vector_id, expected_hash, got, passed))
            if not passed:
                failures += 1

    return results, failures


def main():
    if len(sys.argv) < 2:
        print("Usage: python -m conformance.verifier <vectors.json>", file=sys.stderr)
        sys.exit(1)

    path = sys.argv[1]
    results, failures = verify_vectors(path)

    for name, expected, got, passed in results:
        status = "PASS" if passed else "FAIL"
        print(f"  {name}: {status}")
        if not passed:
            print(f"    expected: {expected}")
            print(f"    got:      {got}")

    if failures > 0:
        print(f"\n{failures} of {len(results)} vectors FAILED", file=sys.stderr)
        sys.exit(1)

    print(f"\nAll {len(results)} vectors: PASS")


if __name__ == "__main__":
    main()

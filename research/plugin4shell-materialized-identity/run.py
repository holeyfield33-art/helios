"""Local research fixture; no network, execution payload, or production integration."""

import json
from pathlib import Path
import sys

HERE = Path(__file__).resolve().parent
sys.path.insert(0, str(HERE.parents[1] / "implementations" / "python"))

from helios import hash_memory_object  # noqa: E402


def load(name):
    return json.loads((HERE / "fixture" / name).read_text(encoding="utf-8"))


def require(condition, message):
    # Explicit exceptions keep the experiment falsifiable under python -O too.
    if not condition:
        raise ValueError(message)


def run():
    expected = load("expected.json")
    substituted = load("substituted.json")
    trusted = load("expected-identity.json")
    expected_hash = trusted["sha256"]
    expected_reference = trusted["reference"]

    def resolve(reference, artifact):
        require(reference == expected_reference, "unexpected requested reference")
        return artifact

    def naive_accepts(artifact):
        return artifact["reference"] == expected_reference

    def verified_accepts(artifact):
        # Use the existing Core implementation, including its schema version 1.
        return hash_memory_object(artifact["object"]) == expected_hash

    original = resolve(expected_reference, expected)
    actual = resolve(expected_reference, substituted)
    require(naive_accepts(original) and verified_accepts(original),
            "positive control rejected the original object")
    require(original["object"] != actual["object"], "fixture objects must differ")
    require(naive_accepts(actual), "naive path did not accept substitution")
    require(not verified_accepts(actual), "verified path accepted substitution")
    reordered = dict(reversed(list(original["object"].items())))
    require(hash_memory_object(reordered) == expected_hash,
            "key ordering changed canonical identity")
    wrong_reference = {**original, "reference": "memory://other"}
    require(not naive_accepts(wrong_reference), "naive negative control failed")
    require(verified_accepts(wrong_reference),
            "content identity unexpectedly authenticated the external reference")
    print("naive reference check: ACCEPTS substituted object")
    print("post-materialization content verification: REJECTS substituted object")
    print(f"expected SHA-256: {expected_hash}")
    print(f"actual SHA-256:   {hash_memory_object(actual['object'])}")
    print("PASS: reference equality does not imply object equality")


if __name__ == "__main__":
    try:
        run()
    except Exception as error:
        print(f"FAIL: {error}", file=sys.stderr)
        sys.exit(1)

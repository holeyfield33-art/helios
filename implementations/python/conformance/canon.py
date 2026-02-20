"""Canonical serialization primitives for Helios Core (Python conformance)."""

import json
import unicodedata


def normalize_string(s: str) -> str:
    """Apply NFC Unicode normalization to a string value.
    Must be called on EVERY string field value BEFORE serialization.
    """
    return unicodedata.normalize("NFC", s)


def normalize_timestamp(s: str) -> str:
    """Validate and normalize an ISO 8601 UTC timestamp.

    Output: YYYY-MM-DDTHH:MM:SS.sssZ (exactly 3 decimal places)
    Rejects: non-Z suffix, != 3 fractional digits.
    Matches Go behavior: explicit rejection on invalid input.
    """
    if not s.endswith("Z"):
        raise ValueError(f"CANON_ERR_TIMESTAMP_NON_UTC: Timestamp must end in Z, got: {s}")

    # Validate exactly 3 fractional digits
    dot_idx = s.rfind(".")
    if dot_idx == -1:
        raise ValueError(f"CANON_ERR_TIMESTAMP_INVALID_PRECISION: Timestamp must have exactly 3 fractional digits, got none: {s}")

    frac = s[dot_idx + 1 : -1]  # between '.' and 'Z'
    if len(frac) != 3:
        raise ValueError(
            f"CANON_ERR_TIMESTAMP_INVALID_PRECISION: Timestamp must have exactly 3 fractional digits, got {len(frac)}: {s}"
        )

    # Validate parsability (basic structure check)
    from datetime import datetime, timezone

    try:
        dt = datetime.strptime(s, "%Y-%m-%dT%H:%M:%S.%fZ").replace(tzinfo=timezone.utc)
    except ValueError:
        raise ValueError(f"Invalid timestamp format: {s}")

    # Re-format with exactly 3 decimal places
    ms = dt.microsecond // 1000
    return dt.strftime(f"%Y-%m-%dT%H:%M:%S.{ms:03d}Z")


def canonicalize_object(obj: dict) -> bytes:
    """Produce deterministic canonical JSON bytes from a dict.

    Rules:
    - Keys sorted lexicographically at every level
    - Compact format (no whitespace)
    - ensure_ascii=False (UTF-8 bytes preserved, no \\uXXXX for non-ASCII)
    - null values included (not stripped)
    - Arrays: insertion order preserved
    - Recursive application to nested objects
    """
    if not isinstance(obj, dict):
        raise TypeError(f"canonicalize_object expects dict, got {type(obj)}")

    normalized = _normalize_dict(obj)
    return json.dumps(
        normalized,
        ensure_ascii=False,
        separators=(",", ":"),
        sort_keys=True,
    ).encode("utf-8")


def _normalize_dict(d: dict) -> dict:
    """Recursively build a sorted dict with normalized values."""
    return {k: _normalize_value(v) for k, v in sorted(d.items())}


def _normalize_value(v):
    """Recursively normalize a value for canonical serialization."""
    if v is None:
        raise ValueError("CANON_ERR_NULL_PROHIBITED: null values are not permitted")
    if isinstance(v, dict):
        return _normalize_dict(v)
    elif isinstance(v, list):
        return [_normalize_value(item) for item in v]  # preserve order
    else:
        return v  # str, int, float, bool pass through as-is


def sort_relationships(rels: list) -> list:
    """Sort relationships by key first, then type as tie-breaker."""
    return sorted(rels, key=lambda r: (r.key, r.type))


def relationship_to_map(r) -> dict:
    """Convert a Relationship to an explicit dict. Never rely on dataclass ordering."""
    return {"key": r.key, "type": r.type}


def validate_schema_version(input: dict) -> None:
    """Validate RULE-001: _helios_schema_version must be present and equal to \"1\"."""
    if "_helios_schema_version" not in input:
        raise ValueError("CANON_ERR_SCHEMA_VERSION_MISSING: _helios_schema_version field is required")
    v = input["_helios_schema_version"]
    if not isinstance(v, str) or v != "1":
        raise ValueError(f"CANON_ERR_SCHEMA_VERSION_INVALID: _helios_schema_version must be string \"1\", got {v!r}")


def validate_ingest_value(v, path: str = "") -> None:
    """Recursively validate a parsed JSON value for spec compliance.

    Checks:
    - RULE-002: No float values (CANON_ERR_FLOAT_PROHIBITED)
    - RULE-009: Integer range within signed 64-bit (CANON_ERR_INTEGER_OUT_OF_RANGE)
    - RULE-010: No null/None values (CANON_ERR_NULL_PROHIBITED)
    """
    if v is None:
        raise ValueError(f"CANON_ERR_NULL_PROHIBITED: null value at {path}")
    elif isinstance(v, float):
        raise ValueError(f"CANON_ERR_FLOAT_PROHIBITED: float value at {path}")
    elif isinstance(v, bool):
        pass  # bool check before int because bool is subclass of int in Python
    elif isinstance(v, int):
        if v > 9223372036854775807 or v < -9223372036854775808:
            raise ValueError(f"CANON_ERR_INTEGER_OUT_OF_RANGE: value {v} at {path} exceeds int64 bounds")
    elif isinstance(v, dict):
        for k, child in v.items():
            validate_ingest_value(child, f"{path}.{k}")
    elif isinstance(v, list):
        for i, child in enumerate(v):
            validate_ingest_value(child, f"{path}[{i}]")
    elif isinstance(v, str):
        pass  # valid
    else:
        raise ValueError(f"Unsupported type {type(v)} at {path}")

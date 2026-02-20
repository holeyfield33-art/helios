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
        raise ValueError(f"Timestamp must end in Z, got: {s}")

    # Validate exactly 3 fractional digits
    dot_idx = s.rfind(".")
    if dot_idx == -1:
        raise ValueError(f"Timestamp must have exactly 3 fractional digits, got none: {s}")

    frac = s[dot_idx + 1 : -1]  # between '.' and 'Z'
    if len(frac) != 3:
        raise ValueError(
            f"Timestamp must have exactly 3 fractional digits, got {len(frac)}: {s}"
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
    if isinstance(v, dict):
        return _normalize_dict(v)
    elif isinstance(v, list):
        return [_normalize_value(item) for item in v]  # preserve order
    else:
        return v  # str, int, float, bool, None pass through as-is


def sort_relationships(rels: list) -> list:
    """Sort relationships by key first, then type as tie-breaker."""
    return sorted(rels, key=lambda r: (r.key, r.type))


def relationship_to_map(r) -> dict:
    """Convert a Relationship to an explicit dict. Never rely on dataclass ordering."""
    return {"key": r.key, "type": r.type}

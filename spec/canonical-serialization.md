# Helios Core — Canonical Serialization Specification

**Version:** 1.0  
**Date:** 2026-02-19  
**Authors:** Claude (Anthropic) & ChatGPT (OpenAI) via TMRP  

## 1. Purpose

This specification defines a deterministic serialization format for AI memory objects, enabling a verifiable content hash that proves an object has not been modified. The canonical form ensures identical inputs always produce identical byte sequences, regardless of implementation language.

## 2. Scope

This specification covers:
- Canonical JSON serialization rules
- Unicode normalization requirements
- Timestamp format requirements
- Hash input field selection
- Relationship ordering rules

## 3. Canonical JSON Rules

### 3.1 Key Ordering
All object keys MUST be sorted lexicographically (Unicode code point order) at every nesting level.

### 3.2 Whitespace
The canonical form MUST use compact JSON with no whitespace between tokens. No spaces after colons, no spaces after commas, no newlines.

### 3.3 Null Values
Null values MUST be included in the serialized output. A field with value `null` serializes as `"field":null`. Fields MUST NOT be omitted when their value is null.

### 3.4 Arrays
Array elements MUST preserve their insertion order. Arrays MUST NOT be sorted unless explicitly specified (e.g., relationships).

### 3.5 UTF-8 Encoding
All string values MUST be serialized as raw UTF-8 bytes. Non-ASCII characters MUST NOT be escaped to `\uXXXX` form. Only characters required by the JSON specification to be escaped (control characters, backslash, double quote) are escaped.

### 3.6 Empty Arrays
Empty arrays serialize as `[]`, not `null`. They are included in the hash input.

## 4. Unicode Normalization

All string field VALUES MUST be normalized to NFC (Unicode Normalization Form C) BEFORE serialization. This ensures that equivalent Unicode representations (e.g., precomposed vs. decomposed characters) produce identical canonical bytes.

Normalization applies to:
- `category`
- `key`
- `source`
- `value` (when the value is a string)
- Relationship `key` and `type` fields

Normalization MUST occur on input values, NOT on output bytes.

## 5. Timestamp Format

All timestamps MUST conform to this exact format:

```
YYYY-MM-DDTHH:MM:SS.sssZ
```

Requirements:
- MUST end with `Z` (UTC only)
- MUST have EXACTLY 3 fractional second digits (millisecond precision)
- Timestamps with fewer or more fractional digits MUST be rejected
- Implementations MUST NOT use variable-precision parsers (e.g., Go's `time.RFC3339Nano`)

## 6. Float Representation

Float values in test vectors are chosen such that their shortest round-trip decimal representation is identical across conformant implementations. Any float value whose canonical string form cannot be independently verified to be identical across implementations is outside v1 scope.

## 7. Hash Input Construction

### 7.1 Included Fields (exactly 6)

| Field | Type | Description |
|-------|------|-------------|
| `category` | string | Object category |
| `created_at` | string | Creation timestamp (canonical format) |
| `key` | string | Unique object key |
| `relationships` | array | Sorted relationship objects |
| `source` | string | Origin source |
| `value` | any | Object value (may be null) |

### 7.2 Excluded Fields

The following fields are NOT included in the content hash:
- `updated_at`
- `version`
- `access_count`
- `last_accessed`
- `confidence`

### 7.3 Construction Steps

1. Extract only the 6 included fields
2. Normalize the timestamp to canonical format
3. Sort relationships by `key`, then `type` as tie-breaker
4. Apply NFC normalization to all string values
5. Build a map with exactly 6 keys
6. Apply canonical serialization
7. Compute SHA-256 hash of the canonical bytes
8. Encode as lowercase hexadecimal (64 characters)

## 8. Relationship Canonicalization

Each relationship is an object with exactly two fields: `key` and `type`.

### 8.1 Sorting
Relationships MUST be sorted by `key` (lexicographic), with `type` as a tie-breaker when keys are equal.

### 8.2 Serialization
Each relationship MUST be serialized as an explicit map with sorted keys: `{"key":"...","type":"..."}`. Implementations MUST NOT rely on struct field ordering.

## 9. Content Hash Algorithm

```
canonical_bytes = canonicalize(hash_input_map)
content_hash = hex(sha256(canonical_bytes))
```

The content hash is a 64-character lowercase hexadecimal string representing the SHA-256 digest of the canonical JSON bytes.

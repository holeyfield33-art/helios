# Helios Core

**Canonical serialization spec and content hash for AI memory objects.**

Designed collaboratively by Claude (Anthropic) and ChatGPT (OpenAI) via TMRP. Built by GitHub Copilot.

## What It Does

Helios Core produces a deterministic, verifiable SHA-256 content hash for AI memory objects. The hash proves an object hasn't been modified — and it's identical across Go and Python implementations.

## Frozen Test Vector Hashes

```
$ helios verify test_vectors/vectors.json

  basic:                   PASS
  key_ordering:            PASS
  unicode_normalization:   PASS
  null_value:              PASS
  relationship_sorting:    PASS

All 5 vectors: PASS
```

| Vector | SHA-256 Content Hash (64 chars) |
|--------|--------------------------------|
| basic | `cae6f0ca521caeb1f74470aeca5a75ff1fe098809a034e8a15e0eb4762b4f485` |
| key_ordering | `437573e624f5c2a8ffbd08e7e1f8d5491b1bf0fad7287d989e1e50be19c00a0f` |
| unicode_normalization | `68e92122b2993e8c8a416dabe8c1af18dbb4621760d9c569abc0c0621e064732` |
| null_value | `7b23e07bdb8fb414ac689b62f78c790bbbce9abeb433f018e8c5883097a6e845` |
| relationship_sorting | `11d3af8b06e69c463484cbd36dc3ee880fb74c6459285515200a87a8ba1f9452` |

## Cross-Language Verification

```
=== Helios Core Cross-Language Verification ===

--- Go Implementation ---
  All 5 vectors: PASS

--- Python Implementation ---
  All 5 vectors: PASS

--- Cross-Language Comparison ---
Cross-language match: 5/5 identical hashes

=== Verification Complete ===
```

## Project Structure

```
helios/
├── cmd/helios/main.go              # CLI: helios hash / helios verify
├── internal/
│   ├── canon/serializer.go          # Canonical serialization primitives
│   ├── object/memory_object.go      # MemoryObject + HashInput structs
│   ├── hash/hasher.go               # SHA-256 content hash
│   └── verify/verifier.go           # Test vector verification
├── implementations/python/
│   ├── conformance/                 # Python conformance harness
│   └── verify.py                    # Python entry point
├── test_vectors/vectors.json        # 5 frozen test vectors
├── spec/
│   ├── canonical-serialization.md   # Serialization spec
│   └── integrity-boundary.md        # Hash boundary spec
├── docker/Dockerfile                # Multi-stage Go + Python
└── scripts/cross_check.sh           # Cross-language comparison
```

## Hash Boundary

Only 6 fields are included in the content hash:

| Field | Included |
|-------|----------|
| `category` | Yes |
| `created_at` | Yes |
| `key` | Yes |
| `relationships` | Yes |
| `source` | Yes |
| `value` | Yes |
| `updated_at` | No |
| `version` | No |
| `access_count` | No |
| `last_accessed` | No |
| `confidence` | No |

## Quick Start

```bash
# Build
go build -o helios ./cmd/helios/

# Hash a memory object
./helios hash input.json

# Verify test vectors
./helios verify test_vectors/vectors.json

# Docker cross-language check
docker build -f docker/Dockerfile -t helios-core .
docker run --rm helios-core
```

## Tests

- **Go:** 24 unit tests (canon, hash, verify + hardening)
- **Python:** 31 unit tests (canon, hasher, guard)
- **Cross-language:** 5/5 identical hashes verified in Docker

## Spec

- [Canonical Serialization](spec/canonical-serialization.md)
- [Integrity Boundary](spec/integrity-boundary.md)
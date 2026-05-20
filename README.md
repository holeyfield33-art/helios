# Helios Core

![CI](https://github.com/holeyfield33-art/helios/actions/workflows/test.yml/badge.svg)
[![PyPI version](https://img.shields.io/pypi/v/helios-core?logo=pypi&logoColor=white)](https://pypi.org/project/helios-core/)

**Canonical serialization spec and content hash for AI memory objects.**

## What It Does

Helios Core produces a deterministic, verifiable SHA-256 content hash for AI
memory objects. The hash proves an object hasn't been modified — and it's
identical across Go and Python implementations.

## Frozen Test Vector Hashes

```text
$ helios verify test_vectors/vectors.json

  POS-001:                 PASS
  POS-002:                 PASS
  POS-003:                 PASS
  POS-004:                 PASS
  POS-005:                 PASS
  NEG-001:                 PASS
  NEG-002:                 PASS
  NEG-003:                 PASS
  NEG-004:                 PASS
  NEG-005:                 PASS
  NEG-006:                 PASS
  NEG-007:                 PASS
  NEG-008:                 PASS
  NEG-009:                 PASS
  NEG-010:                 PASS
  NEG-011:                 PASS
  NEG-012:                 PASS

All 17 vectors: PASS
```

| Vector | SHA-256 Content Hash (64 chars) |
| ------ | -------------------------------- |
| POS-001 | `c3262407645dcdbd1cede212fa0448a3adb2f915f762540c32e0050bbf65e781` |
| POS-002 | `694cafaa80dd0121a4c4415ac44793fee17104d02756b3c1456dd79fc467c1d0` |
| POS-003 | `d7b4f1c46600c6b7f6733e866455cfa3c5646b6e63625a2107a6a57a36be486c` |
| POS-004 | `5e43f9576ec448e9111856b8e0f95593e4aa427ba9ec71cb3a6b574a91719558` |
| POS-005 | `84c6d544a9ee3b9c1bd48a17d8835f25a7df62cd520f78f12fa49810b9e35945` |
| NEG-001..NEG-012 | Rejection vectors, no expected hash |

## Cross-Language Verification

```text
=== Helios Core Cross-Language Verification ===

--- Go Implementation ---
  All 17 vectors: PASS

--- Python Implementation ---
  All 17 vectors: PASS

--- Cross-Language Comparison ---
Cross-language match: 17/17 identical outcomes

=== Verification Complete ===
```

## Project Structure

```text
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
├── test_vectors/vectors.json        # 17 frozen test vectors
├── spec/
│   ├── canonical-serialization.md   # Serialization spec
│   └── integrity-boundary.md        # Hash boundary spec
├── docker/Dockerfile                # Multi-stage Go + Python
└── scripts/cross_check.sh           # Cross-language comparison
```

## Hash Boundary

Only 6 fields are included in the content hash:

| Field | Included |
| ----- | -------- |
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

# Install local commit/push guards (once per clone)
bash scripts/install_hooks.sh

# Run the full quality gate (errors and warnings fail)
bash scripts/quality_gate.sh

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
- **Cross-language:** 17/17 identical outcomes verified in Docker

## Spec

- [Canonical Serialization](spec/canonical-serialization.md)
- [Integrity Boundary](spec/integrity-boundary.md)

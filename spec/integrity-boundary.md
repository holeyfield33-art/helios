# Helios Core — Integrity Boundary Specification

**Version:** 1.0  
**Date:** 2026-02-19  
**Authors:** Claude (Anthropic) & ChatGPT (OpenAI) via TMRP  

## 1. Purpose

This document defines the integrity boundary for Helios Core memory objects — which fields participate in the content hash and which are explicitly excluded. The boundary ensures that operational metadata changes do not invalidate the content hash, while any modification to the semantic content of a memory object is detectable.

## 2. Integrity Boundary

### 2.1 Inside the Boundary (hashed)

These 6 fields constitute the semantic identity of a memory object:

| Field | Rationale |
|-------|-----------|
| `category` | Defines the object's classification |
| `created_at` | Immutable creation timestamp |
| `key` | Unique identifier |
| `relationships` | Semantic links to other objects |
| `source` | Origin of the memory |
| `value` | The actual content |

A change to ANY of these fields produces a different content hash.

### 2.2 Outside the Boundary (not hashed)

These fields are operational metadata that may change without affecting identity:

| Field | Rationale |
|-------|-----------|
| `updated_at` | Changes on every write |
| `version` | Incremented on updates |
| `access_count` | Changes on every read |
| `last_accessed` | Changes on every access |
| `confidence` | May be adjusted over time |

Modifying these fields MUST NOT affect the content hash.

## 3. Verification Protocol

### 3.1 Hash Verification
Given a memory object and its claimed content hash:
1. Extract the 6 hash-input fields
2. Compute the canonical content hash
3. Compare to the claimed hash
4. If they differ, the object has been modified

### 3.2 Cross-Language Verification
The same memory object MUST produce the same content hash regardless of which conformant implementation computes it. This is verified through:
- Shared test vectors with frozen expected hashes
- Docker-based cross-language comparison
- The `cross_check.sh` script that exits non-zero on any divergence

## 4. Security Considerations

- The content hash does NOT provide authentication — it only detects modification
- The hash is computed over the canonical JSON bytes, not the original input format
- SHA-256 is used for its collision resistance and widespread availability
- The hash is deterministic: same input always produces same output

## 5. Failure Modes

See the TMRP audit decision log for the 11 identified failure modes and their mitigations. Critical failure modes include:
1. Timestamp precision variance
2. Map iteration order non-determinism
3. NFC normalization timing
4. Hash input field leakage
5. Null value omission
6. Relationship sort order incompleteness
7. Float representation divergence

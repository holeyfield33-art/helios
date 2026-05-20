# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] — 2026-02-20

### Added

- **Frozen canonical serialization primitive** with immutable spec_version 1
- **17 frozen test vectors** (5 positive, 12 negative) for cross-language verification
- **Go reference implementation** with canonical serialization and SHA-256 content hashing
- **Python reference implementation** (Python 3.10+) with identical behavior verified against Go
- **Public API** (`helios-core` PyPI package with `hash_memory_object` function)
- **CLI tool** for computing hashes and verifying test vectors
- **GitHub Actions CI** for Go and Python tests on push/PR to main
- **Release automation** for building Go binaries and publishing Python packages
- **Community templates** including issue templates, PR template, CODE_OF_CONDUCT
- **Security policy** with vulnerability reporting guidance

### Specification

- **Spec Version**: 1 (immutable)
- **Vectors Version**: 3 (frozen)
- **Hash Algorithm**: SHA-256 (Applied to canonical JSON representation)
- **Included Fields** (6): category, created_at, key, relationships, source, value
- **Excluded Fields** (5): updated_at, version, access_count, last_accessed, confidence

### Technical Details

- **Hash Boundary**: 6 required fields + optional relationships array
- **Canonical Form**: Lexicographically sorted keys, no whitespace, strict float rejection
- **Test Coverage**: 5 positive vectors with expected hashes, 12 negative vectors testing rejection cases
- **Cross-Language Guarantee**: All test vectors verify identical behavior across Go and Python

### Status

**FROZEN**: This is a trust primitive with immutable spec_version 1. All future changes are additive only (new language implementations, documentation, tooling). The canonical serialization algorithm and test vectors cannot be modified.

---

For the full specification, see [spec/canonical-serialization.md](spec/canonical-serialization.md).
For development setup and contribution guidelines, see [CONTRIBUTING.md](CONTRIBUTING.md).

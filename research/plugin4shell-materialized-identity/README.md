# Materialized identity — research draft

Date: 2026-09-18. Isolated, non-product experiment inspired by the
[ASI research intake](https://github.com/holeyfield33-art/agent-security-index/blob/main/docs/PLUGIN4SHELL-RESEARCH.md).
This is a model of a trust-boundary error, not a reproduction of a coding-agent exploit.

**Proposition: reference equality does not imply object equality.**

Materialized Identity Invariant:
A security decision MUST be bound to the cryptographic identity of the artifact actually materialized for use, not merely to the identifier, reference, commit, locator, or expected digest requested before resolution.

Helios Core currently implements this principle for canonical AI memory objects. It does not authenticate provenance and does not currently verify Git/plugin artifacts.

## Run

From the repository root, with Python 3.10 or later:

```sh
python research/plugin4shell-materialized-identity/run.py
```

No installation, external dependency, network request, Git operation, or executable
payload is needed. The script imports the repository's existing Python public API.
It leaves the shipping CLI and package API untouched.

## Design and falsification

`fixture/expected.json` and `fixture/substituted.json` carry the same external
reference and memory key, but different `value` fields. A local resolver model
returns B for the requested reference. The naive verifier compares only the
returned reference; the content verifier canonicalizes and hashes the actual
memory object through `hash_memory_object` and compares it with the fixed digest
in `fixture/expected-identity.json`.

The expected digest was computed from A before the experiment and checked into
the fixture. The run never regenerates it from resolver output. That manifest is
an assumed trusted input, not authenticated by this experiment.

Expected output includes:

```text
naive reference check: ACCEPTS substituted object
post-materialization content verification: REJECTS substituted object
PASS: reference equality does not imply object equality
```

Every failed condition returns a nonzero exit code and prints FAIL, including
under Python optimization. Positive control: A passes both paths. Additional
controls: reordered object keys preserve the hash; a wrong reference fails the
naive check but still passes content verification for A, exposing the distinction
between integrity and authenticated reference binding. Making B identical to A,
or changing the stored expected digest, falsifies the expected result.

## What the result means

Helios proves:

- deterministic content identity
- modification detection

These statements apply within the canonical six-field integrity boundary, given
a trusted expected digest and SHA-256's collision-resistance assumption. They
do not assert byte-for-byte identity of raw JSON or cover excluded metadata.
The existing implementation also includes its fixed schema-version marker `1`;
this experiment neither replaces nor modifies that behavior.

Helios does not prove:

- who authorized the expected hash
- who published the artifact
- that the reference-to-hash mapping is trustworthy
- that a signer is trusted

If an attacker can replace both the object and the trusted expected digest,
hash comparison alone cannot establish authorization. No provenance signing,
marketplace enforcement, package verification, or runtime integration is built.
The fixture supports the invariant in a small memory-object model; it does not
validate any vendor patch or demonstrate exploitation in the wild.

Core 1.0, spec version 1, the existing frozen vectors, and the six-field boundary
remain unchanged. See [canonical serialization](../../spec/canonical-serialization.md)
and [integrity boundary](../../spec/integrity-boundary.md).

## Validation record — 2026-09-18

- Fixture: PASS with ordinary Python and `python -O`.
- Manual in-memory falsification checks: replacing B with A fails; replacing
  hashing with a constant expected digest also fails. No fixture file was changed
  by these checks.
- Existing Python suite: 47/47 PASS, warnings treated as errors.
- Go 1.26.0 under WSL: `go vet ./...` and `go test ./...` PASS.
- Existing frozen vectors: 17/17 PASS in each implementation. The repository
  cross-language script reports 17 matching outcomes. Its checked-out Windows
  CRLF endings initially prevented Bash execution; an LF-only temporary copy
  beside the original ran successfully and was removed. The tracked script was
  not changed. The comparison checks result lines, not a new byte-level proof.
- `git diff --exit-code` confirms no changes to Core implementations, CLI, spec,
  vectors, dependency manifests, or root README. `git diff --check` passes.

Added files are this README, `run.py`, and the three JSON files under `fixture/`.
The experiment supports the Materialized Identity Invariant within its stated
model. Authenticating the expected hash remains outside the result.

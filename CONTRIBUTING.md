# Contributing to Helios Core

## Scope

Helios is a frozen trust primitive. The canonical serialization spec and the 17
frozen test vectors at `spec_version 1` will not change. Contributions that
modify the spec or test vectors will not be accepted unless they fix a provable
correctness bug in the existing specification.

## What contributions are welcome

* New language implementations that pass all 17 frozen test vectors
* Documentation improvements
* CI / tooling improvements
* Bug reports with reproducible test cases

## Development setup

```bash
# Clone the repository
git clone https://github.com/holeyfield33-art/helios.git
cd helios

# Go — install dependencies and run tests
go test ./...

# Python — install the package in editable mode and run tests
pip install -e .
pip install pytest
pytest implementations/python/tests -v
```

## Submitting a PR

- [ ] All tests pass (`go test ./...` and `pytest implementations/python/tests -v`)
- [ ] `scripts/cross_check.sh` passes (requires the Docker image; see `docker/Dockerfile`)
- [ ] One logical change per PR — keep diffs focused and reviewable

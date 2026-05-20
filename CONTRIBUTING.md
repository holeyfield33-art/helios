# Contributing

## Scope

Helios is a frozen trust primitive. The canonical serialization and hashing behavior defined by spec version 1, including the frozen test vectors, is not expected to change. Contributions that modify the specification or test vectors are not accepted unless they correct a provable correctness bug.

## What contributions are welcome

- New language implementations that pass all 17 frozen vectors
- Documentation improvements
- CI / tooling improvements
- Bug reports with reproducible test cases

## Development setup

```bash
git clone https://github.com/holeyfield33-art/helios.git
cd helios
python -m venv .venv
. .venv/bin/activate
pip install -e .[test]
bash scripts/install_hooks.sh
pytest implementations/python/tests -v
go test ./...
```

## Commit and merge safety policy

- Local commits and pushes are blocked if quality checks fail.
- Install once per clone: `bash scripts/install_hooks.sh`
- Manual gate run: `bash scripts/quality_gate.sh`
- CI also enforces this gate. Python warnings are treated as errors.

## Submitting a PR

- [ ] Tests pass locally
- [ ] `bash scripts/quality_gate.sh` passes
- [ ] `bash scripts/cross_check.sh` passes
- [ ] PR contains one logical change

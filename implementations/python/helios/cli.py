"""CLI entry points for helios-core.

Registered via [project.scripts] in pyproject.toml so the installed commands
are namespaced under the public ``helios`` package rather than the internal
``conformance`` package.
"""


def verify_main() -> None:
    """Entry point for the ``helios-verify`` console script."""
    from conformance.verifier import main

    main()

#!/usr/bin/env python3
"""Standalone entry point for cross_check.sh — runs vector verification."""

import sys
import os

# Add implementations/python to path so conformance package is importable
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from conformance.verifier import main

if __name__ == "__main__":
    main()

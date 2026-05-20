"""
Helios: Deterministic canonical serialization and SHA-256 content hashing for AI memory objects.

This module provides a clean public API for hashing memory objects using the frozen Helios spec.

Example:
    from helios import hash_memory_object
    
    obj = {
        "category": "project",
        "created_at": "2025-01-15T10:30:00.000Z",
        "key": "test/example",
        "relationships": [],
        "source": "user",
        "value": "Example memory"
    }
    content_hash = hash_memory_object(obj)
    print(content_hash)  # SHA-256 hex digest
"""

from conformance.hasher import hash_memory_object, MemoryObject
from conformance.verifier import verify_vectors, VerificationResult

__version__ = "1.0.0"
__all__ = [
    "hash_memory_object",
    "MemoryObject",
    "verify_vectors",
    "VerificationResult",
]

__author__ = "Ashura Joseph Holeyfield"
__license__ = "MIT"

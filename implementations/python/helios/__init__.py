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

from conformance.hasher import content_hash as _content_hash
from conformance.hasher import MemoryObject
from conformance.verifier import verify_vectors


def hash_memory_object(obj):
    """Hash a memory object and return its SHA-256 content hash as hex digest.
    
    Args:
        obj: A MemoryObject or dict with required fields (category, created_at, key, relationships, source, value).
        
    Returns:
        str: SHA-256 content hash (64-character hex digest).
        
    Raises:
        ValueError: If the object violates spec constraints (e.g., floats, nulls).
    """
    return _content_hash(obj)


__version__ = "1.0.0"
__all__ = [
    "hash_memory_object",
    "MemoryObject",
    "verify_vectors",
]

__author__ = "Ashura Joseph Holeyfield"
__license__ = "MIT"

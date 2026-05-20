"""
helios-core: Deterministic canonical serialization and SHA-256 content hashing
for AI memory objects.

Public API::

    from helios import hash_memory_object, canonicalize, MemoryObject, Relationship

    # Dict-based (most common):
    h = hash_memory_object({
        "category": "project",
        "created_at": "2025-01-15T10:30:00.000Z",
        "key": "my/object",
        "relationships": [{"key": "other/obj", "type": "related_to"}],
        "source": "user",
        "value": "hello",
    })

    # Typed-object API:
    obj = MemoryObject(
        category="project",
        created_at="2025-01-15T10:30:00.000Z",
        key="my/object",
        relationships=[Relationship(key="other/obj", type="related_to")],
        source="user",
        value="hello",
    )
    h = hash_memory_object(obj)
"""

from conformance.canon import canonicalize_object as _canonicalize_object
from conformance.hasher import content_hash as _content_hash
from conformance.objects import HashInput, MemoryObject, Relationship
from conformance.verifier import verify_vectors


def hash_memory_object(obj) -> str:
    """Compute the deterministic SHA-256 content hash for a memory object.

    Accepts either a :class:`MemoryObject` instance or a plain ``dict`` with
    the six hash-input fields: ``category``, ``created_at``, ``key``,
    ``relationships``, ``source``, ``value``.

    Returns:
        str: 64-character lowercase hex SHA-256 digest.

    Raises:
        ValueError: For spec violations (CANON_ERR_FLOAT_PROHIBITED,
            CANON_ERR_NULL_PROHIBITED, invalid timestamp, etc.)
        TypeError: If *obj* is neither a ``MemoryObject`` nor a ``dict``.
    """
    if isinstance(obj, MemoryObject):
        return _content_hash(obj)
    if isinstance(obj, dict):
        rels = [
            Relationship(key=r["key"], type=r["type"])
            for r in obj.get("relationships", [])
        ]
        mem_obj = MemoryObject(
            category=obj.get("category", ""),
            created_at=obj.get("created_at", ""),
            key=obj.get("key", ""),
            relationships=rels,
            source=obj.get("source", ""),
            value=obj.get("value"),
        )
        return _content_hash(mem_obj)
    raise TypeError(
        f"hash_memory_object expects MemoryObject or dict, got {type(obj).__name__}"
    )


def canonicalize(obj: dict) -> bytes:
    """Produce deterministic canonical JSON bytes from a dict.

    Keys are sorted lexicographically at every level. Output is compact UTF-8
    with no whitespace. Non-ASCII characters are preserved as raw UTF-8 bytes,
    not escaped to ``\\uXXXX``.

    Returns:
        bytes: Canonical UTF-8 JSON bytes.

    Raises:
        ValueError: If *obj* contains null values (CANON_ERR_NULL_PROHIBITED).
        TypeError: If *obj* is not a dict.
    """
    return _canonicalize_object(obj)


def verify(vectors_path: str):
    """Verify test vectors from a JSON file against this implementation.

    Args:
        vectors_path: Path to a ``vectors.json`` file.

    Returns:
        Tuple of ``(results, failures)`` where *results* is a list of
        ``(vector_id, expected, got, passed)`` tuples and *failures* is the
        count of failed vectors.
    """
    return verify_vectors(vectors_path)


__version__ = "1.0.0"
__all__ = [
    "hash_memory_object",
    "canonicalize",
    "verify",
    "verify_vectors",
    "MemoryObject",
    "Relationship",
    "HashInput",
]

__author__ = "Ashura Joseph Holeyfield"
__license__ = "MIT"

"""Memory object types for Helios Core (Python conformance)."""

from dataclasses import dataclass
from typing import Any


@dataclass
class Relationship:
    key: str
    type: str


@dataclass
class MemoryObject:
    """Full memory object with all fields."""

    # Included in hash (6 fields):
    category: str
    created_at: str
    key: str
    relationships: list  # list of Relationship
    source: str
    value: Any  # nullable — None must NOT be omitted

    # Excluded from hash:
    updated_at: str = ""
    version: int = 0
    access_count: int = 0
    last_accessed: str = ""
    confidence: float = 0.0


@dataclass
class HashInput:
    """Contains ONLY the 6 fields included in the content hash.
    CRITICAL: None value must serialize as null, never be omitted.
    """

    category: str
    created_at: str
    key: str
    relationships: list  # list of Relationship
    source: str
    value: Any  # nullable


def new_hash_input(obj: MemoryObject) -> HashInput:
    """Extract only the 6 hash-relevant fields from a MemoryObject."""
    return HashInput(
        category=obj.category,
        created_at=obj.created_at,
        key=obj.key,
        relationships=obj.relationships,
        source=obj.source,
        value=obj.value,
    )

"""SHA-256 content hash for Helios Core (Python conformance)."""

import hashlib

from conformance.canon import (
    canonicalize_object,
    normalize_string,
    normalize_timestamp,
    relationship_to_map,
    sort_relationships,
)
from conformance.objects import MemoryObject, new_hash_input


def content_hash(obj: MemoryObject) -> str:
    """Compute the deterministic content hash for a MemoryObject.

    Steps:
      1. Extract HashInput (6 fields only)
      2. Normalize timestamp
      3. Sort relationships by key, then type
      4. NFC-normalize all string fields
      5. Build explicit field map
      6. Canonicalize → SHA-256 → hex
    """
    inp = new_hash_input(obj)

    # Step 2: Normalize timestamp
    inp.created_at = normalize_timestamp(inp.created_at)

    # Step 3: Sort relationships
    sorted_rels = sort_relationships(inp.relationships)

    # Step 4: NFC-normalize string fields
    inp.category = normalize_string(inp.category)
    inp.key = normalize_string(inp.key)
    inp.source = normalize_string(inp.source)
    if isinstance(inp.value, str):
        inp.value = normalize_string(inp.value)

    # NFC-normalize relationship strings
    rel_maps = []
    for r in sorted_rels:
        rel_maps.append({
            "key": normalize_string(r.key),
            "type": normalize_string(r.type),
        })

    # Step 5: Build explicit field map with exactly 7 keys (6 data + schema version)
    fields = {
        "_helios_schema_version": "1",
        "category": inp.category,
        "created_at": inp.created_at,
        "key": inp.key,
        "relationships": rel_maps,
        "source": inp.source,
        "value": inp.value,
    }

    # Step 6: Canonicalize → SHA-256 → hex
    canonical = canonicalize_object(fields)
    return hashlib.sha256(canonical).hexdigest()

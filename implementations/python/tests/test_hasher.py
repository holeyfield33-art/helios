"""Unit tests for conformance/hasher.py."""

import pytest

from conformance.hasher import content_hash
from conformance.objects import MemoryObject, Relationship


def _base_object():
    return MemoryObject(
        category="project",
        created_at="2025-01-15T10:30:00.000Z",
        key="test/basic_memory",
        relationships=[Relationship(key="project/helios", type="related_to")],
        source="user",
        value="This is a test memory for hash verification.",
    )


class TestContentHash:
    def test_excluded_fields_do_not_affect_hash(self):
        obj1 = _base_object()
        obj2 = _base_object()
        obj2.updated_at = "2099-12-31T23:59:59.999Z"
        obj2.version = 999
        obj2.access_count = 999999
        obj2.last_accessed = "2099-12-31T23:59:59.999Z"
        obj2.confidence = 0.001
        assert content_hash(obj1) == content_hash(obj2)

    def test_value_change_changes_hash(self):
        obj1 = _base_object()
        obj2 = _base_object()
        obj2.value = "A completely different value."
        assert content_hash(obj1) != content_hash(obj2)

    def test_nfd_nfc_hash_parity(self):
        obj1 = _base_object()
        obj1.value = "caf\u00e9"  # NFC

        obj2 = _base_object()
        obj2.value = "cafe\u0301"  # NFD

        assert content_hash(obj1) == content_hash(obj2)

    def test_nil_value_included(self):
        obj = _base_object()
        obj.value = None
        h = content_hash(obj)
        assert len(h) == 64

        # nil and string "null" must differ
        obj2 = _base_object()
        obj2.value = "null"
        assert content_hash(obj) != content_hash(obj2)

    def test_hash_stability(self):
        obj = _base_object()
        h1 = content_hash(obj)
        h2 = content_hash(obj)
        assert h1 == h2

    def test_empty_relationships(self):
        obj = _base_object()
        obj.relationships = []
        h = content_hash(obj)
        assert len(h) == 64

    def test_hash_is_64_hex_chars(self):
        obj = _base_object()
        h = content_hash(obj)
        assert len(h) == 64
        assert all(c in "0123456789abcdef" for c in h)

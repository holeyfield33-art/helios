"""Hardening / guard tests for Helios Core Python conformance."""

import dataclasses
import json
import os

from conformance.canon import canonicalize_object
from conformance.hasher import content_hash
from conformance.objects import HashInput, MemoryObject, Relationship


class TestGoldenCanonicalBytes:
    def test_golden_bytes(self):
        fields = {
            "category": "test",
            "created_at": "2025-01-01T00:00:00.000Z",
            "key": "golden/test",
            "relationships": [],
            "source": "unit_test",
            "value": "hello",
        }
        result = canonicalize_object(fields)
        expected = b'{"category":"test","created_at":"2025-01-01T00:00:00.000Z","key":"golden/test","relationships":[],"source":"unit_test","value":"hello"}'
        assert result == expected


class TestHashInputFieldGuard:
    def test_exactly_6_fields(self):
        spec_fields = sorted(["category", "created_at", "key", "relationships", "source", "value"])
        struct_fields = sorted([f.name for f in dataclasses.fields(HashInput)])
        assert spec_fields == struct_fields, f"HashInput fields diverge: {struct_fields}"


class TestVectorHashMatchesFrozenValue:
    def test_basic_vector_hash(self):
        obj = MemoryObject(
            category="project",
            created_at="2025-01-15T10:30:00.000Z",
            key="test/basic_memory",
            relationships=[Relationship(key="project/helios", type="related_to")],
            source="user",
            value="This is a test memory for hash verification.",
        )
        h = content_hash(obj)
        frozen = "cae6f0ca521caeb1f74470aeca5a75ff1fe098809a034e8a15e0eb4762b4f485"
        assert h == frozen, f"hash mismatch: got {h}, frozen {frozen}"

    def test_all_5_frozen_hashes(self):
        """Verify all 5 frozen hashes from the Go implementation."""
        expected_hashes = {
            "basic": "cae6f0ca521caeb1f74470aeca5a75ff1fe098809a034e8a15e0eb4762b4f485",
            "key_ordering": "437573e624f5c2a8ffbd08e7e1f8d5491b1bf0fad7287d989e1e50be19c00a0f",
            "unicode_normalization": "68e92122b2993e8c8a416dabe8c1af18dbb4621760d9c569abc0c0621e064732",
            "null_value": "7b23e07bdb8fb414ac689b62f78c790bbbce9abeb433f018e8c5883097a6e845",
            "relationship_sorting": "11d3af8b06e69c463484cbd36dc3ee880fb74c6459285515200a87a8ba1f9452",
        }

        # Find vectors.json relative to this test file
        test_dir = os.path.dirname(os.path.abspath(__file__))
        vectors_path = os.path.join(test_dir, "..", "..", "..", "test_vectors", "vectors.json")

        with open(vectors_path) as f:
            data = json.load(f)

        for vec in data["vectors"]:
            name = vec["name"]
            inp = vec["input"]

            relationships = []
            for r in inp.get("relationships", []):
                relationships.append(Relationship(key=r["key"], type=r["type"]))

            obj = MemoryObject(
                category=inp.get("category", ""),
                created_at=inp.get("created_at", ""),
                key=inp.get("key", ""),
                relationships=relationships,
                source=inp.get("source", ""),
                value=inp.get("value"),
            )

            computed = content_hash(obj)
            assert computed == expected_hashes[name], (
                f"Vector '{name}': expected {expected_hashes[name]}, got {computed}"
            )


class TestHashStabilityAcrossExcludedFields:
    def test_excluded_fields_ignored(self):
        obj1 = MemoryObject(
            category="project",
            created_at="2025-01-15T10:30:00.000Z",
            key="test/stability",
            relationships=[Relationship(key="project/helios", type="related_to")],
            source="user",
            value="Stability test",
            updated_at="2025-01-15T12:00:00.000Z",
            version=1,
            access_count=0,
            last_accessed="2025-01-15T10:30:00.000Z",
            confidence=1.0,
        )

        obj2 = MemoryObject(
            category="project",
            created_at="2025-01-15T10:30:00.000Z",
            key="test/stability",
            relationships=[Relationship(key="project/helios", type="related_to")],
            source="user",
            value="Stability test",
            updated_at="2099-12-31T23:59:59.999Z",
            version=999,
            access_count=1000000,
            last_accessed="2099-12-31T23:59:59.999Z",
            confidence=0.001,
        )

        assert content_hash(obj1) == content_hash(obj2)

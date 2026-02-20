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
            "_helios_schema_version": "1",
            "category": "test",
            "created_at": "2025-01-01T00:00:00.000Z",
            "key": "golden/test",
            "relationships": [],
            "source": "unit_test",
            "value": "hello",
        }
        result = canonicalize_object(fields)
        expected = b'{"_helios_schema_version":"1","category":"test","created_at":"2025-01-01T00:00:00.000Z","key":"golden/test","relationships":[],"source":"unit_test","value":"hello"}'
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
        frozen = "c3262407645dcdbd1cede212fa0448a3adb2f915f762540c32e0050bbf65e781"
        assert h == frozen, f"hash mismatch: got {h}, frozen {frozen}"

    def test_all_positive_frozen_hashes(self):
        """Verify all positive frozen hashes from frozen_vectors_v3."""
        expected_hashes = {
            "POS-001": "c3262407645dcdbd1cede212fa0448a3adb2f915f762540c32e0050bbf65e781",
            "POS-002": "694cafaa80dd0121a4c4415ac44793fee17104d02756b3c1456dd79fc467c1d0",
            "POS-003": "d7b4f1c46600c6b7f6733e866455cfa3c5646b6e63625a2107a6a57a36be486c",
            "POS-004": "5e43f9576ec448e9111856b8e0f95593e4aa427ba9ec71cb3a6b574a91719558",
            "POS-005": "84c6d544a9ee3b9c1bd48a17d8835f25a7df62cd520f78f12fa49810b9e35945",
        }

        # Find vectors.json relative to this test file
        test_dir = os.path.dirname(os.path.abspath(__file__))
        vectors_path = os.path.join(test_dir, "..", "..", "..", "test_vectors", "vectors.json")

        with open(vectors_path) as f:
            data = json.load(f)

        for vec in data["vectors"]:
            vector_id = vec["vector_id"]
            if vec["vector_type"] != "positive":
                continue

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
            assert computed == expected_hashes[vector_id], (
                f"Vector '{vector_id}': expected {expected_hashes[vector_id]}, got {computed}"
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

"""Unit tests for conformance/canon.py."""

import pytest
import unicodedata

from conformance.canon import (
    canonicalize_object,
    normalize_string,
    normalize_timestamp,
    relationship_to_map,
    sort_relationships,
)
from conformance.objects import Relationship


class TestNormalizeString:
    def test_nfc_normalization(self):
        nfd = "cafe\u0301"  # NFD: e + combining acute
        nfc = "caf\u00e9"   # NFC: precomposed é
        assert normalize_string(nfd) == normalize_string(nfc)
        assert normalize_string(nfc) == "café"

    def test_ascii_passthrough(self):
        assert normalize_string("hello") == "hello"

    def test_already_nfc(self):
        s = "café"
        assert normalize_string(s) == s


class TestNormalizeTimestamp:
    def test_valid_3_decimals(self):
        assert normalize_timestamp("2025-01-15T10:30:00.000Z") == "2025-01-15T10:30:00.000Z"

    def test_valid_with_ms(self):
        assert normalize_timestamp("2025-01-15T10:30:00.123Z") == "2025-01-15T10:30:00.123Z"

    def test_rejects_no_z(self):
        with pytest.raises(ValueError, match="CANON_ERR_TIMESTAMP_NON_UTC"):
            normalize_timestamp("2025-01-15T10:30:00.123+00:00")

    def test_rejects_no_decimals(self):
        with pytest.raises(ValueError, match="CANON_ERR_TIMESTAMP_INVALID_PRECISION"):
            normalize_timestamp("2025-01-15T10:30:00Z")

    def test_rejects_4_decimals(self):
        with pytest.raises(ValueError, match="CANON_ERR_TIMESTAMP_INVALID_PRECISION"):
            normalize_timestamp("2025-01-15T10:30:00.1234Z")

    def test_rejects_1_decimal(self):
        with pytest.raises(ValueError, match="CANON_ERR_TIMESTAMP_INVALID_PRECISION"):
            normalize_timestamp("2025-01-15T10:30:00.1Z")


class TestCanonicalizeObject:
    def test_key_ordering(self):
        result = canonicalize_object({"z": 1, "a": 2, "m": 3})
        assert result == b'{"a":2,"m":3,"z":1}'

    def test_nested_key_ordering(self):
        result = canonicalize_object({"z": {"y": 1, "b": 2}, "a": "x"})
        assert result == b'{"a":"x","z":{"b":2,"y":1}}'

    def test_null_rejection_top_level(self):
        with pytest.raises(ValueError, match="CANON_ERR_NULL_PROHIBITED"):
            canonicalize_object({"a": None, "b": "x"})

    def test_null_rejection_nested(self):
        with pytest.raises(ValueError, match="CANON_ERR_NULL_PROHIBITED"):
            canonicalize_object({"outer": {"inner": None}})

    def test_array_preservation(self):
        result = canonicalize_object({"arr": [3, 1, 2]})
        assert result == b'{"arr":[3,1,2]}'

    def test_utf8_preserved(self):
        result = canonicalize_object({"value": "café"})
        assert b"caf\\u" not in result  # no unicode escapes
        assert "café".encode("utf-8") in result

    def test_empty_array(self):
        result = canonicalize_object({"relationships": []})
        assert result == b'{"relationships":[]}'

    def test_bool_values(self):
        result = canonicalize_object({"a": True, "b": False})
        assert result == b'{"a":true,"b":false}'


class TestRelationships:
    def test_sort_by_key_then_type(self):
        rels = [
            Relationship(key="z", type="ref"),
            Relationship(key="a", type="child"),
            Relationship(key="a", type="parent"),
        ]
        sorted_rels = sort_relationships(rels)
        assert sorted_rels[0] == Relationship(key="a", type="child")
        assert sorted_rels[1] == Relationship(key="a", type="parent")
        assert sorted_rels[2] == Relationship(key="z", type="ref")

    def test_relationship_to_map(self):
        r = Relationship(key="test_key", type="depends_on")
        m = relationship_to_map(r)
        assert m == {"key": "test_key", "type": "depends_on"}
        assert len(m) == 2

    def test_relationship_object_key_order(self):
        r = Relationship(key="x", type="ref")
        result = canonicalize_object(relationship_to_map(r))
        assert result == b'{"key":"x","type":"ref"}'


class TestIngestValidation:
    """Tests for validate_ingest_value (RULE-002, RULE-009, RULE-010)."""

    def test_rejects_float(self):
        from conformance.canon import validate_ingest_value
        with pytest.raises(ValueError, match="CANON_ERR_FLOAT_PROHIBITED"):
            validate_ingest_value(3.14)

    def test_rejects_null(self):
        from conformance.canon import validate_ingest_value
        with pytest.raises(ValueError, match="CANON_ERR_NULL_PROHIBITED"):
            validate_ingest_value(None)

    def test_rejects_nested_null(self):
        from conformance.canon import validate_ingest_value
        with pytest.raises(ValueError, match="CANON_ERR_NULL_PROHIBITED"):
            validate_ingest_value({"outer": {"inner": None}})

    def test_rejects_integer_overflow(self):
        from conformance.canon import validate_ingest_value
        with pytest.raises(ValueError, match="CANON_ERR_INTEGER_OUT_OF_RANGE"):
            validate_ingest_value(9223372036854775808)  # int64 max + 1

    def test_rejects_negative_integer_overflow(self):
        from conformance.canon import validate_ingest_value
        with pytest.raises(ValueError, match="CANON_ERR_INTEGER_OUT_OF_RANGE"):
            validate_ingest_value(-9223372036854775809)  # int64 min - 1

    def test_accepts_valid_integer(self):
        from conformance.canon import validate_ingest_value
        validate_ingest_value(9223372036854775807)  # int64 max — should pass

    def test_accepts_valid_int64_min(self):
        from conformance.canon import validate_ingest_value
        validate_ingest_value(-9223372036854775808)  # int64 min — should pass

    def test_accepts_string(self):
        from conformance.canon import validate_ingest_value
        validate_ingest_value("hello")  # should pass

    def test_accepts_bool(self):
        from conformance.canon import validate_ingest_value
        validate_ingest_value(True)  # should pass

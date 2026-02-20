package hash

import (
	"testing"

	"github.com/holeyfield33-art/helios/internal/object"
)

func baseObject() object.MemoryObject {
	return object.MemoryObject{
		Category:  "project",
		CreatedAt: "2025-01-15T10:30:00.000Z",
		Key:       "test/basic_memory",
		Relationships: []object.Relationship{
			{Key: "project/helios", Type: "related_to"},
		},
		Source: "user",
		Value:  "This is a test memory for hash verification.",
		// Excluded fields:
		UpdatedAt:    "2025-01-15T12:00:00.000Z",
		Version:      3,
		AccessCount:  42,
		LastAccessed: "2025-01-16T08:00:00.000Z",
		Confidence:   0.95,
	}
}

func TestExcludedFieldsDoNotAffectHash(t *testing.T) {
	obj1 := baseObject()
	obj2 := baseObject()

	// Change all excluded fields
	obj2.UpdatedAt = "2099-12-31T23:59:59.999Z"
	obj2.Version = 999
	obj2.AccessCount = 999999
	obj2.LastAccessed = "2099-12-31T23:59:59.999Z"
	obj2.Confidence = 0.01

	h1, err := ContentHash(obj1)
	if err != nil {
		t.Fatalf("hash1 failed: %v", err)
	}
	h2, err := ContentHash(obj2)
	if err != nil {
		t.Fatalf("hash2 failed: %v", err)
	}

	if h1 != h2 {
		t.Errorf("excluded fields affected hash:\n  h1=%s\n  h2=%s", h1, h2)
	}
	if len(h1) != 64 {
		t.Errorf("hash should be 64 hex chars, got %d", len(h1))
	}
}

func TestValueChangeChangesHash(t *testing.T) {
	obj1 := baseObject()
	obj2 := baseObject()
	obj2.Value = "A completely different value."

	h1, err := ContentHash(obj1)
	if err != nil {
		t.Fatalf("hash1 failed: %v", err)
	}
	h2, err := ContentHash(obj2)
	if err != nil {
		t.Fatalf("hash2 failed: %v", err)
	}

	if h1 == h2 {
		t.Error("different values should produce different hashes")
	}
}

func TestNFDandNFCProduceSameHash(t *testing.T) {
	obj1 := baseObject()
	obj1.Value = "caf\u00e9" // NFC: é as single code point

	obj2 := baseObject()
	obj2.Value = "cafe\u0301" // NFD: e + combining acute accent

	h1, err := ContentHash(obj1)
	if err != nil {
		t.Fatalf("hash1 failed: %v", err)
	}
	h2, err := ContentHash(obj2)
	if err != nil {
		t.Fatalf("hash2 failed: %v", err)
	}

	if h1 != h2 {
		t.Errorf("NFC and NFD forms should produce same hash:\n  NFC=%s\n  NFD=%s", h1, h2)
	}
}

func TestNilValueIncluded(t *testing.T) {
	obj1 := baseObject()
	obj1.Value = nil

	h, err := ContentHash(obj1)
	if err != nil {
		t.Fatalf("hash failed: %v", err)
	}
	if len(h) != 64 {
		t.Errorf("hash should be 64 hex chars, got %d", len(h))
	}

	// Ensure nil value produces a different hash than string "null"
	obj2 := baseObject()
	obj2.Value = "null"
	h2, err := ContentHash(obj2)
	if err != nil {
		t.Fatalf("hash2 failed: %v", err)
	}
	if h == h2 {
		t.Error("nil value and string 'null' should produce different hashes")
	}
}

func TestHashStability(t *testing.T) {
	obj := baseObject()
	h1, err := ContentHash(obj)
	if err != nil {
		t.Fatalf("hash1 failed: %v", err)
	}
	// Compute again — must be identical
	h2, err := ContentHash(obj)
	if err != nil {
		t.Fatalf("hash2 failed: %v", err)
	}
	if h1 != h2 {
		t.Errorf("hash is not stable across calls:\n  h1=%s\n  h2=%s", h1, h2)
	}
}

func TestEmptyRelationships(t *testing.T) {
	obj := baseObject()
	obj.Relationships = []object.Relationship{}

	h, err := ContentHash(obj)
	if err != nil {
		t.Fatalf("hash failed: %v", err)
	}
	if len(h) != 64 {
		t.Errorf("hash should be 64 hex chars, got %d", len(h))
	}
}

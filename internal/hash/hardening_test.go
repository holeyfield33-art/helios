package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/holeyfield33-art/helios/internal/canon"
	"github.com/holeyfield33-art/helios/internal/object"
)

// TestGoldenCanonicalBytes verifies that a known input produces expected canonical bytes.
// This catches any accidental change to the serialization logic.
func TestGoldenCanonicalBytes(t *testing.T) {
	// Build a minimal known object
	fields := map[string]interface{}{
		"category":      "test",
		"created_at":    "2025-01-01T00:00:00.000Z",
		"key":           "golden/test",
		"relationships": []interface{}{},
		"source":        "unit_test",
		"value":         "hello",
	}

	canonical, err := canon.CanonicalizeObject(fields)
	if err != nil {
		t.Fatal(err)
	}

	// The expected canonical form (keys sorted, compact, UTF-8 preserved):
	expected := `{"category":"test","created_at":"2025-01-01T00:00:00.000Z","key":"golden/test","relationships":[],"source":"unit_test","value":"hello"}`

	if string(canonical) != expected {
		t.Errorf("golden bytes mismatch:\n  expected: %s\n  got:      %s", expected, string(canonical))
	}

	// Also verify the hash of these exact bytes
	sum := sha256.Sum256(canonical)
	hash := hex.EncodeToString(sum[:])
	if len(hash) != 64 {
		t.Errorf("hash should be 64 hex chars, got %d", len(hash))
	}

	// The hash of the golden bytes should be stable
	sum2 := sha256.Sum256([]byte(expected))
	hash2 := hex.EncodeToString(sum2[:])
	if hash != hash2 {
		t.Errorf("golden hash mismatch:\n  from canonical: %s\n  from expected:  %s", hash, hash2)
	}
}

// TestHashInputFieldGuard uses reflection to verify HashInput has exactly the 6 spec fields.
func TestHashInputFieldGuard(t *testing.T) {
	expectedFields := map[string]string{
		"Category":      "string",
		"CreatedAt":     "string",
		"Key":           "string",
		"Relationships": "[]object.Relationship",
		"Source":        "string",
		"Value":         "interface {}",
	}

	typ := reflect.TypeOf(object.HashInput{})
	if typ.NumField() != 6 {
		t.Errorf("HashInput should have exactly 6 fields, got %d", typ.NumField())
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		expectedType, ok := expectedFields[field.Name]
		if !ok {
			t.Errorf("unexpected field in HashInput: %s", field.Name)
			continue
		}
		if field.Type.String() != expectedType {
			t.Errorf("field %s: expected type %s, got %s", field.Name, expectedType, field.Type.String())
		}
		// Verify NO omitempty in JSON tag
		tag := field.Tag.Get("json")
		if tag == "" {
			t.Errorf("field %s has no json tag", field.Name)
		}
		if contains(tag, "omitempty") {
			t.Errorf("field %s has omitempty — this breaks null serialization", field.Name)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestHashStabilityAcrossExcludedFields is the named test from the DoD checklist.
func TestHashStabilityAcrossExcludedFields(t *testing.T) {
	obj1 := object.MemoryObject{
		Category:  "project",
		CreatedAt: "2025-01-15T10:30:00.000Z",
		Key:       "test/stability",
		Relationships: []object.Relationship{
			{Key: "project/helios", Type: "related_to"},
		},
		Source:       "user",
		Value:        "Stability test",
		UpdatedAt:    "2025-01-15T12:00:00.000Z",
		Version:      1,
		AccessCount:  0,
		LastAccessed: "2025-01-15T10:30:00.000Z",
		Confidence:   1.0,
	}

	obj2 := object.MemoryObject{
		Category:  "project",
		CreatedAt: "2025-01-15T10:30:00.000Z",
		Key:       "test/stability",
		Relationships: []object.Relationship{
			{Key: "project/helios", Type: "related_to"},
		},
		Source:       "user",
		Value:        "Stability test",
		UpdatedAt:    "2099-12-31T23:59:59.999Z",
		Version:      999,
		AccessCount:  1000000,
		LastAccessed: "2099-12-31T23:59:59.999Z",
		Confidence:   0.001,
	}

	h1, err := ContentHash(obj1)
	if err != nil {
		t.Fatalf("hash1: %v", err)
	}
	h2, err := ContentHash(obj2)
	if err != nil {
		t.Fatalf("hash2: %v", err)
	}

	if h1 != h2 {
		t.Errorf("excluded fields changed the hash:\n  h1=%s\n  h2=%s", h1, h2)
	}
}

// TestVectorHashMatchesFrozenValue verifies the basic vector hash against its frozen value.
func TestVectorHashMatchesFrozenValue(t *testing.T) {
	obj := object.MemoryObject{
		Category:  "project",
		CreatedAt: "2025-01-15T10:30:00.000Z",
		Key:       "test/basic_memory",
		Relationships: []object.Relationship{
			{Key: "project/helios", Type: "related_to"},
		},
		Source: "user",
		Value:  "This is a test memory for hash verification.",
	}

	h, err := ContentHash(obj)
	if err != nil {
		t.Fatalf("hash: %v", err)
	}

	frozen := "cae6f0ca521caeb1f74470aeca5a75ff1fe098809a034e8a15e0eb4762b4f485"
	if h != frozen {
		t.Errorf("hash does not match frozen value:\n  got:    %s\n  frozen: %s", h, frozen)
	}
}

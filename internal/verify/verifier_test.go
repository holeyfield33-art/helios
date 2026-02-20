package verify

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVerifierExitsOnHashMismatch(t *testing.T) {
	// Create a temp vector file with a deliberately wrong hash (64 zeros)
	vectorJSON := `{
  "spec_version": "helios-canonical-serialization-v1",
  "vectors_version": "3",
  "vectors": [
    {
      "vector_id": "TEST-MISMATCH",
      "description": "Vector with wrong expected hash",
      "vector_type": "positive",
      "expected_outcome": "accept",
      "input": {
        "_helios_schema_version": "1",
        "category": "test",
        "created_at": "2025-01-15T10:30:00.000Z",
        "key": "test/mismatch",
        "relationships": [],
        "source": "user",
        "value": "test value"
      },
      "hash": "0000000000000000000000000000000000000000000000000000000000000000"
    }
  ]
}`

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "bad_vectors.json")
	if err := os.WriteFile(path, []byte(vectorJSON), 0644); err != nil {
		t.Fatal(err)
	}

	results, err := VerifyVectors(path)
	if err == nil {
		t.Fatal("expected error for hash mismatch, got nil")
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].Pass {
		t.Error("expected result to be a failure")
	}

	if results[0].Got == "0000000000000000000000000000000000000000000000000000000000000000" {
		t.Error("computed hash should not be all zeros")
	}
}

func TestVerifierPassesOnCorrectHash(t *testing.T) {
	// First compute the actual hash, then create a vector with it
	vectorJSON := `{
  "spec_version": "helios-canonical-serialization-v1",
  "vectors_version": "3",
  "vectors": [
    {
      "vector_id": "TEST-SELF",
      "description": "Compute hash then verify",
      "vector_type": "positive",
      "expected_outcome": "accept",
      "input": {
        "_helios_schema_version": "1",
        "category": "test",
        "created_at": "2025-01-15T10:30:00.000Z",
        "key": "test/self_check",
        "relationships": [],
        "source": "user",
        "value": "hello world"
      },
      "hash": "PLACEHOLDER"
    }
  ]
}`

	// Write with placeholder, run to get hash, rewrite with real hash
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "vectors.json")
	if err := os.WriteFile(path, []byte(vectorJSON), 0644); err != nil {
		t.Fatal(err)
	}

	// Get the actual hash from first run
	results, _ := VerifyVectors(path)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	actualHash := results[0].Got

	// Rewrite with correct hash
	correctedJSON := `{
  "spec_version": "helios-canonical-serialization-v1",
  "vectors_version": "3",
  "vectors": [
    {
      "vector_id": "TEST-SELF",
      "description": "Compute hash then verify",
      "vector_type": "positive",
      "expected_outcome": "accept",
      "input": {
        "_helios_schema_version": "1",
        "category": "test",
        "created_at": "2025-01-15T10:30:00.000Z",
        "key": "test/self_check",
        "relationships": [],
        "source": "user",
        "value": "hello world"
      },
      "hash": "` + actualHash + `"
    }
  ]
}`
	if err := os.WriteFile(path, []byte(correctedJSON), 0644); err != nil {
		t.Fatal(err)
	}

	results2, err := VerifyVectors(path)
	if err != nil {
		t.Fatalf("expected pass, got error: %v", err)
	}
	if len(results2) != 1 || !results2[0].Pass {
		t.Error("expected verification to pass")
	}
}

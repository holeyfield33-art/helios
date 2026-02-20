// Package verify implements test vector verification for Helios Core.
package verify

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/holeyfield33-art/helios/internal/hash"
	"github.com/holeyfield33-art/helios/internal/object"
)

// TestVector represents a single test vector from vectors.json.
type TestVector struct {
	Name                string                 `json:"name"`
	Description         string                 `json:"description"`
	Input               map[string]interface{} `json:"input"`
	ExpectedContentHash string                 `json:"expected_content_hash"`
}

// VectorsFile is the top-level structure of vectors.json.
type VectorsFile struct {
	Vectors []TestVector `json:"vectors"`
}

// VerifyResult holds the result of verifying a single vector.
type VerifyResult struct {
	Name     string
	Expected string
	Got      string
	Pass     bool
}

// VerifyVectors loads a vectors JSON file, computes the hash for each vector,
// and compares to the expected hash. Returns an error if ANY vector mismatches.
func VerifyVectors(path string) ([]VerifyResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read vectors file: %w", err)
	}

	dec := json.NewDecoder(strings.NewReader(string(data)))
	dec.UseNumber()

	var vf VectorsFile
	if err := dec.Decode(&vf); err != nil {
		return nil, fmt.Errorf("failed to parse vectors file: %w", err)
	}

	results := make([]VerifyResult, len(vf.Vectors))
	var failures int

	for i, vec := range vf.Vectors {
		obj, err := inputToMemoryObject(vec.Input)
		if err != nil {
			return nil, fmt.Errorf("vector %q: %w", vec.Name, err)
		}

		got, err := hash.ContentHash(obj)
		if err != nil {
			return nil, fmt.Errorf("vector %q hash failed: %w", vec.Name, err)
		}

		pass := got == vec.ExpectedContentHash
		results[i] = VerifyResult{
			Name:     vec.Name,
			Expected: vec.ExpectedContentHash,
			Got:      got,
			Pass:     pass,
		}

		if !pass {
			failures++
		}
	}

	if failures > 0 {
		return results, fmt.Errorf("%d of %d vectors failed verification", failures, len(vf.Vectors))
	}

	return results, nil
}

// inputToMemoryObject converts a raw JSON map into a MemoryObject.
func inputToMemoryObject(input map[string]interface{}) (object.MemoryObject, error) {
	obj := object.MemoryObject{}

	if v, ok := input["category"].(string); ok {
		obj.Category = v
	}
	if v, ok := input["created_at"].(string); ok {
		obj.CreatedAt = v
	}
	if v, ok := input["key"].(string); ok {
		obj.Key = v
	}
	if v, ok := input["source"].(string); ok {
		obj.Source = v
	}
	obj.Value = input["value"]

	if rels, ok := input["relationships"].([]interface{}); ok {
		for _, r := range rels {
			if rm, ok := r.(map[string]interface{}); ok {
				rel := object.Relationship{}
				if k, ok := rm["key"].(string); ok {
					rel.Key = k
				}
				if t, ok := rm["type"].(string); ok {
					rel.Type = t
				}
				obj.Relationships = append(obj.Relationships, rel)
			}
		}
	}
	if _, exists := input["relationships"]; exists && obj.Relationships == nil {
		obj.Relationships = []object.Relationship{}
	}

	if v, ok := input["updated_at"].(string); ok {
		obj.UpdatedAt = v
	}
	if v, ok := input["version"]; ok {
		switch vv := v.(type) {
		case json.Number:
			n, _ := vv.Int64()
			obj.Version = int(n)
		case float64:
			obj.Version = int(vv)
		}
	}
	if v, ok := input["access_count"]; ok {
		switch vv := v.(type) {
		case json.Number:
			n, _ := vv.Int64()
			obj.AccessCount = int(n)
		case float64:
			obj.AccessCount = int(vv)
		}
	}
	if v, ok := input["last_accessed"].(string); ok {
		obj.LastAccessed = v
	}
	if v, ok := input["confidence"]; ok {
		switch vv := v.(type) {
		case json.Number:
			f, _ := vv.Float64()
			obj.Confidence = f
		case float64:
			obj.Confidence = vv
		}
	}

	return obj, nil
}

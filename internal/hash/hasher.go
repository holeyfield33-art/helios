// Package hash implements the Helios Core content hash.
package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/holeyfield33-art/helios/internal/canon"
	"github.com/holeyfield33-art/helios/internal/object"
)

// ContentHash computes the deterministic content hash for a MemoryObject.
// Steps:
//  1. Extract HashInput (6 fields only)
//  2. Normalize timestamp
//  3. Sort relationships by key, then type
//  4. NFC-normalize all string fields
//  5. Build explicit field map
//  6. Canonicalize → SHA-256 → hex
func ContentHash(obj object.MemoryObject) (string, error) {
	// Step 0: Null prohibition check (RULE-010)
	if obj.Value == nil {
		return "", fmt.Errorf("CANON_ERR_NULL_PROHIBITED: null values are not permitted")
	}

	// Step 1: Extract only the 6 hash-relevant fields
	inp := object.NewHashInput(obj)

	// Step 2: Normalize timestamp
	ts, err := canon.NormalizeTimestamp(inp.CreatedAt)
	if err != nil {
		return "", fmt.Errorf("timestamp normalization failed: %w", err)
	}
	inp.CreatedAt = ts

	// Step 3: Sort relationships by key, then type as tie-breaker
	sortedRels := make([]map[string]interface{}, len(inp.Relationships))
	relMaps := make([]map[string]interface{}, len(inp.Relationships))
	for i, r := range inp.Relationships {
		relMaps[i] = canon.RelationshipToMap(r.Key, r.Type)
	}
	sorted := canon.SortRelationships(relMaps)
	copy(sortedRels, sorted)

	// Step 4: NFC-normalize string fields
	inp.Category = canon.NormalizeString(inp.Category)
	inp.Key = canon.NormalizeString(inp.Key)
	inp.Source = canon.NormalizeString(inp.Source)

	// NFC-normalize string values in relationships
	for i := range sortedRels {
		if k, ok := sortedRels[i]["key"].(string); ok {
			sortedRels[i]["key"] = canon.NormalizeString(k)
		}
		if t, ok := sortedRels[i]["type"].(string); ok {
			sortedRels[i]["type"] = canon.NormalizeString(t)
		}
	}

	// NFC-normalize Value if it's a string
	var normalizedValue interface{} = inp.Value
	if s, ok := inp.Value.(string); ok {
		normalizedValue = canon.NormalizeString(s)
	}

	// Step 5: Build EXPLICIT field map with exactly 6 keys
	// Keys must match the canonical JSON field names
	relsInterface := make([]interface{}, len(sortedRels))
	for i, r := range sortedRels {
		relsInterface[i] = r
	}

	fields := map[string]interface{}{
		"_helios_schema_version": "1",
		"category":               inp.Category,
		"created_at":             inp.CreatedAt,
		"key":                    inp.Key,
		"relationships":          relsInterface,
		"source":                 inp.Source,
		"value":                  normalizedValue,
	}

	// Step 6: Canonicalize → SHA-256 → hex
	canonical, err := canon.CanonicalizeObject(fields)
	if err != nil {
		return "", fmt.Errorf("canonicalization failed: %w", err)
	}

	sum := sha256.Sum256(canonical)
	return hex.EncodeToString(sum[:]), nil
}

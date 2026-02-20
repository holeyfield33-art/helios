// Package canon provides canonical serialization primitives for Helios Core.
// All normalization happens on INPUT values BEFORE serialization.
package canon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

// normalizeString applies NFC Unicode normalization to a string.
// Must be called on EVERY string field value before serialization.
func NormalizeString(s string) string {
	return norm.NFC.String(s)
}

// normalizeTimestamp validates and normalizes an ISO 8601 UTC timestamp
// to exactly YYYY-MM-DDTHH:MM:SS.sssZ (3 decimal places).
// Rejects timestamps not ending in Z or not having exactly 3 fractional digits.
func NormalizeTimestamp(s string) (string, error) {
	if !strings.HasSuffix(s, "Z") {
		return "", fmt.Errorf("CANON_ERR_TIMESTAMP_NON_UTC: timestamp must end in Z, got: %s", s)
	}

	// Validate exactly 3 fractional digits
	dotIdx := strings.LastIndex(s, ".")
	if dotIdx == -1 {
		return "", fmt.Errorf("CANON_ERR_TIMESTAMP_INVALID_PRECISION: timestamp must have exactly 3 fractional digits, got none: %s", s)
	}
	// Extract fractional part (between '.' and 'Z')
	frac := s[dotIdx+1 : len(s)-1] // strip trailing Z
	if len(frac) != 3 {
		return "", fmt.Errorf("CANON_ERR_TIMESTAMP_INVALID_PRECISION: timestamp must have exactly 3 fractional digits, got %d: %s", len(frac), s)
	}

	// Parse with explicit format — NEVER use time.RFC3339Nano
	t, err := time.Parse("2006-01-02T15:04:05.000Z", s)
	if err != nil {
		return "", fmt.Errorf("invalid timestamp format: %w", err)
	}

	return t.Format("2006-01-02T15:04:05.000Z"), nil
}

// CanonicalizeObject produces a deterministic JSON byte representation of a map.
// Keys are sorted lexicographically at every level. null values are preserved.
// UTF-8 is preserved (no \uXXXX escaping for non-ASCII). Arrays maintain insertion order.
func CanonicalizeObject(obj map[string]interface{}) ([]byte, error) {
	return canonicalizeValue(obj)
}

func canonicalizeValue(v interface{}) ([]byte, error) {
	switch val := v.(type) {
	case nil:
		return nil, fmt.Errorf("CANON_ERR_NULL_PROHIBITED: null values are not permitted")
	case bool:
		if val {
			return []byte("true"), nil
		}
		return []byte("false"), nil
	case json.Number:
		return []byte(val.String()), nil
	case float64:
		// Use strconv for shortest round-trip representation
		return []byte(strconv.FormatFloat(val, 'f', -1, 64)), nil
	case int:
		return []byte(strconv.Itoa(val)), nil
	case int64:
		return []byte(strconv.FormatInt(val, 10)), nil
	case string:
		return canonicalizeString(val)
	case map[string]interface{}:
		return canonicalizeMap(val)
	case []interface{}:
		return canonicalizeArray(val)
	default:
		return nil, fmt.Errorf("unsupported type: %T", v)
	}
}

// canonicalizeString writes a JSON string with UTF-8 preserved.
// Only characters that MUST be escaped in JSON are escaped.
func canonicalizeString(s string) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('"')
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		switch {
		case r == '"':
			buf.WriteString(`\"`)
		case r == '\\':
			buf.WriteString(`\\`)
		case r == '\b':
			buf.WriteString(`\b`)
		case r == '\f':
			buf.WriteString(`\f`)
		case r == '\n':
			buf.WriteString(`\n`)
		case r == '\r':
			buf.WriteString(`\r`)
		case r == '\t':
			buf.WriteString(`\t`)
		case r < 0x20:
			// Control characters must be escaped
			buf.WriteString(fmt.Sprintf(`\u%04x`, r))
		default:
			// Write raw UTF-8 bytes — do NOT escape to \uXXXX
			buf.Write([]byte(s[i : i+size]))
		}
		i += size
	}
	buf.WriteByte('"')
	return buf.Bytes(), nil
}

// canonicalizeMap serializes a map with explicitly sorted keys.
func canonicalizeMap(m map[string]interface{}) ([]byte, error) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		keyBytes, err := canonicalizeString(k)
		if err != nil {
			return nil, err
		}
		buf.Write(keyBytes)
		buf.WriteByte(':')

		valBytes, err := canonicalizeValue(m[k])
		if err != nil {
			return nil, err
		}
		buf.Write(valBytes)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// canonicalizeArray serializes an array, preserving insertion order.
func canonicalizeArray(arr []interface{}) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, v := range arr {
		if i > 0 {
			buf.WriteByte(',')
		}
		valBytes, err := canonicalizeValue(v)
		if err != nil {
			return nil, err
		}
		buf.Write(valBytes)
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

// SortRelationships sorts relationships by Key first, then Type as tie-breaker.
func SortRelationships(rels []map[string]interface{}) []map[string]interface{} {
	sorted := make([]map[string]interface{}, len(rels))
	copy(sorted, rels)
	sort.SliceStable(sorted, func(i, j int) bool {
		ki, _ := sorted[i]["key"].(string)
		kj, _ := sorted[j]["key"].(string)
		if ki != kj {
			return ki < kj
		}
		ti, _ := sorted[i]["type"].(string)
		tj, _ := sorted[j]["type"].(string)
		return ti < tj
	})
	return sorted
}

// RelationshipToMap converts key/type strings to an explicit map.
// NEVER rely on struct field ordering.
func RelationshipToMap(key, typ string) map[string]interface{} {
	return map[string]interface{}{
		"key":  key,
		"type": typ,
	}
}

// ValidateSchemaVersion checks RULE-001: _helios_schema_version must be present and equal to "1".
func ValidateSchemaVersion(input map[string]interface{}) error {
	v, exists := input["_helios_schema_version"]
	if !exists {
		return fmt.Errorf("CANON_ERR_SCHEMA_VERSION_MISSING: _helios_schema_version field is required")
	}
	s, ok := v.(string)
	if !ok || s != "1" {
		return fmt.Errorf("CANON_ERR_SCHEMA_VERSION_INVALID: _helios_schema_version must be string \"1\", got %v", v)
	}
	return nil
}

// ValidateIngestValue recursively validates a parsed JSON value for spec compliance.
// Checks: RULE-002 (no floats), RULE-009 (integer range), RULE-010 (no nulls).
// Expects values from json.Decoder with UseNumber().
func ValidateIngestValue(v interface{}) error {
	return validateIngest(v, "")
}

func validateIngest(v interface{}, path string) error {
	switch val := v.(type) {
	case nil:
		return fmt.Errorf("CANON_ERR_NULL_PROHIBITED: null value at %s", path)
	case float64:
		return fmt.Errorf("CANON_ERR_FLOAT_PROHIBITED: float value at %s", path)
	case json.Number:
		s := val.String()
		// Check for float indicators: decimal point or scientific notation
		if strings.Contains(s, ".") || strings.Contains(s, "e") || strings.Contains(s, "E") {
			return fmt.Errorf("CANON_ERR_FLOAT_PROHIBITED: numeric value %q at %s contains decimal or exponent", s, path)
		}
		// Check integer range (signed 64-bit)
		_, err := val.Int64()
		if err != nil {
			return fmt.Errorf("CANON_ERR_INTEGER_OUT_OF_RANGE: value %q at %s exceeds int64 bounds", s, path)
		}
	case map[string]interface{}:
		for k, child := range val {
			childPath := path + "." + k
			if err := validateIngest(child, childPath); err != nil {
				return err
			}
		}
	case []interface{}:
		for i, child := range val {
			childPath := fmt.Sprintf("%s[%d]", path, i)
			if err := validateIngest(child, childPath); err != nil {
				return err
			}
		}
	case string, bool:
		// Valid types, no checks needed
	default:
		return fmt.Errorf("unsupported type %T at %s", v, path)
	}
	return nil
}

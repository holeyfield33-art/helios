package canon

import (
	"encoding/json"
	"strings"
	"testing"

	"golang.org/x/text/unicode/norm"
)

func TestKeyOrdering(t *testing.T) {
	obj := map[string]interface{}{
		"zebra":  1,
		"alpha":  2,
		"middle": 3,
	}
	result, err := CanonicalizeObject(obj)
	if err != nil {
		t.Fatal(err)
	}
	expected := `{"alpha":2,"middle":3,"zebra":1}`
	if string(result) != expected {
		t.Errorf("expected %s, got %s", expected, string(result))
	}
}

func TestNullInclusion(t *testing.T) {
	obj := map[string]interface{}{
		"key":   "test",
		"value": nil,
	}
	result, err := CanonicalizeObject(obj)
	if err != nil {
		t.Fatal(err)
	}
	s := string(result)
	if !strings.Contains(s, `"value":null`) {
		t.Errorf("null value must be included, got: %s", s)
	}
}

func TestArrayPreservation(t *testing.T) {
	obj := map[string]interface{}{
		"items": []interface{}{"c", "a", "b"},
	}
	result, err := CanonicalizeObject(obj)
	if err != nil {
		t.Fatal(err)
	}
	expected := `{"items":["c","a","b"]}`
	if string(result) != expected {
		t.Errorf("expected %s, got %s", expected, string(result))
	}
}

func TestUTF8Preserved(t *testing.T) {
	obj := map[string]interface{}{
		"name": "héllo wörld 日本語",
	}
	result, err := CanonicalizeObject(obj)
	if err != nil {
		t.Fatal(err)
	}
	s := string(result)
	// Must contain raw UTF-8, not \uXXXX escapes
	if strings.Contains(s, `\u`) {
		t.Errorf("UTF-8 should be preserved, not escaped: %s", s)
	}
	if !strings.Contains(s, "héllo") {
		t.Errorf("should contain raw UTF-8 héllo, got: %s", s)
	}
	if !strings.Contains(s, "日本語") {
		t.Errorf("should contain raw UTF-8 日本語, got: %s", s)
	}
}

func TestTimestampExactly3Decimals(t *testing.T) {
	result, err := NormalizeTimestamp("2025-01-15T10:30:00.123Z")
	if err != nil {
		t.Fatal(err)
	}
	expected := "2025-01-15T10:30:00.123Z"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestTimestampRejects4Decimals(t *testing.T) {
	_, err := NormalizeTimestamp("2025-01-15T10:30:00.1234Z")
	if err == nil {
		t.Error("expected error for 4 decimal places, got nil")
	}
}

func TestTimestampRejectsNoZ(t *testing.T) {
	_, err := NormalizeTimestamp("2025-01-15T10:30:00.123+00:00")
	if err == nil {
		t.Error("expected error for non-Z suffix, got nil")
	}
}

func TestTimestampRejectsNoDecimals(t *testing.T) {
	_, err := NormalizeTimestamp("2025-01-15T10:30:00Z")
	if err == nil {
		t.Error("expected error for no decimal places, got nil")
	}
}

func TestNFCNormalization(t *testing.T) {
	// NFD form of "é" is e + combining acute accent (U+0065 U+0301)
	nfd := norm.NFD.String("é")
	// NFC form is single code point (U+00E9)
	nfc := NormalizeString(nfd)

	if nfc != "é" {
		t.Errorf("expected NFC é, got %q", nfc)
	}
	// Verify input was NFD and output is NFC
	if nfd == nfc {
		t.Error("NFD and NFC forms should differ in byte representation for this test")
	}
}

func TestNestedObjectKeyOrdering(t *testing.T) {
	obj := map[string]interface{}{
		"outer_b": map[string]interface{}{
			"inner_z": 1,
			"inner_a": 2,
		},
		"outer_a": "first",
	}
	result, err := CanonicalizeObject(obj)
	if err != nil {
		t.Fatal(err)
	}
	expected := `{"outer_a":"first","outer_b":{"inner_a":2,"inner_z":1}}`
	if string(result) != expected {
		t.Errorf("expected %s, got %s", expected, string(result))
	}
}

func TestRelationshipSorting(t *testing.T) {
	rels := []map[string]interface{}{
		{"key": "b_key", "type": "related_to"},
		{"key": "a_key", "type": "depends_on"},
		{"key": "a_key", "type": "blocks"},
	}
	sorted := SortRelationships(rels)

	// a_key/blocks should come before a_key/depends_on (tie-break by type)
	k0, _ := sorted[0]["key"].(string)
	t0, _ := sorted[0]["type"].(string)
	k1, _ := sorted[1]["key"].(string)
	t1, _ := sorted[1]["type"].(string)
	k2, _ := sorted[2]["key"].(string)

	if k0 != "a_key" || t0 != "blocks" {
		t.Errorf("expected a_key/blocks first, got %s/%s", k0, t0)
	}
	if k1 != "a_key" || t1 != "depends_on" {
		t.Errorf("expected a_key/depends_on second, got %s/%s", k1, t1)
	}
	if k2 != "b_key" {
		t.Errorf("expected b_key third, got %s", k2)
	}
}

func TestRelationshipToMap(t *testing.T) {
	m := RelationshipToMap("test_key", "depends_on")
	if m["key"] != "test_key" {
		t.Errorf("expected key=test_key, got %v", m["key"])
	}
	if m["type"] != "depends_on" {
		t.Errorf("expected type=depends_on, got %v", m["type"])
	}
	if len(m) != 2 {
		t.Errorf("expected exactly 2 fields, got %d", len(m))
	}
}

func TestEmptyArrayHandling(t *testing.T) {
	obj := map[string]interface{}{
		"relationships": []interface{}{},
	}
	result, err := CanonicalizeObject(obj)
	if err != nil {
		t.Fatal(err)
	}
	expected := `{"relationships":[]}`
	if string(result) != expected {
		t.Errorf("expected %s, got %s", expected, string(result))
	}
}

func TestJSONNumberHandling(t *testing.T) {
	// Simulate what json.Decoder with UseNumber produces
	input := `{"count":42,"ratio":3.14}`
	dec := json.NewDecoder(strings.NewReader(input))
	dec.UseNumber()
	var obj map[string]interface{}
	if err := dec.Decode(&obj); err != nil {
		t.Fatal(err)
	}
	result, err := CanonicalizeObject(obj)
	if err != nil {
		t.Fatal(err)
	}
	expected := `{"count":42,"ratio":3.14}`
	if string(result) != expected {
		t.Errorf("expected %s, got %s", expected, string(result))
	}
}

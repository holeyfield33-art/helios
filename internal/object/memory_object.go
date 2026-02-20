// Package object defines the memory object types for Helios Core.
package object

// Relationship represents a typed link between memory objects.
type Relationship struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}

// MemoryObject is the full memory object with all fields.
// Some fields are excluded from the content hash.
type MemoryObject struct {
	// Included in hash (6 fields):
	Category      string         `json:"category"`
	CreatedAt     string         `json:"created_at"`
	Key           string         `json:"key"`
	Relationships []Relationship `json:"relationships"`
	Source        string         `json:"source"`
	Value         interface{}    `json:"value"`

	// Excluded from hash:
	UpdatedAt    string  `json:"updated_at"`
	Version      int     `json:"version"`
	AccessCount  int     `json:"access_count"`
	LastAccessed string  `json:"last_accessed"`
	Confidence   float64 `json:"confidence"`
}

// HashInput contains ONLY the 6 fields included in the content hash.
// CRITICAL: NO omitempty on any field — nil Value must serialize as null.
type HashInput struct {
	Category      string         `json:"category"`
	CreatedAt     string         `json:"created_at"`
	Key           string         `json:"key"`
	Relationships []Relationship `json:"relationships"`
	Source        string         `json:"source"`
	Value         interface{}    `json:"value"`
}

// NewHashInput extracts only the 6 hash-relevant fields from a MemoryObject.
func NewHashInput(obj MemoryObject) HashInput {
	return HashInput{
		Category:      obj.Category,
		CreatedAt:     obj.CreatedAt,
		Key:           obj.Key,
		Relationships: obj.Relationships,
		Source:        obj.Source,
		Value:         obj.Value,
	}
}

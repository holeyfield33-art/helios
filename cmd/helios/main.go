package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/holeyfield33-art/helios/internal/hash"
	"github.com/holeyfield33-art/helios/internal/object"
	"github.com/holeyfield33-art/helios/internal/verify"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "hash":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: helios hash <file.json>")
			os.Exit(1)
		}
		if err := runHash(os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "verify":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: helios verify <vectors.json>")
			os.Exit(1)
		}
		if err := runVerify(os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Helios Core — Canonical Hash Tool")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  helios hash <file.json>      Compute content hash for a memory object")
	fmt.Fprintln(os.Stderr, "  helios verify <vectors.json>  Verify test vectors")
}

func runHash(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	dec := json.NewDecoder(strings.NewReader(string(data)))
	dec.UseNumber()

	var input map[string]interface{}
	if err := dec.Decode(&input); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	obj := mapToMemoryObject(input)
	h, err := hash.ContentHash(obj)
	if err != nil {
		return fmt.Errorf("hash computation failed: %w", err)
	}

	fmt.Println(h)
	return nil
}

func runVerify(path string) error {
	results, err := verify.VerifyVectors(path)

	for _, r := range results {
		status := "PASS"
		if !r.Pass {
			status = "FAIL"
		}
		fmt.Printf("  %s: %s\n", r.Name, status)
		if !r.Pass {
			fmt.Printf("    expected: %s\n", r.Expected)
			fmt.Printf("    got:      %s\n", r.Got)
		}
	}

	if err != nil {
		return err
	}

	fmt.Printf("\nAll %d vectors: PASS\n", len(results))
	return nil
}

func mapToMemoryObject(input map[string]interface{}) object.MemoryObject {
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

	return obj
}

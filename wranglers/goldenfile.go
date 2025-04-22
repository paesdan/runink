package wranglers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// LoadGoldenFileWithPath loads a golden file and unmarshals its JSON contents into the provided destination.
func LoadGoldenFileWithPath(path string, target interface{}) error {
	bytes, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("load golden file failed: %w", err)
	}
	if err := json.Unmarshal(bytes, target); err != nil {
		return fmt.Errorf("decode golden json: %w", err)
	}
	return nil
}

// SaveGoldenFileWithPath serializes the provided data and writes it to the golden file path.
func SaveGoldenFileWithPath(path string, data interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create golden dir: %w", err)
	}
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal golden json: %w", err)
	}
	return os.WriteFile(filepath.Clean(path), bytes, 0644)
}

// GoldenFilePathFromScenario creates a reusable golden file path given a base dir, scenario, and contract hash.
func GoldenFilePathFromScenario(baseDir, feature, scenario, contractHash string) string {
	cleanFeature := sanitizeFileName(feature)
	cleanScenario := sanitizeFileName(scenario)
	return filepath.Join(baseDir, cleanFeature, fmt.Sprintf("%s_%s.golden.json", cleanScenario, contractHash[:8]))
}

// sanitizeFileName ensures the scenario and feature names are file-safe.
func sanitizeFileName(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]+`)
	clean := re.ReplaceAllString(strings.ToLower(name), "_")
	return strings.Trim(clean, "_")
}

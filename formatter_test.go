package autocomplete

import (
	"os"
	"testing"
)

func TestDefaultFormatter(t *testing.T) {
	var _ Formatter = (*DefaultFormat)(nil)
	fmtr := DefaultFormat{}

	// Failing JSON file, invalid format.
	byts, err := os.ReadFile("icons.json")

	keywords, err := fmtr.FormatRead(byts, "icons.json")
	if err == nil {
		t.Errorf("Expected non-nil, got %v", err)
	}

	// Passing JSON
	byts, cleanup := testJsonFile(t, "test.json")
	keywords, err = fmtr.FormatRead(byts, "test.json")
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if len(keywords) != 3 {
		t.Errorf("Expected 3, got %v", len(keywords))
	}
	cleanup()

	// Passing TXT
	byts, cleanup = testTxtFile(t, "test.txt")
	keywords, err = fmtr.FormatRead(byts, "test.txt")
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if len(keywords) != 5 {
		t.Errorf("Expected 5, got %v", len(keywords))
	}

}

func TestKeywordListFormatter(t *testing.T) {

}

func TestDetectFileType(t *testing.T) {

	_, cleanup := testJsonFile(t, "sample.json")

	cleanup()

}

func testJsonFile(t *testing.T, filename string) ([]byte, func()) {
	t.Helper()
	fData := []byte(`["keyword1", "keyword2", "keyword3"]`)
	file, err := os.Create(filename)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	_, err = file.Write(fData)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if err := file.Close(); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	file.Close()

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	return data, func() {
		os.Remove(file.Name())
	}

}

func testTxtFile(t *testing.T, filename string) ([]byte, func()) {
	t.Helper()
	fileData := []string{"keywords", "keyword1", "keyword2", "keyword3"}
	file, err := os.Create(filename)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	for _, line := range fileData {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	}

	if err := file.Close(); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	return data, func() {
		os.Remove(file.Name())
	}

}

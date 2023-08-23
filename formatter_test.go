package autocomplete

import (
	"encoding/csv"
	"fmt"
	"os"
	"testing"
)

func TestDefaultFormatter(t *testing.T) {
	var _ Formatter = (*DefaultFormat)(nil)
	fmtr := DefaultFormat{}

	t.Run("Default format should read from json", func(t *testing.T) {
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

	})

	t.Run("Default format should read from txt", func(t *testing.T) {
		// Passing TXT
		byts, cleanup := testTxtFile(t, "test.txt")
		keywords, err := fmtr.FormatRead(byts, "test.txt")
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
		// TODO: Sensitive of empty lines
		// Should probably see about parsing those out.
		if len(keywords) != 4 {
			t.Errorf("Expected 4, got %v", len(keywords))
		}

		cleanup()

	})

	t.Run("Default format should read from csv", func(t *testing.T) {
		// Passing csv
		byts, cleanup := testCsvFile(t, "test.csv")
		keywords, err := fmtr.FormatRead(byts, "test.csv")
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
		fmt.Printf("%v\n", keywords)
		if len(keywords) != 3 {
			t.Errorf("Expected 3, got %v", len(keywords))
		}

		cleanup()

	})

	t.Run("Default format should read from yaml", func(t *testing.T) {
		// Passing TXT
		byts, cleanup := testTxtFile(t, "test.yaml")
		keywords, err := fmtr.FormatRead(byts, "test.yaml")
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
		if len(keywords) != 4 {
			t.Errorf("Expected 4, got %v", len(keywords))
		}

		cleanup()

	})

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

	for i, line := range fileData {
		var lineContent string
		if i == len(fileData)-1 {
			lineContent = line
		} else {
			lineContent = line + "\n"
		}
		_, err := file.WriteString(lineContent)
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

func testCsvFile(t *testing.T, filename string) ([]byte, func()) {
	t.Helper()

	fileData := []string{"keywords", "keyword1", "keyword2", "keyword3"}
	file, err := os.Create(filename)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	// NOTE: The first record here is considered the filedname/key
	wtr := csv.NewWriter(file)
	for _, kwd := range fileData {
		if err := wtr.Write([]string{kwd}); err != nil {
			t.Errorf("there was a problem writing to the csv: %v", err)
		}
	}

	wtr.Flush()

	if err := file.Close(); err != nil {
		t.Errorf("failed to close csv file: %v", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Errorf("failed to read csv file: %v", err)

	}

	return data, func() {
		os.Remove(file.Name())
	}
}

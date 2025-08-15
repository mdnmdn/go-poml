package poml

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBasicRender(t *testing.T) {
	pomlString := `<poml><p>Hello</p></poml>`
	result, err := RenderFromString(pomlString, nil)
	if err != nil {
		t.Fatalf("RenderFromString() failed: %v", err)
	}

	expected := "Hello"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestDocument(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "test-document")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	// Create a dummy PDF file
	pdfPath := filepath.Join(tmpDir, "test.pdf")
	pdfContent := "This is a test PDF"
	if err := os.WriteFile(pdfPath, []byte(pdfContent), 0644); err != nil {
		t.Fatalf("Failed to write dummy PDF: %v", err)
	}

	// Create the POML string with a document tag
	pomlString := `<poml><document src="` + pdfPath + `" /></poml>`

	// Render the POML string
	result, err := RenderFromString(pomlString, nil)
	if err != nil {
		t.Fatalf("RenderFromString() failed: %v", err)
	}

	// The expected result should be the content of the PDF file
	if result != pdfContent {
		t.Errorf("Expected '%s', got '%s'", pdfContent, result)
	}
}

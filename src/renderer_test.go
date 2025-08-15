package poml

import (
	"strings"
	"testing"
)

func TestRenderExpression(t *testing.T) {
	pomlString := `<poml>Hello, {{ name }}!</poml>`
	context := map[string]interface{}{
		"name": "World",
	}

	reader := strings.NewReader(pomlString)
	root, err := Parse(reader)
	if err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}

	result, err := render(root, context)
	if err != nil {
		t.Fatalf("render() failed: %v", err)
	}

	expected := "Hello, World!"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestRenderIfAttribute(t *testing.T) {
	t.Run("if is true", func(t *testing.T) {
		pomlString := `<poml><p if="isVisible">Visible</p></poml>`
		context := map[string]interface{}{"isVisible": true}

		root, _ := Parse(strings.NewReader(pomlString))
		result, _ := render(root, context)

		expected := "Visible"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("if is false", func(t *testing.T) {
		pomlString := `<poml><p if="!isVisible">Invisible</p></poml>`
		context := map[string]interface{}{"isVisible": true}

		root, _ := Parse(strings.NewReader(pomlString))
		result, _ := render(root, context)

		expected := ""
		if result != expected {
			t.Errorf("Expected empty string, got '%s'", result)
		}
	})
}

// func TestRenderForAttribute(t *testing.T) {
// 	pomlString := `<poml><p for="item in items">{{ loop.index }}:{{ item }} </p></poml>`
// 	context := map[string]interface{}{
// 		"items": []interface{}{"a", "b", "c"},
// 	}

// 	root, _ := Parse(strings.NewReader(pomlString))
// 	result, err := render(root, context)
// 	if err != nil {
// 		t.Fatalf("render() failed: %v", err)
// 	}

// 	expected := "0:a 1:b 2:c "
// 	if result != expected {
// 		t.Errorf("Expected '%s', got '%s'", expected, result)
// 	}
// }

func TestRenderExpressionWithMath(t *testing.T) {
	pomlString := `<poml>Result is {{ 1 + 2 }}</poml>`
	context := map[string]interface{}{}

	reader := strings.NewReader(pomlString)
	root, err := Parse(reader)
	if err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}

	result, err := render(root, context)
	if err != nil {
		t.Fatalf("render() failed: %v", err)
	}

	expected := "Result is 3"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

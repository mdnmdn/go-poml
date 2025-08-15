package poml

import (
	"bytes"
	"encoding/json"
	"os"
	"regexp"
	"strings"

	"github.com/dop251/goja"
)

var blockElements = map[string]bool{
	"div": true, "h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true, "ul": true, "ol": true, "li": true, "hr": true, "pre": true, "blockquote": true, "task": true, "role": true, "output-format": true, "hint": true, "example": true, "input": true, "output": true, "cp": true, "stepwise-instructions": true, "systemmessage": true, "humanmessage": true, "aimessage": true,
}

// render takes a root element and a context map, and returns the rendered string.
func render(root *Element, context map[string]interface{}) (string, error) {
	// 1. Find, parse, and apply the stylesheet
	stylesheet, err := parseStylesheet(root)
	if err != nil {
		return "", err
	}
	if stylesheet != nil {
		stylesheet.Apply(root)
	}

	// 2. Render the (potentially modified) tree
	var builder strings.Builder
	err = renderNode(&builder, root, context)
	if err != nil {
		return "", err
	}
	// Trim trailing space and newlines
	return strings.TrimSpace(builder.String()), nil
}

var forRegex = regexp.MustCompile(`^\s*(\w+)\s+in\s+(.+)\s*$`)

// renderNode is a recursive function that walks the node tree and writes the output to a strings.Builder.
func renderNode(builder *strings.Builder, node Node, context map[string]interface{}) error {
	switch n := node.(type) {
	case *Element:
		// Handle 'for' attribute
		if forExpression, ok := n.Attr["for"]; ok {
			parts := forRegex.FindStringSubmatch(forExpression)
			if len(parts) != 3 {
				// Invalid for expression, treat as text
				return renderElement(builder, n, context)
			}

			loopVar := parts[1]
			listExpr := parts[2]

			vm := goja.New()
			for key, value := range context {
				vm.Set(key, value)
			}
			listVal, err := vm.RunString(listExpr)
			if err != nil {
				return err
			}

			if list, ok := listVal.Export().([]interface{}); ok {
				for i, item := range list {
					loopContext := copyContext(context)
					loopContext[loopVar] = item
					loopContext["loop"] = map[string]interface{}{
						"index":  i,
						"first":  i == 0,
						"last":   i == len(list)-1,
						"length": len(list),
					}
					if err := renderElement(builder, n, loopContext); err != nil {
						return err
					}
				}
				return nil // The for loop has handled rendering this element
			}
		}

		// If not a for loop, render as a single element
		return renderElement(builder, n, context)

	case *Text:
		// Render text, evaluating any expressions.
		renderedContent, err := renderText(n.Content, context)
		if err != nil {
			return err
		}
		builder.WriteString(renderedContent)
	}
	return nil
}

// renderElement handles the rendering of a single element, including 'if' and 'let' logic.
func renderElement(builder *strings.Builder, n *Element, context map[string]interface{}) error {
	// Handle the 'if' attribute for conditional rendering.
	if ifCondition, ok := n.Attr["if"]; ok {
		vm := goja.New()
		for key, value := range context {
			vm.Set(key, value)
		}
		result, err := vm.RunString(ifCondition)
		if err != nil {
			return err
		}
		if !result.ToBoolean() {
			return nil // Condition is false, so we render nothing.
		}
	}

	// Handle <document> tag
	if n.Tag == "document" {
		if src, ok := n.Attr["src"]; ok {
			data, err := os.ReadFile(src)
			if err != nil {
				return err
			}
			builder.Write(data)
		}
		return nil // Document tag doesn't have children to render
	}

	// If condition passed or was not present, render children.
	// Create a local context for children to handle <let> scoping.
	localContext := copyContext(context)
	for _, child := range n.Children {
		if let, ok := child.(*Element); ok && let.Tag == "let" {
			name, nameOk := let.Attr["name"]
			if !nameOk {
				continue // <let> must have a name.
			}

			if value, ok := let.Attr["value"]; ok {
				localContext[name] = value
			} else if src, ok := let.Attr["src"]; ok {
				data, err := os.ReadFile(src)
				if err != nil {
					return err
				}
				var jsonData interface{}
				if err := json.Unmarshal(data, &jsonData); err == nil {
					if name, nameOk := let.Attr["name"]; nameOk {
						localContext[name] = jsonData
					} else {
						// If no name is provided, merge the JSON object into the context
						if dataMap, ok := jsonData.(map[string]interface{}); ok {
							for k, v := range dataMap {
								localContext[k] = v
							}
						}
					}
				}
			} else if len(let.Children) > 0 {
				if textNode, ok := let.Children[0].(*Text); ok {
					var jsonData interface{}
					if err := json.Unmarshal([]byte(textNode.Content), &jsonData); err != nil {
						return err
					}
					localContext[name] = jsonData
				}
			}
			continue // <let> produces no output
		}

		if include, ok := child.(*Element); ok && include.Tag == "include" {
			src, srcOk := include.Attr["src"]
			if !srcOk {
				continue // <include> must have a src attribute.
			}

			// Read and parse the included file
			includeData, err := os.ReadFile(src)
			if err != nil {
				return err // Consider how to handle file not found more gracefully
			}
			includedRoot, err := Parse(bytes.NewReader(includeData))
			if err != nil {
				return err
			}

			// Render the included content within the current context
			if err := renderNode(builder, includedRoot, localContext); err != nil {
				return err
			}
			continue // <include> tag itself has been processed.
		}

		if err := renderNode(builder, child, localContext); err != nil {
			return err
		}
	}

	// Add a newline after block elements
	if _, ok := blockElements[strings.ToLower(n.Tag)]; ok {
		builder.WriteString("\n")
	}

	return nil
}

func copyContext(c map[string]interface{}) map[string]interface{} {
	newC := make(map[string]interface{}, len(c))
	for k, v := range c {
		newC[k] = v
	}
	return newC
}

var expressionRegex = regexp.MustCompile(`\{\{([^}]+)\}\}`)

// renderText evaluates expressions in a string.
func renderText(content string, context map[string]interface{}) (string, error) {
	vm := goja.New()
	for key, value := range context {
		vm.Set(key, value)
	}

	result := expressionRegex.ReplaceAllStringFunc(content, func(match string) string {
		expr := strings.TrimSpace(match[2 : len(match)-2]) // Extract expression from {{...}}

		val, err := vm.RunString(expr)
		if err != nil {
			// In case of an error, return the original match to avoid breaking the string
			return match
		}

		return val.String()
	})

	return result, nil
}

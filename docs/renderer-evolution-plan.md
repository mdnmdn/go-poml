# POML Go Renderer Evolution Plan

This document outlines the features and improvements required for the Go POML renderer to achieve full functionality and pass all ported tests.

## 1. Advanced Whitespace Handling

**Problem:** The current renderer does not correctly handle whitespace, leading to formatting issues in the output. It collapses all whitespace, which is incorrect for many use cases.

**Proposed Solution:**
- Implement a whitespace handling model similar to HTML/XML.
- By default, collapse consecutive whitespace characters into a single space.
- Preserve all whitespace within elements that have the `whiteSpace="preserve"` attribute (similar to CSS `white-space: pre`). The `<pre>` tag is a good candidate for this behavior by default.
- Correctly handle newlines and indentation to produce human-readable output.

## 2. Block vs. Inline Elements

**Problem:** The renderer treats all elements the same, without distinguishing between block and inline elements. This is a major cause of the formatting issues.

**Proposed Solution:**
- Introduce a distinction between block and inline elements in the renderer.
- Block elements (e.g., `<div>`, `<p>`, `<h1>`) should be rendered on new lines.
- Inline elements (e.g., `<span>`, `<a>`) should be rendered on the same line as their surrounding content.
- This will require a more sophisticated rendering loop that is aware of the element's display type.

## 3. `syntax="xml"` Support

**Problem:** The renderer currently ignores the `syntax="xml"` attribute on the `<poml>` tag. The output is always plain text.

**Proposed Solution:**
- When `syntax="xml"` is present, the renderer should output a well-formed XML document that mirrors the structure of the POML input.
- This means that the output should include the tags, attributes, and content of the POML file, correctly formatted as XML.

## 4. Enhanced `<let>` Tag Functionality

**Problem:** The `<let>` tag with a `src` attribute does not correctly handle JSON files when no `name` attribute is provided.

**Proposed Solution:**
- When `<let src="..." />` is used without a `name` attribute, and the `src` file is a JSON object, the renderer should merge the fields of the JSON object into the current rendering context.
- This will allow expressions in the POML file to directly access the fields of the JSON object, as is the case in the Python implementation.

## 5. List Rendering (`<ul>`, `<ol>`, `<li>`)

**Problem:** The renderer does not have specific logic for handling lists.

**Proposed Solution:**
- Implement rendering logic for `<ul>` (unordered lists) and `<ol>` (ordered lists).
- `<li>` (list item) elements should be rendered with appropriate indentation and markers (e.g., bullets for `<ul>`, numbers for `<ol>`).

## 6. Table Rendering (`<table>`, `<tr>`, `<th>`, `<td>`)

**Problem:** The renderer does not support tables. The POML examples use tables, but the tests were not ported.

**Proposed Solution:**
- Implement rendering logic for tables.
- This should include support for `<table>`, `<tr>` (table row), `<th>` (table header), and `<td>` (table data) elements.
- The renderer should be able to format the table content in a readable way, with columns and rows.

## 7. `for` loop rendering

**Problem:** The `for` loop rendering adds extra newlines.

**Proposed Solution:**
- The rendering of elements inside a `for` loop should not add extra newlines unless the elements themselves are block-level elements that require a newline. The loop itself should not introduce formatting.

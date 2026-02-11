package email

import (
	"testing"
)

func TestRender_PlainText(t *testing.T) {
	body := &Body{
		Text: "This is plain text.\nWith multiple lines.\n",
		HTML: "",
	}

	result, err := Render(body)
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestRender_HTMLOnly(t *testing.T) {
	body := &Body{
		Text: "",
		HTML: `<!DOCTYPE html>
<html>
<body>
<h1>Heading</h1>
<p>This is a <strong>bold</strong> paragraph.</p>
<ul>
<li>Item 1</li>
<li>Item 2</li>
</ul>
</body>
</html>`,
	}

	result, err := Render(body)
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Should contain text representation of HTML
	// At minimum should have "Heading" and "paragraph" text
}

func TestRender_HTMLAndText(t *testing.T) {
	body := &Body{
		Text: "Plain text version",
		HTML: `<html><body><h1>HTML Version</h1><p>Better formatting.</p></body></html>`,
	}

	result, err := Render(body)
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Should prefer HTML and convert to markdown/ANSI
}

func TestRender_EmptyBody(t *testing.T) {
	body := &Body{
		Text: "",
		HTML: "",
	}

	result, err := Render(body)
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}

	if result != "" {
		t.Errorf("Expected empty result for empty body, got '%s'", result)
	}
}

func TestRender_HTMLWithComplexFormatting(t *testing.T) {
	body := &Body{
		HTML: `
<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
<h1>Main Title</h1>
<h2>Subtitle</h2>
<p>This is a paragraph with <em>emphasis</em> and <strong>strong</strong> text.</p>
<blockquote>A quoted section</blockquote>
<pre><code>Code block example</code></pre>
<a href="https://example.com">Link text</a>
</body>
</html>`,
	}

	result, err := Render(body)
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Should successfully convert complex HTML to readable format
}

func TestRender_MalformedHTML(t *testing.T) {
	body := &Body{
		Text: "Fallback text",
		HTML: "<html><body><p>Unclosed tag<body></html>",
	}

	result, err := Render(body)
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}

	// Should gracefully handle malformed HTML and potentially fall back to text
	if result == "" {
		t.Error("Expected non-empty result even with malformed HTML")
	}
}

func TestRender_HTMLWithEntities(t *testing.T) {
	body := &Body{
		HTML: `<p>&lt;div&gt; &amp; &quot;quotes&quot; &copy; 2024</p>`,
	}

	result, err := Render(body)
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Should properly decode HTML entities
}

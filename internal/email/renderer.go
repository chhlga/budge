package email

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/JohannesKaufmann/html-to-markdown"
	"github.com/charmbracelet/glamour"
)

func sanitizeMarkdown(markdown string) string {
	markdown = strings.ReplaceAll(markdown, "\r\n", "\n")

	linkPattern := regexp.MustCompile(`\[(.*?)\]\s*\([^)]*\)`)
	markdown = linkPattern.ReplaceAllString(markdown, "$1")

	imagePattern := regexp.MustCompile(`!\[(.*?)\]\s*\([^)]*\)`)
	markdown = imagePattern.ReplaceAllStringFunc(markdown, func(match string) string {
		altText := imagePattern.FindStringSubmatch(match)[1]
		if altText == "" {
			return "[Image]"
		}
		return fmt.Sprintf("[Image: %s]", altText)
	})

	markdown = regexp.MustCompile(`\n{3,}`).ReplaceAllString(markdown, "\n\n")

	markdown = regexp.MustCompile(`\n\s*\n`).ReplaceAllString(markdown, "\n\n")

	return strings.TrimSpace(markdown)
}

// Render converts an email body to ANSI-formatted text
// It tries HTML→Markdown→ANSI first, falls back to plain text
func Render(body *Body) (string, error) {
	if body == nil {
		return "", nil
	}

	// If we have HTML, try to convert it
	if body.HTML != "" {
		rendered, err := renderHTML(body.HTML)
		if err == nil && rendered != "" {
			return rendered, nil
		}
		// If HTML rendering fails, fall back to plain text
	}

	// Use plain text as-is or fallback
	if body.Text != "" {
		return renderPlainText(body.Text)
	}

	return "", nil
}

func renderHTML(html string) (string, error) {
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(html)
	if err != nil {
		return "", fmt.Errorf("failed to convert HTML to markdown: %w", err)
	}

	markdown = sanitizeMarkdown(markdown)

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create glamour renderer: %w", err)
	}

	ansi, err := renderer.Render(markdown)
	if err != nil {
		return "", fmt.Errorf("failed to render markdown to ANSI: %w", err)
	}

	return strings.TrimSpace(ansi), nil
}

func renderPlainText(text string) (string, error) {
	// For plain text, we can still use glamour to add some formatting
	// Treat it as preformatted text
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		// If glamour fails, return plain text as-is
		return text, nil
	}

	// Render as plain text (no markdown processing)
	ansi, err := renderer.Render(text)
	if err != nil {
		// If rendering fails, return plain text as-is
		return text, nil
	}

	return strings.TrimSpace(ansi), nil
}

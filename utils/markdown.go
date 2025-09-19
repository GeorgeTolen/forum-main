package utils

import (
	"bytes"

	"github.com/yuin/goldmark"
)

func RenderMarkdown(input string) string {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(input), &buf); err != nil {
		return input
	}
	return buf.String()
}

package main

import (
	"context"
	"strings"

	"go.lsp.dev/protocol"
	"go.uber.org/zap"
	"zochi.space/cc-tweaked/utils"
)

func (h *handler) Hover(ctx context.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	h.logger.Info("Hover called", zap.Any("params", params))

	docURI := params.TextDocument.URI
	content, exists := h.documents[docURI]
	if !exists {
		return &protocol.Hover{}, nil
	}

	line := params.Position.Line
	character := params.Position.Character
	lines := utils.SplitLines(content)
	if int(line) >= len(lines) {
		return nil, nil
	}

	currentLine := lines[line]
	if int(character) > len(currentLine) {
		return nil, nil
	}

	word := getWordAtPosition(currentLine, int(character))
	if word == "" {
		return &protocol.Hover{}, nil
	}

	h.logger.Info("Hover word", zap.String("word", word))

	hoverInfo := getHoverInfo(word)
	if hoverInfo == nil {
		return &protocol.Hover{}, nil
	}

	return hoverInfo, nil
}

// Helper to get word at cursor position
func getWordAtPosition(line string, character int) string {
	// Find start of word
	start := character
	for start > 0 && isIdentifierChar(rune(line[start-1])) {
		start--
	}

	// Find end of word
	end := character
	for end < len(line) && isIdentifierChar(rune(line[end])) {
		end++
	}

	// Check if there's a dot before (for module.method)
	if start > 0 && line[start-1] == '.' {
		// Find the module name before the dot
		moduleStart := start - 1
		for moduleStart > 0 && isIdentifierChar(rune(line[moduleStart-1])) {
			moduleStart--
		}
		return line[moduleStart:end]
	}

	return line[start:end]
}

func getHoverInfo(word string) *protocol.Hover {
	// Search for the word in completion items
	for _, item := range completionItems {
		// Check both the full label and without module prefix
		if item.Label == word {
			return createHoverFromItem(item)
		}

		// Also check if word matches the method part (e.g., "print" matches "term.print")
		if strings.Contains(item.Label, ".") {
			parts := strings.Split(item.Label, ".")
			if len(parts) == 2 && parts[1] == word {
				// Check if the line before cursor has the module name
				return createHoverFromItem(item)
			}
		}
	}
	return nil
}

func createHoverFromItem(item protocol.CompletionItem) *protocol.Hover {
	var content strings.Builder

	// Start with the signature from Detail
	if item.Detail != "" {
		content.WriteString("```lua\n")
		content.WriteString(item.Detail)
		content.WriteString("\n```\n\n")
	}

	// Add the documentation
	if item.Documentation != nil {
		if markupContent, ok := item.Documentation.(*protocol.MarkupContent); ok {
			// Skip the code block in documentation since we already added it from Detail
			docValue := markupContent.Value
			// Remove the first code block if it exists in documentation
			if strings.HasPrefix(docValue, "```lua\n") {
				// Find the end of the first code block
				endOfCodeBlock := strings.Index(docValue[8:], "```")
				if endOfCodeBlock != -1 {
					// Skip past the code block and the following newlines
					docValue = strings.TrimLeft(docValue[8+endOfCodeBlock+3:], "\n")
				}
			}
			content.WriteString(docValue)
		} else if strDoc, ok := item.Documentation.(string); ok {
			content.WriteString(strDoc)
		}
	}

	if content.Len() == 0 {
		return nil
	}

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.Markdown,
			Value: content.String(),
		},
	}
}

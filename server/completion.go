package main

import (
	"context"
	"strings"

	"go.lsp.dev/protocol"
	"go.uber.org/zap"
	"zochi.space/cc-tweaked/utils"
)

var completionItems = []protocol.CompletionItem{
	{
		Label:            "print",
		Kind:             protocol.CompletionItemKindFunction,
		Detail:           "print(text: string): nil",
		Documentation:    "Prints a message to the console",
		InsertText:       "print(${1:text})$0",
		InsertTextFormat: protocol.InsertTextFormatSnippet,
		SortText:         "!print", // Changed to !
	},
	{
		Label:  "term.print",
		Kind:   protocol.CompletionItemKindFunction,
		Detail: "term.print(text: string): nil",
		Documentation: &protocol.MarkupContent{
			Kind:  protocol.Markdown,
			Value: "Prints a message to the terminal.",
		},
		InsertText:       "print(${1:text})$0",
		InsertTextFormat: protocol.InsertTextFormatSnippet,
		SortText:         "!print", // Changed to ! (this was still "0print")
	},
	{
		Label:    "term",
		Kind:     protocol.CompletionItemKindModule,
		Detail:   "the terminal module that is used for all computers in computercraft",
		SortText: "!term", // Changed to !
	},
}

func (h *handler) Completion(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	h.logger.Info("Completion called", zap.Any("params", params))

	// gets the string data from the current active document.
	docURI := params.TextDocument.URI
	content, exists := h.documents[docURI]
	if !exists {
		return &protocol.CompletionList{
			IsIncomplete: false,
			Items:        []protocol.CompletionItem{},
		}, nil
	}

	line := params.Position.Line
	character := params.Position.Character

	lines := utils.SplitLines(content)
	if int(line) >= len(lines) {
		return &protocol.CompletionList{
			IsIncomplete: false,
			Items:        []protocol.CompletionItem{},
		}, nil
	}

	currentLine := lines[line]
	if int(character) > len(currentLine) {
		character = uint32(len(currentLine))
	}
	textBeforeCursor := currentLine[:character]

	items := GetCompletionItems(textBeforeCursor)

	return &protocol.CompletionList{
		IsIncomplete: false,
		Items:        items,
	}, nil

}

func GetCompletionItems(text string) []protocol.CompletionItem {

	trimmed := strings.TrimSpace(text)

	if strings.Contains(trimmed, ".") {
		parts := strings.Split(trimmed, ".")
		if len(parts) >= 2 {
			// example getting "term" from "term."
			moduleName := parts[len(parts)-2]
			partialMethod := parts[len(parts)-1]

			// extracts just the module name and removes any non-identifier characters
			for i := len(moduleName) - 1; i >= 0; i-- {
				if !isIdentifierChar(rune(moduleName[i])) {
					moduleName = moduleName[i+1:]
					break
				}
			}

			var filtered []protocol.CompletionItem
			prefix := moduleName + "."

			for _, item := range completionItems {
				if strings.HasPrefix(item.Label, prefix) {
					methodName := strings.TrimPrefix(item.Label, prefix)

					if strings.HasPrefix(methodName, partialMethod) {
						itemCopy := item
						itemCopy.Label = methodName
						itemCopy.SortText = "!" + methodName // Changed to !

						if itemCopy.InsertText == "" {
							itemCopy.InsertText = methodName
						}

						filtered = append(filtered, itemCopy)
					}
				}
			}
			return filtered
		}

	}

	if trimmed == "" {
		return completionItems
	}

	// Filter by what's been typed
	var filtered []protocol.CompletionItem
	for _, item := range completionItems {
		// Only show top-level items (modules and global functions)
		// Don't show "term.print" when typing "term" - only show "term"
		if !strings.Contains(item.Label, ".") && strings.HasPrefix(item.Label, trimmed) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

func isIdentifierChar(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_'
}

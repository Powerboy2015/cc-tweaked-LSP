package main

import "go.lsp.dev/protocol"

func GetCompletionItems() []protocol.CompletionItem {
	return []protocol.CompletionItem{
		{
			Label:         "print",
			Kind:          protocol.CompletionItemKindFunction,
			Detail:        "Prints a message to the console",
			Documentation: "Use this function to print messages.",
		},
		{
			Label:  "term.print",
			Kind:   protocol.CompletionItemKindFunction,
			Detail: "Prints a message to the terminal",
			Documentation: &protocol.MarkupContent{
				Kind:  protocol.Markdown,
				Value: "Moves the turtle forward one block. Returns `true` if successful, `false` otherwise.",
			},
		},
	}
}

package main

import "go.lsp.dev/protocol"

// Separate file in order to keep track of the completion items
//  in order to keep an proper overview of them.

func GetCompletionItems() []protocol.CompletionItem {
	return []protocol.CompletionItem{
		{
			Label:  "print",
			Kind:   3, // Function
			Detail: "print(value: any): void",
			Documentation: protocol.MarkupContent{
				Kind:  protocol.MarkupKind("CompletionItem"),
				Value: "Prints the given value to the console.",
			},
		},
		{
			Label:  "turrtle.back",
			Kind:   protocol.CompletionItemKindMethod,
			Detail: "turtle.back(steps: number): boolean",
			Documentation: protocol.MarkupContent{
				Kind:  protocol.MarkupKind("CompletionItem"),
				Value: "Moves the turtle backward by the specified number of steps.",
			},
		},
	}
}

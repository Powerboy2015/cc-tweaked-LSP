package jsonparser

import (
	"embed"
	"encoding/json"
	"fmt"

	"go.lsp.dev/protocol"
)

//go:embed data/cc-tweaked.json
var apiData embed.FS

type APIDefinition struct {
	Modules []Module   `json:"modules"`
	Globals []Function `json:"globals"`
}

type Module struct {
	Name          string     `json:"name"`
	Kind          string     `json:"kind"`
	Description   string     `json:"description"`
	Documentation string     `json:"documentation"`
	Functions     []Function `json:"functions"`
}

type Function struct {
	Name        string      `json:"name"`
	Signature   string      `json:"signature"`
	Description string      `json:"description"`
	Parameters  []Parameter `json:"parameters"`
	Returns     []Return    `json:"returns"`
	Example     string      `json:"example"`
}

type Parameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Optional    bool   `json:"optional,omitempty"`
}

type Return struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Optional    bool   `json:"optional,omitempty"`
}

// LoadAPIs loads the API definitions from the embedded JSON file
func LoadAPIs() (*APIDefinition, error) {
	data, err := apiData.ReadFile("data/cc-tweaked.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read apis.json: %w", err)
	}

	var apis APIDefinition
	if err := json.Unmarshal(data, &apis); err != nil {
		return nil, fmt.Errorf("failed to parse apis.json: %w", err)
	}

	return &apis, nil
}

// BuildCompletionItems converts API definitions to completion items
func BuildCompletionItems(apis *APIDefinition) []protocol.CompletionItem {
	var items []protocol.CompletionItem

	// Add global functions
	for _, fn := range apis.Globals {
		items = append(items, functionToCompletionItem(fn, ""))
	}

	// Add modules and their functions
	for _, module := range apis.Modules {
		// Add the module itself
		items = append(items, protocol.CompletionItem{
			Label:  module.Name,
			Kind:   protocol.CompletionItemKindModule,
			Detail: "module",
			Documentation: &protocol.MarkupContent{
				Kind:  protocol.Markdown,
				Value: formatModuleDocumentation(module),
			},
			SortText: "!" + module.Name,
		})

		// Add module functions
		for _, fn := range module.Functions {
			items = append(items, functionToCompletionItem(fn, module.Name))
		}
	}

	return items
}

func functionToCompletionItem(fn Function, moduleName string) protocol.CompletionItem {
	label := fn.Name
	if moduleName != "" {
		label = moduleName + "." + fn.Name
	}

	insertText := buildInsertText(fn)

	return protocol.CompletionItem{
		Label:  label,
		Kind:   protocol.CompletionItemKindFunction,
		Detail: fn.Signature,
		Documentation: &protocol.MarkupContent{
			Kind:  protocol.Markdown,
			Value: formatFunctionDocumentation(fn),
		},
		InsertText:       insertText,
		InsertTextFormat: protocol.InsertTextFormatSnippet,
		SortText:         "!" + fn.Name,
	}
}

func buildInsertText(fn Function) string {
	if len(fn.Parameters) == 0 {
		return fn.Name + "()$0"
	}

	insertText := fn.Name + "("
	for i, param := range fn.Parameters {
		if i > 0 {
			insertText += ", "
		}
		insertText += fmt.Sprintf("${%d:%s}", i+1, param.Name)
	}
	insertText += ")$0"

	return insertText
}

func formatFunctionDocumentation(fn Function) string {
	doc := fmt.Sprintf("```lua\n%s\n```\n\n%s", fn.Signature, fn.Description)

	if len(fn.Parameters) > 0 {
		doc += "\n\n**Parameters:**"
		for _, param := range fn.Parameters {
			optional := ""
			if param.Optional {
				optional = " (optional)"
			}
			doc += fmt.Sprintf("\n- `%s`: %s%s", param.Name, param.Description, optional)
		}
	}

	if len(fn.Returns) > 0 {
		doc += "\n\n**Returns:**"
		for _, ret := range fn.Returns {
			optional := ""
			if ret.Optional {
				optional = " (optional)"
			}
			doc += fmt.Sprintf("\n- `%s` (%s): %s%s", ret.Name, ret.Type, ret.Description, optional)
		}
	}

	if fn.Example != "" {
		doc += fmt.Sprintf("\n\n**Example:**\n```lua\n%s\n```", fn.Example)
	}

	return doc
}

func formatModuleDocumentation(module Module) string {
	return fmt.Sprintf("```lua\n%s\n```\n\n%s\n\n%s", module.Name, module.Description, module.Documentation)
}

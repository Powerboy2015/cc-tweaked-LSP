# CC-Tweaked Language Support for VS Code

A Visual Studio Code extension that provides comprehensive language support for CC-Tweaked Lua programming in Minecraft ComputerCraft.

## Features

-   ✅ **Autocomplete** for all CC-Tweaked APIs (turtle, peripheral, term, monitor, printer, etc.)
-   ✅ **Hover documentation** with detailed API descriptions and parameter information
-   ✅ **Type checking** for peripherals returned by `peripheral.wrap()`
-   ✅ **Intelligent parameter hints** while typing function calls
-   ✅ **Full API coverage** for all CC-Tweaked modules and peripherals

## How It Works

This extension integrates with the Lua Language Server (Sumneko) by:

1. **Generating EmmyLua type definitions** from CC-Tweaked API documentation
2. **Providing class definitions** for all peripheral types (monitor, printer, modem, drive, etc.)
3. **Adding type-safe overloads** for `peripheral.wrap()` to automatically infer peripheral types
4. **Injecting definitions** into the Lua LSP's workspace library

No custom LSP server needed - it extends the existing Lua LSP!

## Usage

### Basic Example

```lua
-- Terminal API (global)
term.clear()
term.setCursorPos(1, 1)
term.write("Hello, World!")

-- Turtle API (global)
turtle.forward()
turtle.turnRight()
turtle.dig()
```

### Peripheral Wrapping with Type Inference

```lua
-- Automatically typed as monitor
local mon = peripheral.wrap("monitor_0")
mon:setTextScale(2)
mon:clear()
mon:write("CC-Tweaked rocks!")

-- Automatically typed as printer
local printer = peripheral.wrap("printer_0")
printer:newPage()
printer:write("Printing...")
printer:endPage()
```

### Type Annotations

```lua
-- Explicit type annotation (optional, auto-inferred from peripheral name)
local mon = peripheral.wrap("monitor_0") ---@type monitor
mon:setCursorPos(1, 1)
```

## Project Structure

```
vscode-cc-tweaked/
├── src/
│   ├── extension.ts              # Main extension entry point
│   └── scripts/
│       └── generate-definitions.ts  # Generates Lua definitions from JSON
├── data/
│   └── cc-tweaked.json           # CC-Tweaked API documentation (source of truth)
├── definitions/
│   └── cc-tweaked.lua            # Generated EmmyLua type definitions
├── out/                          # Compiled TypeScript output
├── package.json                  # Extension manifest
├── tsconfig.json                 # TypeScript configuration
└── README.md
```

## Development

### Prerequisites

-   Node.js (version 16 or later)
-   npm (Node package manager)
-   Visual Studio Code

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/yourusername/vscode-cc-tweaked
    cd vscode-cc-tweaked
    ```

2. Install dependencies:

    ```bash
    npm install
    ```

3. Compile the extension:
    ```bash
    npm run compile
    ```

### Running the Extension

1. Open the project in Visual Studio Code
2. Press `F5` to launch the Extension Development Host
3. Open a `.lua` file and test CC-Tweaked autocomplete

### Updating API Definitions

The extension uses [`data/cc-tweaked.json`](data/cc-tweaked.json) as the source of truth for all CC-Tweaked APIs.

To add or update API definitions:

1. Edit [`data/cc-tweaked.json`](data/cc-tweaked.json)
2. Run the generator:
    ```bash
    npm run generate-defs
    ```
3. Recompile:
    ```bash
    npm run compile
    ```

### Scripts

-   `npm run compile` - Compile TypeScript and generate definitions
-   `npm run watch` - Watch for changes and recompile
-   `npm run generate-defs` - Generate Lua definitions from JSON
-   `npm run package` - Package extension as `.vsix` file
-   `npm run publish` - Publish to VS Code Marketplace

## API Coverage

### Global APIs

-   `term` - Terminal manipulation
-   `turtle` - Turtle movement and actions
-   `peripheral` - Peripheral detection and wrapping
-   `fs` - File system operations
-   `os` - Operating system functions
-   `textutils` - Text formatting and serialization
-   `redstone` / `rs` - Redstone control
-   `colors` / `colours` - Color constants
-   And many more...

### Peripheral Types

-   `monitor` - Advanced displays
-   `printer` - Text/graphic printing
-   `modem` - Network communication
-   `drive` - Disk drive access
-   `speaker` - Sound playback
-   `command` - Command block integration

## Requirements

-   [Lua Language Server](https://marketplace.visualstudio.com/items?itemName=sumneko.lua) extension (recommended for best experience)

## Contributing

Contributions are welcome! Here's how you can help:

1. **Add missing APIs** - Update [`data/cc-tweaked.json`](data/cc-tweaked.json) with new/missing functions
2. **Improve documentation** - Enhance function descriptions and examples
3. **Fix bugs** - Report issues or submit fixes
4. **Add features** - Suggest new functionality

### Contributing Steps

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Update [`data/cc-tweaked.json`](data/cc-tweaked.json) if adding/modifying APIs
5. Run `npm run compile` to regenerate definitions
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## TODO

-   [ ] Add support for ComputerCraft events
-   [ ] Include code snippets for common patterns
-   [ ] Add examples for all peripheral types
-   [ ] Support for custom peripheral APIs
-   [ ] Integration with CC:Tweaked documentation links

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

-   [CC-Tweaked](https://tweaked.cc/) - The amazing ComputerCraft mod
-   [Lua Language Server](https://github.com/LuaLS/lua-language-server) - Powering the LSP integration
-   ComputerCraft community for inspiration and support

## Links

-   [CC-Tweaked Documentation](https://tweaked.cc/)
-   [VS Code Extension API](https://code.visualstudio.com/api)
-   [Report Issues](https://github.com/yourusername/vscode-cc-tweaked/issues)

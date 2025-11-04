# vscode-cc-tweaked

## Overview

This project implements a Language Server Protocol (LSP) server and a Visual Studio Code extension for the ComputerCraft: Tweaked mod. The server is written in Go, while the client extension is developed using TypeScript.

## Project Structure

```
vscode-cc-tweaked
├── client
│   ├── src
│   │   └── extension.ts
│   ├── package.json
│   └── tsconfig.json
├── server
│   ├── main.go
│   ├── completion.go
│   └── go.mod
├── package.json
├── tsconfig.json
└── README.md
```

## Getting Started

### Prerequisites

-   Go (version 1.16 or later)
-   Node.js (version 14 or later)
-   npm (Node package manager)

### Installation

1. Clone the repository:

    ```
    git clone <repository-url>
    cd vscode-cc-tweaked
    ```

2. Install server dependencies:

    ```
    cd server
    go mod tidy
    ```

3. Install client dependencies:
    ```
    cd client
    npm install
    ```

### Running the Language Server

To run the language server, execute the following command in the `server` directory:

```
go run main.go
```

### Running the Extension

To run the extension in Visual Studio Code:

1. Open the `client` directory in Visual Studio Code.
2. Press `F5` to start debugging the extension. This will open a new instance of Visual Studio Code with the extension loaded.

## Usage

Once the extension is running, it will automatically connect to the language server. You can start using the features provided by the LSP, such as code completion and error checking.

## Contributing

Contributions are welcome! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and commit them.
4. Push your branch and create a pull request.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.

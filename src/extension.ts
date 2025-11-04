import * as vscode from 'vscode';
import { LanguageClient, LanguageClientOptions, ServerOptions, Executable } from 'vscode-languageclient/node';
import * as path from 'path';

let client: LanguageClient;

export function activate(context: vscode.ExtensionContext) {
    // Get configuration
    const config = vscode.workspace.getConfiguration('ccTweaked');
    const customServerPath = config.get<string>('serverPath');
    
    // Determine server executable path
    let serverExecutable: string;
    if (customServerPath && customServerPath.trim() !== '') {
        serverExecutable = customServerPath;
    } else {
        // Use bundled server (adjust path based on where you bundle the Go executable)
        serverExecutable = path.join(context.extensionPath, 'server', 'cc-tweaked-lsp.exe');
    }

    const run: Executable = {
        command: serverExecutable,
        args: []
    };

    const serverOptions: ServerOptions = {
        run,
        debug: run
    };

    const clientOptions: LanguageClientOptions = {
        documentSelector: [{ scheme: 'file', language: 'lua' }],
        synchronize: {
            fileEvents: vscode.workspace.createFileSystemWatcher('**/*.lua')
        }
    };

    client = new LanguageClient(
        'ccTweakedLanguageServer',
        'CC Tweaked Language Server',
        serverOptions,
        clientOptions
    );

    client.start();
}

export function deactivate(): Thenable<void> | undefined {
    if (!client) {
        return undefined;
    }
    return client.stop();
}
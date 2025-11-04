import * as vscode from 'vscode';
import * as path from 'path';

export function activate(context: vscode.ExtensionContext) {
    const defsPath = path.join(context.extensionPath, 'definitions');
    
    // Add CC-Tweaked definitions to Lua language server
    const config = vscode.workspace.getConfiguration('Lua');
    const library = config.get<string[]>('workspace.library') || [];
    
    if (!library.includes(defsPath)) {
        config.update('workspace.library', [...library, "out", defsPath], 
            vscode.ConfigurationTarget.Global);
        
        vscode.window.showInformationMessage('CC-Tweaked API definitions loaded!');
    }
}

export function deactivate() {}
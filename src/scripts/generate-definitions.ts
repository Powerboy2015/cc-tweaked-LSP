const fs = require('fs');
const path = require('path');

const jsonPath = path.join(__dirname, '..', '..', 'data',  'cc-tweaked.json');
const outputPath = path.join(__dirname, '..', 'definitions', 'cc-tweaked.lua');
// const snippetsOutputPath = path.join(__dirname, '..', 'snippets', 'cc-tweaked.json');

const data: ccDefs = JSON.parse(fs.readFileSync(jsonPath, 'utf8'));

let output = '---@meta\n\n';
let snippets: snippetList  = {};

data.globals.forEach(func => {
    output += generateFunctionDefinition(func);
    // snippets[func.name] = generateSnippet(func);
});


// These are the peripheral modules that are not globally available and are mostly returned by a return statement.
data.modules.forEach(mod => {
    output += `---@class ${mod.name}\n`;
    output += `---${mod.description.replace(/\n/g," ")}\n`;
    output += `local ${mod.name} = {}\n\n`;

    mod.functions.forEach(func => {
        output += generateFunctionDefinition(func, mod.name);
        // snippets[`${mod.name}.${func.name}`] = generateSnippet(func, mod.name);
    });

});

// These are the peripheral modules that are globally available, e.g. "term" or "os".
data.global_modules.forEach(mod => {
    output += `---@class ${mod.name}\n`;
    output += `---${mod.description.replace(/\n/g," ")}\n\n`;
    output += `${mod.name} = {}\n\n`;

    mod.functions.forEach(func => {
        output += generateFunctionDefinition(func, mod.name);
        // snippets[`${mod.name}.${func.name}`] = generateSnippet(func, mod.name);
    });
});




// Write definitions
fs.mkdirSync(path.dirname(outputPath), { recursive: true });
fs.writeFileSync(outputPath, output);
console.log('✅ Generated definitions/cc-tweaked.lua');

// // Write snippets
// fs.mkdirSync(path.dirname(snippetsOutputPath), { recursive: true });
// fs.writeFileSync(snippetsOutputPath, JSON.stringify(snippets, null, 2));
// console.log('✅ Generated snippets/cc-tweaked.json');

function generateFunctionDefinition(func: FunctionDef, moduleName?: string): string {
    let result = '';

    // Add description
    if (func.description) {
        result += `--${func.description.replace(/\n/g," ")}\n`;
    }

    // Adds parameters
    func.parameters.forEach(param => {
        result += `---@param ${param.name} ${param.type} ${param.description.replace(/\n/g," ")}\n`;
    });

    // adds return values
    func.returns.forEach(ret => {
        result += `---@return ${ret.type} ${ret.name} ${ret.description.replace(/\n/g," ")}\n`;
    });

    // adds function signature
    const funcName = moduleName ? `${moduleName}.${func.name}` : func.name;
    const params = func.parameters.map(p => p.name).join(', ');
    result += `function ${funcName}(${params}) end\n\n`;

    return result;
}

function generateSnippet(func: FunctionDef, moduleName?: string): Snippet {
    const funcName = moduleName ? `${moduleName}.${func.name}` : func.name;
    
    let body = funcName + '(';
    
    if (func.parameters.length > 0) {
        body += func.parameters.map((param, index) => {
            return `\${${index + 1}:${param.name}}`;
        }).join(', ');
    }
    
    body += ')$0';
    
    return {
        prefix: funcName,
        body: body,
        description: func.description
    };
}

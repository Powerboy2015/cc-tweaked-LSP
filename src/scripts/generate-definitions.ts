const fs = require('fs');
const path = require('path');

const jsonPath = path.join(__dirname, '..', '..', 'data',  'cc-tweaked.json');
const outputPath = path.join(__dirname, '..', 'definitions', 'cc-tweaked.lua');

const data: ccDefs = JSON.parse(fs.readFileSync(jsonPath, 'utf8'));

let output = '---@meta\n\n';

data.globals.forEach(func => {
    output += generateFunctionDefinition(func);
});

data.modules.forEach(mod => {
    output += `---@class ${mod.name}\n`;
    output += `---${mod.description}\n\n`;
    output += `${mod.name} = {}\n\n`;

    mod.functions.forEach(func => {
        output += generateFunctionDefinition(func, mod.name);
    });

});

fs.mkdirSync(path.dirname(outputPath), { recursive: true });
fs.writeFileSync(outputPath, output);
console.log('✅ Generated definitions/cc-tweaked.lua');

function generateFunctionDefinition(func: FunctionDef, moduleName?: string): string {
    let result = '';

    // Add description
    if (func.description) {
        result += `--${func.description}\n`;
    }

    // Adds parameters
    func.parameters.forEach(param => {
        result += `---@param ${param.name} ${param.type} ${param.description}\n`;
    });

    // adds return values
    func.returns.forEach(ret => {
        result += `---@return ${ret.type} ${ret.description}\n`;
    });

    // adds function signature
    const funcName = moduleName ? `${moduleName}.${func.name}` : func.name;
    const params = func.parameters.map(p => p.name).join(', ');
    result += `function ${funcName}(${params}) end\n\n`;

    return result;
}

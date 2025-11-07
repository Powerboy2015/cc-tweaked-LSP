const fs = require("fs");
const path = require("path");

const dataDir = path.join(__dirname, "..", "..", "data");
const outputDir = path.join(__dirname, "..", "definitions");
// const snippetsOutputPath = path.join(__dirname, '..', 'snippets', 'cc-tweaked.json');

const jsonFiles = fs
    .readdirSync(dataDir)
    .filter((file: string) => file.endsWith(".json"));

if (jsonFiles.length === 0) {
    console.error("❌ No JSON files found in data directory");
    process.exit(1);
}

jsonFiles.forEach((file: string) => {
    const jsonPath = path.join(dataDir, file);
    const outputFileName = file.replace(".json", ".lua");
    const outputPath = path.join(outputDir, outputFileName);

    try {
        const data: ccDefs = JSON.parse(fs.readFileSync(jsonPath, "utf8"));

        let output = "---@meta\n\n";
        let snippets: snippetList = {};

        data.globals.forEach((func) => {
            output += generateFunctionDefinition(func);
            // snippets[func.name] = generateSnippet(func);
        });

        // These are the peripheral modules that are not globally available and are mostly returned by a return statement.
        data.modules.forEach((mod) => {
            output += `---@class ${mod.name}\n`;
            output += `---${mod.description.replace(/\n/g, " ")}\n`;
            output += `${mod.name} = {}\n\n`;

            mod.functions.forEach((func) => {
                output += generateFunctionDefinition(func, mod.name);
                // snippets[`${mod.name}.${func.name}`] = generateSnippet(func, mod.name);
            });
        });

        // These are the return types.
        data.peripherals.forEach((mod) => {
            const extender = mod.extends ? ` : ${mod.extends}` : "";
            output += `---@class ${mod.name}${extender}\n`;
            output += `---${mod.description.replace(/\n/g, " ")}\n`;
            output += `local ${mod.name} = {}\n\n`;

            mod.functions.forEach((func) => {
                output += generateFunctionDefinition(func, mod.name);
            });
        });

        data.types.forEach((type) => {
            output += `---@class ${type.name}\n`;
            output += `---${type.description.replace(/\n/g, " ")}\n`;
            type.fields.forEach((field) => {
                output += `---@field ${field.name} ${field.type} ${field.description.replace(
                    /\n/g,
                    " "
                )}\n`;
            });
            output += `\n`;
        });

        // Write definitions
        fs.mkdirSync(path.dirname(outputPath), { recursive: true });
        fs.writeFileSync(outputPath, output);
        console.log(`✅ Generated ${outputFileName}`);
    } catch (err) {
        console.error(`❌ Error processing file ${file}:`, err);
    }
});

// // Write snippets
// fs.mkdirSync(path.dirname(snippetsOutputPath), { recursive: true });
// fs.writeFileSync(snippetsOutputPath, JSON.stringify(snippets, null, 2));
// console.log('✅ Generated snippets/cc-tweaked.json');

function generateFunctionDefinition(
    func: FunctionDef,
    moduleName?: string
): string {
    let result = "";

    // Add description
    if (func.description) {
        result += `--${func.description.replace(/\n/g, " ")}\n`;
    }

    if (func.name === "wrap" || func.name === "find") {
        result += `---@generic T: monitor|printer|modem|drive|speaker|command\n`;
    }

    // Adds parameters
    func.parameters.forEach((param) => {
        result += `---@param ${param.name} ${
            param.type
        } ${param.description.replace(/\n/g, " ")}\n`;
    });

    // adds return values
    func.returns.forEach((ret) => {
        result += `---@return ${ret.type} ${ret.name} ${ret.description.replace(
            /\n/g,
            " "
        )}\n`;
    });

    // adds function signature
    const funcName = moduleName ? `${moduleName}.${func.name}` : func.name;
    const params = func.parameters.map((p) => p.name).join(", ");
    result += `function ${funcName}(${params}) end\n\n`;

    return result;
}

function generateSnippet(func: FunctionDef, moduleName?: string): Snippet {
    const funcName = moduleName ? `${moduleName}.${func.name}` : func.name;

    let body = funcName + "(";

    if (func.parameters.length > 0) {
        body += func.parameters
            .map((param, index) => {
                return `\${${index + 1}:${param.name}}`;
            })
            .join(", ");
    }

    body += ")$0";

    return {
        prefix: funcName,
        body: body,
        description: func.description,
    };
}

type ccDefs ={
    modules: module[];
    globals: FunctionDef[];
}

type module = {
    name: string;
    kind: string;
    description: string;
    documentation: string;
    functions: FunctionDef[];
}

type FunctionDef = {
    name: string;
    signature: string;
    description: string;
    parameters: ParameterDef[];
    returns: ParameterDef[];
    example?: string;
}

type ParameterDef = {
    name: string;
    type: string;
    description: string;
}
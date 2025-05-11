import type { AdvanedFilterOperator } from '../../core-features/dynamic-filter/filter';

//STRING
export const eqStringOperator: AdvanedFilterOperator = {
    label: 'Equal (case sensitive)',
    name: '_', // no suffix
    field: 'string'
};

export const eqStringFoldOperator: AdvanedFilterOperator = {
    label: 'Equal',
    name: 'EqualFold',
    field: 'string'
};

export const neqOperator: AdvanedFilterOperator = {
    label: 'Not Equal',
    name: 'NEQ',
    field: 'string'
};

export const containsFoldOperator: AdvanedFilterOperator = {
    label: 'Contains',
    name: 'ContainsFold',
    field: 'string'
};

export const containsOperator: AdvanedFilterOperator = {
    label: 'Contains (case sensitive)',
    name: 'Contains',
    field: 'string'
};

export const hasPrefixOperator: AdvanedFilterOperator = {
    label: 'Starts with',
    name: 'HasPrefix',
    field: 'string'
};

export const hasSufixOperator: AdvanedFilterOperator = {
    label: 'Ends with',
    name: 'HasSufix',
    field: 'string'
};


//NUMBER
export const eqNumberOperator: AdvanedFilterOperator = {
    label: 'Equal',
    name: '_', // no suffix
    field: 'number'
};

export const neqNumberOperator: AdvanedFilterOperator = {
    label: 'Not Equal',
    name: 'NEQ',
    field: 'number'
};

export const greaterThanNumberOperator: AdvanedFilterOperator = {
    label: 'Greater Than',
    name: 'GT',
    field: 'number'
};


export const greaterThanOrEqualNumberOperator: AdvanedFilterOperator = {
    label: 'Greater Than or Equal',
    name: 'GTE',
    field: 'number'
};

export const lessThanNumberOperator: AdvanedFilterOperator = {
    label: 'Less Than',
    name: 'LT',
    field: 'number'
};

export const lessThanOrEqualNumberOperator: AdvanedFilterOperator = {
    label: 'Less Than or Equal',
    name: 'LTE',
    field: 'number'
};


// BOOLEAN
export const eqBooleanOperator: AdvanedFilterOperator = {
    label: 'Equal',
    name: '_', // no suffix
    field: 'boolean'
};

// DATE
export const greaterThanDateTimeOperator: AdvanedFilterOperator = {
    label: 'Greater Than',
    name: 'GT',
    field: 'date'
};

export const greaterThanOrEqualDateTimeOperator: AdvanedFilterOperator = {
    label: 'Greater Than or Equal',
    name: 'GTE',
    field: 'date'
};

export const lessThanDateTimeOperator: AdvanedFilterOperator = {
    label: 'Less Than',
    name: 'LT',
    field: 'date'
};

export const lessThanOrEqualDateTimeOperator: AdvanedFilterOperator = {
    label: 'Less Than or Equal',
    name: 'LTE',
    field: 'date'
};

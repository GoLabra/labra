import { AdvancedFilter } from "@/core-features/dynamic-filter/filter";
import { containsFoldOperator, eqNumberOperator, eqBooleanOperator } from "@/core-features/dynamic-filter/filter-operators";
import { Field } from "@/lib/apollo/graphql.entities";


export const getAdvancedFiltersFromQuery = (query: string, fields: Field[]): AdvancedFilter[] => {
    if (!query) {
        return [];
    }

    const reusult = [
        ...getStringQuery(query, fields),
        ...getIntegerQuery(query, fields),
        ...getBooleanQuery(query, fields)
    ];

    return reusult;
}

const getStringQuery = (query: string, fields: Field[]): AdvancedFilter[] => {
    return fields.filter(i => 
        i.type == 'ID'
        || i.type == 'ShortText'
        || i.type == 'LongText'
        || i.type == 'RichText'
        || i.type == 'Email'
        || i.type == 'Json'
        || i.type == 'SingleChoice'
    ).map(i => ({
        property: i.name,
        operator: containsFoldOperator.name,
        value: query
    }));
}

const getIntegerQuery = (query: string, fields: Field[]): AdvancedFilter[] => {

    const numberValue = parseInt(query);
    if (isNaN(numberValue)) {
        return [];
    }

    return fields.filter(i => i.type == 'Integer'
        || i.type == 'Decimal'
        || i.type == 'Float'
    ).map(i => ({
        property: i.name,
        operator: eqNumberOperator.name,
        value: numberValue
    }));
}

const getBooleanQuery = (query: string, fields: Field[]): AdvancedFilter[] => {
    const boolField = fields.find(i => i.type == 'Boolean' && i.caption.toLocaleLowerCase() == query.toLocaleLowerCase());
    if(!boolField){
        return [];
    }

    return [{
        property: boolField!.name,
        operator: eqBooleanOperator.name,
        value: true
    }];
}
import { containsFoldOperator, lessThanOrEqualDateTimeOperator, lessThanOrEqualNumberOperator } from "@/core-features/dynamic-filter/filter-operators";
import { Field } from "@/lib/apollo/graphql.entities"
import { ApiFieldTypes, FormFieldTypes } from "@/types/field-type-descriptor";
import { ColumnDefFilter } from "mosaic-data-table";
import { useMemo } from "react";

interface UseDynamicGridFilterProps {
    fields?: Field[],
}
export const useDyamicGridFilter = (props: UseDynamicGridFilterProps): Record<string, ColumnDefFilter | Exclude<ColumnDefFilter["type"], "select">> => {

    return useMemo(() => {

        if(!props.fields){
            return {}
        }

        return props.fields?.reduce((acc: Record<string, ColumnDefFilter | Exclude<ColumnDefFilter["type"], "select">>, field: Field) => {
            const filterDef = fieldToFilter(field);

            if (!filterDef) {
                return acc;
            }

            const [name, filter] = filterDef!;

            if(!filter){
                return acc;
            }
            
            acc[name] = filter;
            return acc;
        }, {});

    }, [props.fields]);
}

const fieldToFilter = (field: Field): [key: string, ColumnDefFilter | Exclude<ColumnDefFilter["type"], "select">] | null => {

    switch (field.type as ApiFieldTypes) {
        case 'ShortText':
        case 'LongText':
        case 'RichText':
        case 'Email':
            return textColumnDef(field.name);
        case 'Integer': 
            return numberColumnDef(field.name);
        case 'DateTime': 
            return dateTimeColumnDef(field.name);
        case 'Date': 
            return dateColumnDef(field.name);
        case 'Time': 
            return timeColumnDef(field.name);
        case 'Boolean': 
            return booleanColumnDef(field.name);
        default:
            return null;
    }
}


const textColumnDef = (name: string): [key: string, ColumnDefFilter | Exclude<ColumnDefFilter["type"], "select">] => {
    return [
        name,
        {
            type: 'string',
            defaultOperator: containsFoldOperator.name,
        }
    ]
}


const numberColumnDef = (name: string): [key: string, ColumnDefFilter | Exclude<ColumnDefFilter["type"], "select">] => {
    return [
        name,
        {
            type: 'number',
            defaultOperator: lessThanOrEqualNumberOperator.name,
        }
    ]
}

const dateColumnDef = (name: string): [key: string, ColumnDefFilter | Exclude<ColumnDefFilter["type"], "select">] => {
    return [
        name,
        {
            type: 'date',
            defaultOperator: lessThanOrEqualDateTimeOperator.name,
        }
    ]
}

const timeColumnDef = (name: string): [key: string, ColumnDefFilter | Exclude<ColumnDefFilter["type"], "select">] => {
    return [
        name,
        {
            type: 'time',
            defaultOperator: lessThanOrEqualDateTimeOperator.name,
        }
    ]
}

const dateTimeColumnDef = (name: string): [key: string, ColumnDefFilter | Exclude<ColumnDefFilter["type"], "select">] => {
    return [
        name,
        {
            type: 'datetime',
            defaultOperator: lessThanOrEqualDateTimeOperator.name,
        }
    ]
}

const booleanColumnDef = (name: string): [key: string, ColumnDefFilter | Exclude<ColumnDefFilter["type"], "select">] => {
    return [
        name,
        {
            type: 'select',
			selectOptions: [{
				label: 'True',
				value: true
			}, {
				label: 'False',
				value: false
			}] as any[]
        }
    ]
}
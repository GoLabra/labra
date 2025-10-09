import { containsFoldOperator, lessThanOrEqualDateTimeOperator, lessThanOrEqualNumberOperator } from "@/core-features/dynamic-filter/filter-operators";
import { useContentManagerSearchFromQuery } from "@/hooks/use-content-manager-search-from-query";
import { Field } from "@/lib/apollo/graphql.entities"
import { ApiFieldTypes, FormFieldTypes } from "@/types/field-type-descriptor";
import { ColumnDefFilter, createFilterRowStore, Filter } from "mosaic-data-table";
import { useMemo } from "react";

interface UseDynamicGridFilterProps {
	search: ReturnType<typeof useContentManagerSearchFromQuery>;
    fields?: Field[];
}
export const useDyamicGridFilter = (props: UseDynamicGridFilterProps) => {

	const filterStore = useMemo(() => createFilterRowStore<any>(props.search.state.filter), []);

    const filterDef = useMemo(() => {

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


	return useMemo(() => ({
		store: filterStore,
		filterDef,
	}), [filterStore, filterDef])
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
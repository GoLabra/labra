import { containsFoldOperator, eqNumberOperator, eqBooleanOperator } from "@/core-features/dynamic-filter/filter-operators";
import { Field } from "@/lib/apollo/graphql.entities";
import { Filter } from "mosaic-data-table";


export const getFiltersFromQuery = (query: string, fields: Field[]): Filter => {
	if (!query) {
		return {};
	}

	const result = {
		...getStringQuery(query, fields),
		...getIntegerQuery(query, fields),
		...getBooleanQuery(query, fields)
	};

	return result;
}

const getStringQuery = (query: string, fields: Field[]): Filter => {

	return fields.reduce((acc, field) => {
		if (['ID', 'ShortText', 'LongText', 'RichText', 'Email', 'Json', 'SingleChoice'].includes(field.type)) {
			acc[field.name] = {
				operator: containsFoldOperator.name,
				value: query
			};
		}
		return acc;
	}, {} as Filter);

	// return fields.filter(i => 
	//     i.type == 'ID'
	//     || i.type == 'ShortText'
	//     || i.type == 'LongText'
	//     || i.type == 'RichText'
	//     || i.type == 'Email'
	//     || i.type == 'Json'
	//     || i.type == 'SingleChoice'
	// ).map(i => ({
	//     property: i.name,
	//     operator: containsFoldOperator.name,
	//     value: query
	// }));
}

const getIntegerQuery = (query: string, fields: Field[]): Filter => {

	const numberValue = parseInt(query);
	if (isNaN(numberValue)) {
		return {};
	}

	return fields.reduce((acc, field) => {
		if (['Integer', 'Decimal', 'Float'].includes(field.type)) {
			acc[field.name] = {
				operator: eqNumberOperator.name,
				value: numberValue
			};
		}
		return acc;
	}, {} as Filter);

	// return fields.filter(i => i.type == 'Integer'
	//     || i.type == 'Decimal'
	//     || i.type == 'Float'
	// ).map(i => ({
	//     property: i.name,
	//     operator: eqNumberOperator.name,
	//     value: numberValue
	// }));
}

const getBooleanQuery = (query: string, fields: Field[]): Filter => {
	const boolField = fields.find(i => i.type == 'Boolean' && i.caption.toLocaleLowerCase() == query.toLocaleLowerCase());
	if (!boolField) {
		return {};
	}

	return fields.reduce((acc, field) => {
		if (field.type == 'Boolean' && field.caption.toLocaleLowerCase() == query.toLocaleLowerCase()) {
			acc[field.name] = {
				operator: eqBooleanOperator.name,
				value: true
			};
		}
		return acc;
	}, {} as Filter);


	// return [{
	// 	property: boolField!.name,
	// 	operator: eqBooleanOperator.name,
	// 	value: true
	// }];
}
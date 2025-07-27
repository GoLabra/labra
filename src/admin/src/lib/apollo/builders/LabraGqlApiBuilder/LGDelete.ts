import { pascalCase } from "change-case";
import { FieldRequest, GplFilter, ILGQuery, ObjectKeys, Unarray, WhereInput } from "./types/types";
import { LGSelectInclude } from "./LGSelectInclude";
import IQueryBuilderOptions from "gql-query-builder/build/IQueryBuilderOptions";
import Fields from "gql-query-builder/build/Fields";
import * as gqlBuilder from 'gql-query-builder'

export class LGDelete<T> implements ILGQuery {

	public readonly isMutation = true;

	private readonly _field: string;
	private readonly _select: FieldRequest<T>[] = [];
	private readonly _filter?: WhereInput<T>;

	private constructor(operation: string, fields: FieldRequest<T>[] , filter?: WhereInput<T>) {
		this._field = operation;
		this._select = fields;
		this._filter = filter;
	}

	public getOperationName = (): string => {
		return `delete${pascalCase(this._field)}`;
	}

	public select = (...fields: FieldRequest<T>[]) => {
		return new LGDelete<T>(this._field, [...this._select, ...fields], this._filter);
	};

	public include<K extends keyof T>(
		key: K extends keyof T ? T[K] extends object ? K : never : never,
		builder: (query: LGSelectInclude<Unarray<T[ObjectKeys<T>]>>) => LGSelectInclude<Unarray<T[ObjectKeys<T>]>>
	): LGDelete<T> {
		const nestedFields = builder(LGSelectInclude.from<Unarray<T[ObjectKeys<T>]>>(key as string));
		return new LGDelete(this._field, [...this._select, nestedFields], this._filter);
	}

	public where = (filter: GplFilter<T>) => {
		const whereInput = !!this._filter ? GplFilter.and(this._filter, filter.expression).expression : filter.expression
		return new LGDelete<T>(this._field, this._select, whereInput);
	};
		
	public buildOptions = (): IQueryBuilderOptions => {

		const operation = this.getOperationName();

		// compute fields
		const fields = this._select.reduce((acc, field) => {
			if(field instanceof LGSelectInclude){
				return acc.concat(field.getField());
			}
			if(typeof field === 'string'){
				acc.push(field);
			}
			return acc;
		}, [] as Fields);

		const queryOptions = {
			operation,
			fields,
			variables: {
				where: {
					type: `${pascalCase(this._field!)}WhereUniqueInput`,
					required: true,
					value: this._filter
				}
			}
		};

		return queryOptions;
	}

	public build = (): { query: string, variables: Record<string, any> } => {
		const queryOptions = this.buildOptions();
		return gqlBuilder.query(queryOptions);
	}

	public getResultData = (response: any) => {
		const fieldName = this.getOperationName();
		return response[fieldName];
	}

	public static from = <T = any>(entityName: string) => {
		return new LGDelete<T>(entityName, [], undefined);
	};
}
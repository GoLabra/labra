import { pascalCase } from "change-case";
import { EntityBaseType, FieldRequest, GplFilter, ILGQuery, ObjectKeys, Unarray } from "./types/types";
import { LGSelectInclude } from "./LGSelectInclude";
import IQueryBuilderOptions from "gql-query-builder/build/IQueryBuilderOptions";
import Fields from "gql-query-builder/build/Fields";
import * as gqlBuilder from 'gql-query-builder'
import pluralize from "pluralize";

export class LGUpsert<T extends EntityBaseType> implements ILGQuery {

	public readonly isMutation = true;
	private readonly _field: string;
	private readonly _select: FieldRequest<T>[] = [];
	private readonly _data: T | T[];

	private constructor(data: T | T[], operation: string, fields: FieldRequest<T>[]) {
		this._data = data;
		this._field = operation;
		this._select = fields;
	}

	public getOperationName = (): string => {
		if(Array.isArray(this._data)){
			return `upsertMany${pascalCase(pluralize(this._field))}`;
		}
		return `upsert${pascalCase(this._field)}`;
	}

	public select = (...fields: FieldRequest<T>[]) => {
		return new LGUpsert<T>(this._data, this._field, [...this._select, ...fields]);
	};

	public include<K extends keyof T>(
		key: K extends keyof T ? T[K] extends object ? K : never : never,
		builder: (query: LGSelectInclude<Unarray<T[ObjectKeys<T>]>>) => LGSelectInclude<Unarray<T[ObjectKeys<T>]>>
	): LGUpsert<T> {
		const nestedFields = builder(LGSelectInclude.from<Unarray<T[ObjectKeys<T>]>>(key as string));
		return new LGUpsert(this._data, this._field, [...this._select, nestedFields]);
	}

	public where = (filter: GplFilter<T>) => {
		return new LGUpsert<T>(this._data, this._field, this._select);
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
				data: {
					type: `Create${pascalCase(this._field)}Input`,
					required: true,
					value: this._data,
					list: Array.isArray(this._data) ? [true] : undefined
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
		if(!response){
			return null;
		}
		
		const fieldName = this.getOperationName();
		return response[fieldName];
	}

	public static from = <T extends EntityBaseType = any>(entityName: string, data: T | T[]) => {
		return new LGUpsert<T>(data, entityName, []);
	};
}
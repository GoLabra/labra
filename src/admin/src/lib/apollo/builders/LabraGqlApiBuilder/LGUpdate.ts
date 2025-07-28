import { pascalCase } from "change-case";
import { EntityBaseType, FieldRequest, GplFilter, ILGQuery, ObjectKeys, Unarray } from "./types/types";
import { LGSelectInclude } from "./LGSelectInclude";
import IQueryBuilderOptions from "gql-query-builder/build/IQueryBuilderOptions";
import Fields from "gql-query-builder/build/Fields";
import * as gqlBuilder from 'gql-query-builder'
import pluralize from "pluralize";

export class LGUpdate<T extends EntityBaseType> implements ILGQuery {

	public readonly isMutation = true;

	private readonly _field: string;
	private readonly _select: FieldRequest<T>[] = [];
	private readonly _filter?: GplFilter<T>;
	private readonly _data: any;

	private constructor(data: any, operation: string, fields: FieldRequest<T>[] , filter?: GplFilter<T>) {
		this._data = data;
		this._field = operation;
		this._select = fields;
		this._filter = filter;
	}

	private isMany = () => {
		return this._filter?.key == null || this._filter?.key == 'id';
	}

	public getOperationName = (): string => {
		if(this.isMany()){
			return `updateMany${pascalCase(pluralize(this._field))}`;
		}
		return `update${pascalCase(this._field)}`;
	}

	public select = (...fields: FieldRequest<T>[]) => {
		return new LGUpdate<T>(this._data, this._field, [...this._select, ...fields], this._filter);
	};

	public include<K extends keyof T>(
		key: K extends keyof T ? T[K] extends object ? K : never : never,
		builder: (query: LGSelectInclude<Unarray<T[ObjectKeys<T>]>>) => LGSelectInclude<Unarray<T[ObjectKeys<T>]>>
	): LGUpdate<T> {
		const nestedFields = builder(LGSelectInclude.from<Unarray<T[ObjectKeys<T>]>>(key as string));
		return new LGUpdate(this._data, this._field, [...this._select, nestedFields], this._filter);
	}

	public where = (filter: GplFilter<T>) => {
		const whereInput = !!this._filter ? GplFilter.and(this._filter, filter) : filter
		return new LGUpdate<T>(this._data, this._field, this._select, whereInput);
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

		let whereType = (() => {
			if(this.isMany()){
				return `${pascalCase(this._field!)}WhereInput`;
			}
			return `${pascalCase(this._field!)}WhereUniqueInput`;
		})();

		const queryOptions = {
			operation,
			fields,
			variables: {
				where: {
					type: whereType,
					required: true,
					value: this._filter
				},
				 data: {
					type: `Update${pascalCase(this._field!)}Input`,
					required: true,
					value: this._data
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

	public static from = <T extends EntityBaseType = any>(entityName: string, data: any) => {
		return new LGUpdate<T>(data, entityName, [], undefined);
	};
}
import * as gqlBuilder from 'gql-query-builder'
import Fields from 'gql-query-builder/build/Fields';
import IQueryBuilderOptions from 'gql-query-builder/build/IQueryBuilderOptions';
import pluralize from "pluralize";
import { EntityBaseType, FieldRequest, GplFilter, GplOrder, ILGQuery, ObjectKeys, Unarray } from './types/types';
import { LGSelectInclude } from './LGSelectInclude';

export class LGConnectionQuery<T extends EntityBaseType> implements ILGQuery {

	public readonly isMutation = false;

	private readonly _field: string;
	private readonly _select: FieldRequest<T>[] = [];
	private readonly _filter?: GplFilter<T>;
	private readonly _order?: GplOrder<T>;
	private readonly _skip?: number;
	private readonly _first?: number;
	private readonly _last?: number;

	private constructor(operation: string, fields: FieldRequest<T>[] , filter?: GplFilter<T>, order?: GplOrder<T>, skip?: number, first?: number, last?: number) {
		this._field = operation;
		this._select = fields;
		this._filter = filter;
		this._order = order;
		this._skip = skip;
		this._first = first;
		this._last = last;
	}

	public select = (...fields: FieldRequest<T>[]) => {
		return new LGConnectionQuery<T>(this._field, [...this._select, ...fields], this._filter, this._order);
	};

	public include<K extends keyof T>(
		key: K extends keyof T ? T[K] extends object ? K : never : never,
		builder: (query: LGSelectInclude<Unarray<T[ObjectKeys<T>]>>) => LGSelectInclude<Unarray<T[ObjectKeys<T>]>>
	): LGConnectionQuery<T> {
		const nestedFields = builder(LGSelectInclude.from<Unarray<T[ObjectKeys<T>]>>(key as string));
		return new LGConnectionQuery(this._field, [...this._select, nestedFields], this._filter, this._order);
	}

	public where = (filter: GplFilter<T>) => {
		const whereInput = !!this._filter ? GplFilter.and(this._filter, filter) : filter;
		return new LGConnectionQuery<T>(this._field, this._select, whereInput);
	};

	public orderByAscending = (field: keyof T) => {		
		return new LGConnectionQuery<T>(this._field, this._select, this._filter, new GplOrder<T>(field, true));
	};

	public orderByDescending = (field: keyof T) => {		
		return new LGConnectionQuery<T>(this._field, this._select, this._filter, new GplOrder<T>(field, false));
	};

	public skip = (skip: number) => {		
		return new LGConnectionQuery<T>(this._field, this._select, this._filter, this._order, skip, this._first, this._last);
	}

	public first = (first: number) => {		
		return new LGConnectionQuery<T>(this._field, this._select, this._filter, this._order, this._skip, first, this._last);
	}

	public last = (last: number) => {		
		return new LGConnectionQuery<T>(this._field, this._select, this._filter, this._order, this._skip, this._first, last);
	}

	public getOperationName = (): string => {
		return `${pluralize(this._field)}Connection`;
	}

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
			fields
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

	public static from = <T extends EntityBaseType = any>(entityName: string) => {
		return new LGConnectionQuery<T>(entityName, [], undefined, undefined, undefined, undefined, undefined);
	};
}
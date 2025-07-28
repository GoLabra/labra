import * as gqlBuilder from 'gql-query-builder'
import Fields from 'gql-query-builder/build/Fields';
import IQueryBuilderOptions from 'gql-query-builder/build/IQueryBuilderOptions';
import pluralize from "pluralize";
import { EntityBaseType, FieldRequest, GplFilter, GplOrder, ILGQuery, ObjectKeys, Unarray } from './types/types';
import { LGSelectInclude } from './LGSelectInclude';
import { LGConnectionQuery } from './LGConnectionQuery';
import { LGDelete } from './LGDelete';
import { LGUpdate } from './LGUpdate';
import { LGCreate } from './LGCreate';
import { pascalCase } from 'change-case';
import { LGUpsert } from './LGUpsert';

export class LGQuery<T extends EntityBaseType>  implements ILGQuery {
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
		return new LGQuery<T>(this._field, [...this._select, ...fields], this._filter, this._order);
	};

	public include<K extends keyof T>(
		key: K extends keyof T ? T[K] extends object ? K : never : never,
		builder: (query: LGSelectInclude<Unarray<T[ObjectKeys<T>]>>) => LGSelectInclude<Unarray<T[ObjectKeys<T>]>>
	): LGQuery<T> {
		const nestedFields = builder(LGSelectInclude.from<Unarray<T[ObjectKeys<T>]>>(key as string));
		return new LGQuery(this._field, [...this._select, nestedFields], this._filter, this._order);
	}

	public where = (filter: GplFilter<T>) => {
		const whereInput = !!this._filter ? GplFilter.and(this._filter, filter) : filter
		return new LGQuery<T>(this._field, this._select, whereInput);
	};

	public orderByAscending = (field: keyof T) => {		
		return new LGQuery<T>(this._field, this._select, this._filter, new GplOrder<T>(field, true), this._skip, this._first, this._last);
	};

	public orderByDescending = (field: keyof T) => {		
		return new LGQuery<T>(this._field, this._select, this._filter, new GplOrder<T>(field, false), this._skip, this._first, this._last);
	};

	public skip = (skip: number) => {		
		return new LGQuery<T>(this._field, this._select, this._filter, this._order, skip, this._first, this._last);
	}

	public first = (first: number) => {		
		return new LGQuery<T>(this._field, this._select, this._filter, this._order, this._skip, first, this._last);
	}

	public last = (last: number) => {		
		return new LGQuery<T>(this._field, this._select, this._filter, this._order, this._skip, this._first, last);
	}

	public getOperationName = (): string => {
		return pluralize(this._field);
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

		// compute filter
		const queryOptions = {
			operation,
			fields,
			variables: {
				where: {
					type: `${pascalCase(this._field)}WhereInput`,
					value: this._filter?.getExpression(),
				},

				...(this._skip != null &&  {
					skip: this._skip
				}),

				...(this._first != null &&  {
					first: this._first
				}),

				...(this._last != null &&  {
					skip: this._last
				}),

				... (this._order != null && {
					orderBy: {
						type: `${pascalCase(this._field)}Order`,
						value: {
							field: this._order.field,
							direction: this._order.ascending ? 'ASC' : 'DESC'
						}
					}
				})
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

	public static from = <T extends EntityBaseType = any>(entityName: string) => {
		return new LGQuery<T>(entityName, [], undefined, undefined, undefined, undefined, undefined);
	};

	public static fromConnection = <T extends EntityBaseType = any>(entityName: string) => {
		return LGConnectionQuery.from<T>(entityName);
	};

	public static update = <T extends EntityBaseType = any>(entityName: string, data: any) => {
		return LGUpdate.from<T>(entityName, data);
	};
	
	public static create = <T extends EntityBaseType = any>(entityName: string, data: any | any[]) => {
		return LGCreate.from<T>(entityName, data);
	};

	public static upsert = <T extends EntityBaseType = any>(entityName: string, data: any | any[]) => {
		return LGUpsert.from<T>(entityName, data);
	};
	
	public static deleteFrom = <T extends EntityBaseType = any>(entityName: string) => {
		return LGDelete.from<T>(entityName);
	};

	public static merge = <T = any>(...queries: ILGQuery[]) => {
		
		if(!queries.length){
			throw new Error("No queries to merge");
		}

		if(queries[0].isMutation == false){
			return gqlBuilder.query(queries.map(i => i.buildOptions()));
		} 
	
		return gqlBuilder.mutation(queries.map(i => i.buildOptions()));
	};
}

// Example data model

type Permission = {
	id: string;
	name2: string;
	role: {
		id: string;
		name2dsdf: string;
	};
};

type User = {
	id: string;
	name: string;
	permissions: Permission[];
};

// Example usage

const query = LGQuery.from<User>('user')
	.where(
		GplFilter.or(
			GplFilter.field('id', '', '123'),
			GplFilter.and(
				GplFilter.field('name', 'ContainsFold', 'john'),
				GplFilter.field('id', 'NotIn', ['456', '789'])
			)
		)
	)
	.orderByAscending('id')
	.skip(0)
	.first(10)
	.include('permissions', q => q.select('id', 'name2')
									.include('role', q => q.select('id', 'name2dsdf')
								)
			)
	.select('id', 'name');

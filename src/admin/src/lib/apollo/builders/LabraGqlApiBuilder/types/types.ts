import IQueryBuilderOptions from "gql-query-builder/build/IQueryBuilderOptions";
import { LGSelectInclude } from "../LGSelectInclude";


export type Unarray<T> = T extends Array<infer U> ? U : T;
export declare type Maybe<T> = null | undefined | T;


export type FilterOperators =
	| ''
	| 'NEQ'
	| 'In'
	| 'NotIn'
	| 'GT'
	| 'GTE'
	| 'LT'
	| 'LTE'
	| 'EqualFold'
	| 'ContainsFold';

export type EntityBaseType = { id: string }

export type ObjectKeys<T> = {
    [K in keyof T]: T[K] extends Maybe<object> ? K : never
}[keyof T];

export type FieldRequest<T = any> =
    | keyof T
    | LGSelectInclude<any>;

// Query Filter class
export class GplFilter<T extends EntityBaseType> {

	private constructor(key?: Extract<keyof T, string>, operator?: FilterOperators, value?: any, not?: boolean, and?: GplFilter<T>[], or?: GplFilter<T>[]) {
		this.key = key;
		this.operator = operator;
		this.value = value;
		this.not = not ?? false;
		this.and = and;
		this.or = or;
	}

	public key?: Extract<keyof T, string>
	public operator?: FilterOperators
	public value?: any
	public not: boolean = false;
	
	public and?: GplFilter<T>[];
	public or?: GplFilter<T>[];

	static and<T extends EntityBaseType>(...filters: GplFilter<T>[]): GplFilter<T> {
		return new GplFilter<T>(undefined, undefined, undefined, undefined, filters); 
	}

	static or<T extends EntityBaseType>(...filters: GplFilter<T>[]): GplFilter<T> {
		return new GplFilter<T>(undefined, undefined, undefined, undefined, undefined, filters); 
	}

	static not<T extends EntityBaseType>(
		key: Extract<keyof T, string>,
		operator: FilterOperators,
		value: any
	): GplFilter<T> {
		return new GplFilter<T>(key, operator, value, true);
	}

	static field<T extends EntityBaseType>(
		key: Extract<keyof T, string>,
		operator: FilterOperators,
		value: any
	): GplFilter<T> {
		return new GplFilter<T>(key, operator, value);
	}

	public getExpression = (): any => {

		if(this.key){
			const field = {[`${this.key}${this.operator}`]: this.value};
			if(this.not){
				return {not: field};
			}
			return field;
		}

		if(this.and){
			return {and: this.and.map(i => i.getExpression())};
		}

		if(this.or){
			return {or: this.or.map(i => i.getExpression())};
		}

		return undefined;
	}
}

export class GplOrder<T> {
	public ascending: boolean;
	public field: keyof T;

	constructor(field: keyof T, ascending: boolean) {
		this.field = field;
		this.ascending = ascending;
	}
}

export interface ILGQuery {
	readonly isMutation: boolean;
	getOperationName: () => string;
	buildOptions: () => IQueryBuilderOptions;
	build: () => { query: string, variables: Record<string, any> };
	getResultData: (response: any) => any;
}
import IQueryBuilderOptions from "gql-query-builder/build/IQueryBuilderOptions";
import { LGSelectInclude } from "../LGSelectInclude";
import { string } from "zod";

export type Unarray<T> = T extends Array<infer U> ? U : T;

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

export type Filters<T> = {
	[K in keyof T as `${Extract<K, string>}${FilterOperators}`]?: T[K] | T[K][];
};

export type WhereInput<T = any> = Filters<T> & {
	not?: WhereInput<T>;
	and?: WhereInput<T>[];
	or?: WhereInput<T>[];
};

export type ObjectKeys<T> = {
    [K in keyof T]: T[K] extends object ? K : never
}[keyof T];

export type FieldRequest<T = any, K extends ObjectKeys<T> = ObjectKeys<T>> =
	| keyof T
	| LGSelectInclude<Unarray<T[K]>>;

// Query Filter class

export class GplFilter<T> {
	private constructor(public readonly expression: WhereInput<T>) { }

	static and<T>(...filters: Array<GplFilter<T> | WhereInput<T> | null | undefined>): GplFilter<T> {
		return new GplFilter<T>({ and: filters.filter(i => !!i).map(f => {
			if(f instanceof GplFilter){
				f.expression
			}
			return f;
		}) } as WhereInput<T>);
	}

	static or<T>(...filters: Array<GplFilter<T> | null | undefined>): GplFilter<T> {
		return new GplFilter<T>({ or: filters.filter(i => !!i).map(f => {
			if(f instanceof GplFilter){
				f.expression
			}
			return f;
		}) } as WhereInput<T>);
	}

	static not<T>(filter: GplFilter<T>): GplFilter<T> {
		return new GplFilter<T>({ not: filter.expression } as WhereInput<T>);
	}

	static field<T, K extends Extract<keyof T, string>, Op extends FilterOperators, V = T[K]>(
		key: K,
		operator: Op,
		value: Op extends 'In' | 'NotIn' ? V[] : V
	): GplFilter<T> {
		const fieldKey = `${key}${operator}` as keyof WhereInput<T>;
		return new GplFilter<T>({ [fieldKey]: value } as WhereInput<T>);
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
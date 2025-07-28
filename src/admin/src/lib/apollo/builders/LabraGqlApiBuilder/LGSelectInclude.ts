import NestedField from "gql-query-builder/build/NestedField";
import { FieldRequest, ObjectKeys, Unarray } from "./types/types";
import Fields from "gql-query-builder/build/Fields";

export class LGSelectInclude<T> {
	private readonly _field: string;
	private readonly _select: FieldRequest<T>[] = [];

	private constructor(operation: string, fields: FieldRequest<T>[]) {
		this._field = operation;
		this._select = fields;
	}

	public static from = <T = any>(operation: string) => {
		return new LGSelectInclude<T>(operation, []);
	}

	public select = (...fields: FieldRequest<T>[]) => {
		return new LGSelectInclude<T>(this._field, [...this._select, ...fields]);
	};

	public include<K extends keyof T>(
		key: K extends keyof T ? T[K] extends object ? K : never : never,
		builder: (query: LGSelectInclude<Unarray<T[ObjectKeys<T>]>>) => LGSelectInclude<Unarray<T[ObjectKeys<T>]>>
	): LGSelectInclude<T> {
		const nestedFields = builder(LGSelectInclude.from<Unarray<T[ObjectKeys<T>]>>(key as string));
		return new LGSelectInclude(this._field, [...this._select, nestedFields]);
	}

	public getField = (): NestedField => {

		const fields = this._select.reduce((acc, field) => {
			if(field instanceof LGSelectInclude){
				return acc.concat(field.getField());
			}

			if(typeof field === 'string'){
				acc.push(field);
			}

			return acc;
		}, [] as Fields);

		return {
			operation: this._field,
			fields,
			variables: []
		}
	}
}

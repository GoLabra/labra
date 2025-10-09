export type GenericEvent<T = any> = {
	target: {
		name: string;
		value: T;
	}
}

export const createGenericEvent = <T = any>(name: string, value: T): GenericEvent<T> => {
	return {
		target: {
			name,
			value
		}
	}
}
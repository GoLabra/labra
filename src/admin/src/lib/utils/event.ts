export type GenericEvent<T = any> = {
	target: {
		name: string;
		value: T;
	}
}
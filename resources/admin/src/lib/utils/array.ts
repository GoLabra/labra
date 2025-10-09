export const arrayDistinct = (array: Array<any>) => {
	let map = new Map();
	array.forEach(value => map.set(value, value));
	return [...map.values()];
};

export function groupByMap<T, K extends PropertyKey, V>(
	array: T[],
	keyFn: (item: T) => K,
	mapFn: (item: T) => V
): Record<K, V[]> {
	return array.reduce((acc, item) => {
		const key = keyFn(item);
		if (!acc[key]) acc[key] = [];
		acc[key].push(mapFn(item));
		return acc;
	}, {} as Record<K, V[]>);
}
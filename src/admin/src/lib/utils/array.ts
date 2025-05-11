export const arrayDistinct = (array: Array<any>) => {
    let map = new Map();
    array.forEach(value => map.set(value, value));
    return [...map.values()];
};
export const objFromArray = (arr: any[], key = 'id'): {} => arr.reduce((
  accumulator,
  current
) => {
  accumulator[current[key]] = current;
  return accumulator;
}, {});

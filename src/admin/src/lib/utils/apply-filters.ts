import type { AdvancedFilter } from '../../core-features/dynamic-filter/filter';

// NOTE: In really you'll use server filter, but you can also use this to filter the data on client.

const contains = (rowValue: string, filterValue: string): boolean => {
  if (!rowValue) {
    return false;
  }

  return rowValue.toLowerCase().includes(filterValue.toLowerCase());
};

const endsWith = (rowValue: string, filterValue: string): boolean => {
  if (!rowValue) {
    return false;
  }

  return rowValue.substring(rowValue.length - filterValue.length) === filterValue;
};

const equals = (rowValue: string | number, filterValue: string | number): boolean => {
  if (!rowValue) {
    return false;
  }

  // Here we evaluate == instead of === because values can be number | string
  // eslint-disable-next-line eqeqeq
  return rowValue == filterValue;
};

const greaterThan = (rowValue: number, filterValue: number): boolean => {
  if (!rowValue) {
    return false;
  }

  return rowValue > filterValue;
};

const isAfter = (rowValue: Date | number, filterValue: Date | number): boolean => {
  if (!rowValue) {
    return false;
  }

  return new Date(rowValue).getTime() > new Date(filterValue).getTime();
};

const isBlank = (rowValue?: any): boolean => {
  if (rowValue === null || typeof rowValue === 'undefined') {
    return true;
  }

  return false;
};

const isPresent = (rowValue?: any): boolean => {
  if (rowValue === null || typeof rowValue === 'undefined') {
    return false;
  }

  return true;
};

const lessThan = (rowValue: number, filterValue: number): boolean => {
  if (!rowValue) {
    return false;
  }

  return rowValue < filterValue;
};

const isBefore = (rowValue: Date | number, filterValue: Date | number): boolean => {
  if (!rowValue) {
    return false;
  }

  return new Date(rowValue).getTime() < new Date(filterValue).getTime();
};

const notContains = (rowValue: string, filterValue: string): boolean => {
  if (!rowValue) {
    return false;
  }

  return !rowValue.includes(filterValue);
};

const notEqual = (rowValue: string | number, filterValue: string | number): boolean => {
  if (!rowValue) {
    return false;
  }

  return rowValue !== filterValue;
};

const startsWith = (rowValue: string, filterValue: string): boolean => {
  if (!rowValue) {
    return false;
  }

  return rowValue.substring(0, filterValue.length) === filterValue;
};

const checkFilter = (rowValue: any, filter: AdvancedFilter): boolean => {
  switch (filter.operator) {
    case 'contains':
      return contains(rowValue as string, filter.value as string);

    case 'endsWith':
      return endsWith(rowValue as string, filter.value as string);

    case 'equals':
      return equals(rowValue as string | number, filter.value as number);

    case 'greaterThan':
      return greaterThan(rowValue as number, filter.value as number);

    case 'isAfter':
      return isAfter(rowValue as Date | number, filter.value as Date);

    case 'isBefore':
      return isBefore(rowValue as Date | number, filter.value as Date);

    case 'isBlank':
      return isBlank(rowValue);

    case 'isPresent':
      return isPresent(rowValue);

    case 'lessThan':
      return lessThan(rowValue as number, filter.value as number);

    case 'notContains':
      return notContains(rowValue as string, filter.value as string);

    case 'notEqual':
      return notEqual(rowValue as string | number, filter.value as number);

    case 'startsWith':
      return startsWith(rowValue as string, filter.value as string);

    default:
      throw new Error('Provided an unknown filter operator');
  }
};

export const applyFilters = <T = any>(rows: T[], filters: AdvancedFilter[] = []): T[] => {
  if (filters.length === 0) {
    return rows;
  }

  return rows.filter((row: any) => {
    let passedAll = true;

    for (let i = 0; i < filters.length; i++) {
      const filter = filters[i];
      const rowValue = filter.property? row[filter.property] : null;
      let passed = true;

      try {
        passed = checkFilter(rowValue, filter);
      } catch (err) {
        console.warn('[Apply Filters] Skipped filter due to error', err);
      }

      if (!passed) {
        passedAll = false;
        break;
      }
    }

    return passedAll;
  });
};

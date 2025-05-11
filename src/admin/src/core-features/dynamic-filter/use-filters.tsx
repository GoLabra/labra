import { useCallback, useEffect, useMemo, useState } from 'react';
import type { AdvancedFilter, AdvanedFilterOperator, AdvancedFilterProperty, AdvancedFilterValue } from './filter';

const emptyFilter: AdvancedFilter = {
  operator: undefined,
  property: '',
  value: undefined
};

const validateFilter = (filter: AdvancedFilter): boolean => {
  // Filters need an operator and a property
  if (!filter.operator || !filter.property) {
    return false;
  }

  // isBlank and isPresent operators cannot have a value
  if (filter.operator === 'isBlank' || filter.operator === 'isPresent') {
    return typeof filter.value === 'undefined';
  }

  // Other operators require a value
  if (typeof filter.value === 'undefined') {
    return false;
  }

  return true;
};

interface FiltersController {
  filters: AdvancedFilter[];
  handleFilterAdd: (index: number) => void;
  handleFilterRemove: (index: number) => void;
  handleFiltersClear: () => void;
  handleOperatorChange: (index: number, name: string) => void;
  handlePropertyChange: (index: number, name: string) => void;
  handleValueChange: (index: number, value: AdvancedFilterValue) => void;
  valid: boolean;
}

export const useFilters = (
  //operators: FilterOperator[] = [],
  properties: AdvancedFilterProperty[] = [],
  initialFilters: AdvancedFilter[] = []
): FiltersController => {
  const [filters, setFilters] = useState<AdvancedFilter[]>([]);

  useEffect(
    () => {
      setFilters(
        initialFilters.length > 0
          ? initialFilters
          : [emptyFilter]
      );
    },
    [initialFilters]
  );

  const valid = useMemo(
    () => {
      let passedAll = true;

      for (let i = 0; i < filters.length; i++) {
        const passed = validateFilter(filters[i]);

        if (!passed) {
          passedAll = false;
          break;
        }
      }

      return passedAll;
    },
    [filters]
  );

  const handleFilterAdd = useCallback(
    (index: number): void => {
      setFilters((prevState) => {
        const filters = [...prevState];

        filters.splice(index, 0, emptyFilter);

        return filters;
      });
    },
    []
  );

  const handleOperatorChange = useCallback(
    (index: number, name: string): void => {
      // Ensure operator is allowed

    //   const operator = operators.find((operator) => operator.name === name);

    //   if (!operator) {
    //     return;
    //   }

      setFilters((prevState) => {
        const filters = [...prevState];

        filters[index] = {
          ...filters[index],
          operator: name
        };

        return filters;
      });
    },
    [setFilters]
  );

  const handlePropertyChange = useCallback(
    (index: number, name: string): void => {
      // Ensure property is allowed

      const property = properties.find((property) => property.name === name);

      if (!property) {
        return;
      }

      setFilters((prevState) => {
        const filters = [...prevState];

        filters[index] = {
          operator: undefined,
          property: name,
          value: undefined
        };

        return filters;
      });
    },
    [properties]
  );

  const handleValueChange = useCallback(
    (index: number, value: AdvancedFilterValue): void => {
      setFilters((prevState) => {
        const filters = [...prevState];

        filters[index] = {
          ...filters[index],
          value
        };

        return filters;
      });
    },
    []
  );

  const handleFiltersClear = useCallback(
    (): void => {
      setFilters([emptyFilter]);
    },
    []
  );

  const handleFilterRemove = useCallback(
    (index: number): void => {
      if (filters.length === 1) {
        setFilters([emptyFilter]);
        return;
      }

      setFilters((prevState) => {
        return prevState.filter((filter, _index) => _index !== index);
      });
    },
    [filters]
  );

  return {
    filters,
    handleFilterAdd,
    handleFilterRemove,
    handleFiltersClear,
    handleOperatorChange,
    handlePropertyChange,
    handleValueChange,
    valid
  };
};

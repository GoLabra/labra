import { FilterOperators } from "lg-query";

export interface AdvanedFilterOperator {
    label: string;
    name: FilterOperators;
    field?: 'date' | 'string' | 'number' | 'boolean';
}

// export interface AdvancedFilterProperty {
//     label: string;
//     name: string;
//     // List of operator names
//     operators: AdvanedFilterOperator[];
// }

// export type AdvancedFilterValue = Date | string | number | boolean | undefined;

// export interface AdvancedFilter {
//     // Property name
//     property?: string;
//     // Operator name
//     operator?: string;
//     // Matching value
//     value: AdvancedFilterValue;
// }

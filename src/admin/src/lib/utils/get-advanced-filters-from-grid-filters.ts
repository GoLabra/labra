import { AdvancedFilter } from "@/core-features/dynamic-filter/filter";
import { Field } from "@/lib/apollo/graphql.entities";
import { Filter } from "mosaic-data-table";

export const getAdvancedFiltersFromGridFilter = (filter: Filter): AdvancedFilter[] => {

    const result = Object.entries(filter).map(([key, value]): AdvancedFilter => {
        
        return {
            property: key,
            operator: value.operator,
            value: value.value
        };
    });

    return result;
}
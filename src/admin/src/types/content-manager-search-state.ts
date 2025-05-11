import { Filter, Order } from "mosaic-data-table";
import { AdvancedFilter as AdvancedFilters } from "../core-features/dynamic-filter/filter";

export interface ContentManagerSearchState {
    //entityName: string;
    page: number;
    rowsPerPage: number;
    sortBy: string | null;
    order: Order;
    query: string;
    advancedFilters: AdvancedFilters[];
    filter: Filter;
}



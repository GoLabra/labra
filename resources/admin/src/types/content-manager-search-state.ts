import { Filter, Order } from "mosaic-data-table";

export interface ContentManagerSearchState {
    //entityName: string;
    page: number;
    rowsPerPage: number;
    sortBy: string | null;
    order: Order;
    query: string;
    filter: Filter;
}



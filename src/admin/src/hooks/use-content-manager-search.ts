import { ContentManagerSearchState } from "@/types/content-manager-search-state";
import { useCallback, useMemo, useState } from "react";
import { AdvancedFilter as AdvancedFilters } from "@/core-features/dynamic-filter/filter";
import { Filter, Order } from "mosaic-data-table";

interface UseContentManagerSearchParams {
    initialFilter?: Filter;
}
export const useContentManagerSearch = (props?: UseContentManagerSearchParams) => {

    const [page, setPage] = useState<number>(1);
    const [rowsPerPage, setRowsPerPage] = useState<number>(10);
    const [sortBy, setSortBy] = useState<string | null>(null)
    const [order, setOrderBy] = useState<string>('desc')
    const [advancedFilters, setAdvancedFilters] = useState<AdvancedFilters[] | null>([]);
    const [filter, setFilter] = useState<Filter|null>(props?.initialFilter ?? null);
    const [query, setQuery] = useState<string>('');

    const handleAdvancedFiltersApply = useCallback((filters: AdvancedFilters[]): void => {
        setAdvancedFilters(filters);
        setPage(1);
    }, [setAdvancedFilters, setPage]);

    const handleAdvancedFiltersClear = useCallback((): void => {
        setAdvancedFilters(null);
        setPage(1);
    }, [setAdvancedFilters, setPage]);


    const handleFiltersApply = useCallback((filter: Filter): void => {
        setFilter(filter);
        setPage(1);
    }, [setFilter, setPage]);

    const handleFiltersClear = useCallback((): void => {
        setFilter(null);
        setPage(1);
    }, [setFilter, setPage]);

    const handlePageChange = useCallback((value: number): void => {
        setPage(value);
    }, [setPage]);

    const handleRowsPerPageChange = useCallback((value: number): void => {
        setRowsPerPage(value);
        setPage(1);
    }, [setRowsPerPage, setPage]);

    const handleSortChange = useCallback((sortBy: string, order: Order): void => {
        setSortBy(sortBy);
        setOrderBy(order);
    }, [setSortBy, setOrderBy]);

    const handleQueryChange = useCallback((query: string): void => {
        setQuery(query);
        setPage(1);
    }, [setPage, setQuery]);

    const searchState: ContentManagerSearchState = useMemo(() => ({
        page,
        rowsPerPage,
        sortBy,
        order: order as Order,
        query,
        advancedFilters: advancedFilters ?? [],
        filter: filter ?? {},
    }), [page, rowsPerPage, sortBy, order, advancedFilters, filter, query]);

    return useMemo(() => ({
        handleQueryChange,
        handleAdvancedFiltersApply: handleAdvancedFiltersApply,
        handleAdvancedFiltersClear: handleAdvancedFiltersClear,
        handleFiltersApply: handleFiltersApply,
        handleFiltersClear: handleFiltersClear,
        handlePageChange,
        handleRowsPerPageChange,
        handleSortChange,
        state: searchState
    }), [handleQueryChange, handleAdvancedFiltersApply, handleAdvancedFiltersClear, handlePageChange, handleRowsPerPageChange, handleSortChange, searchState]);
};
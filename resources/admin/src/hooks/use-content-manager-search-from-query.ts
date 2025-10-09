import { ContentManagerSearchState } from "@/types/content-manager-search-state";
import { useCallback, useMemo } from "react";
import { useQueryState, parseAsInteger, parseAsString, parseAsJson, parseAsBoolean } from 'nuqs';
import { Filter, Order } from "mosaic-data-table";
import Defaults from "@/config/Defaults.json";


export const useContentManagerSearchFromQuery = () => {

    const [page, setPage] = useQueryState('p', parseAsInteger.withDefault(1))
    const [rowsPerPage, setRowsPerPage] = useQueryState('rpp', parseAsInteger.withDefault(Defaults.dataTable.rowsPerPage))
    const [sortBy, setSortBy] = useQueryState('s')
    const [order, setOrderBy] = useQueryState('o', parseAsString.withDefault('desc'))
    const [filter, setFilter] = useQueryState('f', parseAsJson<Filter>());
    const [query, setQuery] = useQueryState('q', parseAsString.withDefault(''));

    const handleFiltersApply = useCallback((filter: Filter): void => {
        const newFilter = Object.keys(filter).length == 0 ? null : filter;
        setFilter(newFilter);
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
        filter: filter ?? {},
    }), [page, rowsPerPage, sortBy, order, filter, query]);

    return useMemo(() => ({
        handleQueryChange,
        handleFiltersApply: handleFiltersApply,
        handleFiltersClear: handleFiltersClear,
        handlePageChange,
        handleRowsPerPageChange,
        handleSortChange,
        state: searchState
    }), [handleQueryChange, handlePageChange, handleRowsPerPageChange, handleSortChange, searchState]);
};
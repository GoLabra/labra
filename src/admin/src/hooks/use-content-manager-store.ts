import { ContentManagerSearchState } from "@/types/content-manager-search-state";
import { UsersStore } from "@/types/content-manager-store-state";
import { useCallback, useEffect, useMemo, useState } from "react";
import { useContentManagerStoreRequest } from "./use-content-manager-store-request";
import { FullEntity } from "@/types/entity";
import { addNotification } from "@/lib/notifications/store";
import { Edge, Field } from "@/lib/apollo/graphql.entities";
import { getAdvancedFiltersFromGridFilter } from "@/lib/utils/get-advanced-filters-from-grid-filters";
import { getAdvancedFiltersFromQuery } from "@/lib/utils/get-filters-from-query";
import { EdgeRequest } from "@/lib/apollo/builders/gqlQueryBuilder";
import { Order } from "mosaic-data-table";
import { AdvancedFilter } from "@/core-features/dynamic-filter/filter";

interface UseContentManagerStoreParams {
    entityName: string
    //fullEntity?: FullEntity | null,

    page: number;
    rowsPerPage: number;
    sortBy: string | null;
    order: Order;
    
    fields?: string[],
    edges?: EdgeRequest[],
    // searchState: ContentManagerSearchState;

    filters?: AdvancedFilter[];
    orFilters?: AdvancedFilter[];

    lazy?: boolean;
}
export const useContentManagerStore = (params: UseContentManagerStoreParams): UsersStore => {

    const { entityName, fields, edges, page, rowsPerPage, sortBy, order, filters, orFilters, lazy } = params;

    const contentManagerStoreRequest = useContentManagerStoreRequest({entityName});
    const [dataLoading, setDataLoading] = useState<boolean>(false);
    const [data, setData] = useState<[]>([]);
    const [dataConnection, setDataConnection] = useState<{ totalCount: number } | null>(null);

    const fetch = useCallback(() => {
        setDataLoading(true);
        contentManagerStoreRequest.fetchData({

            page: page ?? 1,
            rowsPerPage: rowsPerPage ?? 10,
            sortBy: sortBy ?? null,
            order: order ?? 'desc',            
            fields: fields ?? [],
            edges: edges ?? [],
            
            // fields: fields?.map(i => i.name) ?? [],
            // edges: edges?.filter(i => i.relationType !== 'ManyToMany')
            //                         .filter(i => i.relationType !== 'ManyToOne')
            //                         .filter(i => i.relationType !== 'Many')
            //                         .map(i => ({
            //                             name: i.name,
            //                             fields: [i.relatedEntity.displayField.name],
            //                         })) ?? [],
            // searchState: searchState!,
            filters: filters ?? [], // getAdvancedFiltersFromGridFilter(searchState.filter),
            orFilters: orFilters ?? [] //getAdvancedFiltersFromQuery(searchState.query, fields ?? [])
        })
            ?.then((response: any) => {
                setData(response.data as []);
                setDataConnection(response.connection ?? { totalCount: 0 });
            }).finally(() => {
                setDataLoading(false);
            });
    }, [contentManagerStoreRequest, setData, setDataLoading, fields, edges, page, rowsPerPage, sortBy, order, filters, orFilters]);

    const refresh = useCallback(() => {
        fetch();
    }, [fetch]);
    
    const addItem = useCallback((data: any) => {
        contentManagerStoreRequest.addItem(entityName,data)
            ?.then((response) => {
                fetch();
            })
            ?.catch(({message}) => {
                addNotification({ message, type: 'error' });
            });
    }, [ contentManagerStoreRequest.addItem, entityName, fetch, setDataLoading]);


    const addItems = useCallback((data: any[]) => {
        const promises = data.map(item => {
            return contentManagerStoreRequest.addItem(entityName, item)
                .catch(({message}) => {
                    addNotification({ message, type: 'error' });
                });
        });

        Promise.all(promises)
            .then(() => {
                fetch();
            });
    }, [contentManagerStoreRequest.addItem, entityName, fetch, setDataLoading]);

    const updateItem = useCallback((id: string, data: any) => {
        contentManagerStoreRequest.updateItem(entityName, id, data)
            ?.then((response) => {
                fetch();
            })
            ?.catch(({message}) => {
                addNotification({ message, type: 'error' });
            });
    }, [fetch, entityName, setDataLoading]);

    const deleteItem = useCallback((id: string) => {
        contentManagerStoreRequest.deleteItem(entityName, id)
            ?.then((response) => {
                fetch();
            })
            ?.catch(({message}) => {
                addNotification({ message, type: 'error' });
            });
    }, [fetch, entityName, setDataLoading]);

    const deleteBulk = useCallback((ids: Array<string>) => {
        contentManagerStoreRequest.deleteBulk(entityName, ids)
            ?.then((response) => {
                fetch();
            })
            ?.catch(({message}) => {
                addNotification({ message, type: 'error' });
            });
    }, [entityName, fetch]);

    useEffect(() => {
        if(lazy){
            return;
        }

        if(!fields && !edges){
            return;
        }

        // if(!fullEntity){
        //     return;
        // }

        // if (fullEntity.loading) {
        //     return;
        // }

        // if (fullEntity.isError) {
        //     return;
        // }
        
        fetch();
    }, [fields, edges, fetch]);

    const numberOfPages = useMemo(() => {
        if (!dataConnection) {
            return 1;
        }

        if (!dataConnection.totalCount) {
            return 1;
        }

        if (!rowsPerPage) {
            return 1;
        }

        return Math.ceil(dataConnection.totalCount / rowsPerPage);
    }, [dataConnection, rowsPerPage]);

    return useMemo(() => ({
        state: {
            data: data,
            // schemaLoading: !fullEntity || fullEntity.loading,
            // dataLoading: !fullEntity || fullEntity.loading || dataLoading,
            // entityFields: fullEntity?.fields ?? [],
            // entityEdges: fullEntity?.edges ?? [],
            dataLoading,

            dataConnection,
            pagesCount: numberOfPages,
            totalItems: dataConnection?.totalCount ?? 0
        },
        refresh,
        addItem,
        addItems,
        updateItem,
        deleteItem,
        deleteBulk
    }), [data,  dataLoading, dataConnection?.totalCount, addItem, dataConnection, deleteBulk, deleteItem, numberOfPages, updateItem]);
}
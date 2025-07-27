import { ContentManagerSearchState } from "@/types/content-manager-search-state";
import { UsersStore } from "@/types/content-manager-store-state";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { FullEntity } from "@/types/entity";
import { addNotification } from "@/lib/notifications/store";
import { Edge, EntityOwner, Field } from "@/lib/apollo/graphql.entities";
import { getAdvancedFiltersFromGridFilter } from "@/lib/utils/get-advanced-filters-from-grid-filters";
import { getAdvancedFiltersFromQuery } from "@/lib/utils/get-filters-from-query";
import { EdgeRequest } from "@/lib/apollo/builders/gqlQueryBuilder";
import { Order } from "mosaic-data-table";
import { AdvancedFilter } from "@/core-features/dynamic-filter/filter";
import { useApolloClient } from "@apollo/client";
import { LGQuery } from "@/lib/apollo/builders/LabraGqlApiBuilder/LGQuery";
import { GplFilter } from "@/lib/apollo/builders/LabraGqlApiBuilder/types/types";
import { RunQuery } from "./use-lg-query";

interface UseGridContentManagerStoreParams {
	fullEntity?: FullEntity;
	searchState: ContentManagerSearchState;
}
export const useContentManagerStore = (params: UseGridContentManagerStoreParams): UsersStore => {

    const { name:entityName } = params.fullEntity ?? {};

    //const contentManagerStoreRequest = useContentManagerStoreRequest({entityName, entityOwner: owner ?? EntityOwner.User});
	const client = useApolloClient();
    const [dataLoading, setDataLoading] = useState<boolean>(false);
    const [data, setData] = useState<[]>([]);
    const [dataConnection, setDataConnection] = useState<{ totalCount: number } | null>(null);

	const fields = useMemo(() => params.fullEntity?.fields.map(i => i.name), [params.fullEntity?.fields]);
	const edges = useMemo(() => params.fullEntity?.edges.filter(i => i.relationType !== 'ManyToMany')
													.filter(i => i.relationType !== 'ManyToOne')
													.filter(i => i.relationType !== 'Many'), [params.fullEntity?.edges]);

	const dataQuery = useMemo(() => {

		if(!entityName){
			return null;
		}

		if(!fields || !edges){
			return null;
		}

		let query = LGQuery.from<any>(entityName)
							.select( ...fields)
		// add edges							
		query = edges.reduce((query: LGQuery<any>, edge) => {
			return query.include(edge.name, q => q.select('id', edge.relatedEntity.displayField.name));
		}, query);

		// add filters
		query = Object.entries(params.searchState.filter).reduce((query, [key, value]) => { 
			return query.where(GplFilter.field(key, value.operator, value.value));
		}, query);

		// add skip
		query = query.skip(params.searchState.page * params.searchState.rowsPerPage);
		// add first
		query = query.first(params.searchState.rowsPerPage);
		// add order
		if(params.searchState.sortBy) {
			if(params.searchState.order == 'asc'){
				query = query.orderByAscending(params.searchState.sortBy);
			} else {
				query = query.orderByDescending(params.searchState.sortBy);
			}
		}
		
		return query;	
	}, [entityName, fields, edges, params.searchState]);

	const connectionQuery = useMemo(() => {
		if(!entityName){
			return null!;
		}

		let query = LGQuery.fromConnection<any>(entityName)
			     .select('totalCount');

		// add filters
		query = Object.entries(params.searchState.filter).reduce((query, [key, value]) => { 
			return query.where(GplFilter.field(key, value.operator, value.value));
		}, query);

		return query;
	}, [entityName]);

    const fetch = useCallback(() => {

		if(!params.fullEntity){
			return;
		}

		if(!dataQuery){
			return;
		}
		
        setDataLoading(true);

		const apiType = params.fullEntity.owner == EntityOwner.Admin ? 'admin' : 'user';

		RunQuery(client, apiType, dataQuery, connectionQuery).then((response: any) => {
			setData(dataQuery.getResultData(response.data) as []); 
			setDataConnection(connectionQuery.getResultData(response.data));
		}).finally(() => {
			setDataLoading(false);
		});

    }, [client, params.fullEntity, dataQuery, connectionQuery, setData, setDataLoading]);

    const refresh = useCallback(() => {
        fetch();
    }, [fetch]);
    
    const addItem = useCallback((data: any) => {
		if(!params.fullEntity){
			return;
		}

		const query = LGQuery.create(params.fullEntity.name, data)
							.select('id');

		setDataLoading(true);

		const apiType = params.fullEntity.owner == EntityOwner.Admin ? 'admin' : 'user';
       	RunQuery(client, apiType, query)
		.then((response) => {
			refresh();
		}).catch(() => {
			setDataLoading(false);
		});
			
    }, [client, params.fullEntity, refresh, setDataLoading]);


    const addItems = useCallback((data: any[]) => {
	
    }, []);

    const updateItem = useCallback((id: string, data: any) => {

		if(!params.fullEntity){
			return;
		}

		const query = LGQuery.update(params.fullEntity.name, data)
							.where(GplFilter.field('id', '', id))
							.select('id');

		setDataLoading(true);

		const apiType = params.fullEntity.owner == EntityOwner.Admin ? 'admin' : 'user';
       	RunQuery(client, apiType, query)
		.then((response) => {
			refresh();
		}).catch(() => {
			setDataLoading(false);
		});

    }, [client, params.fullEntity, refresh, setDataLoading]);

    const deleteItem = useCallback((id: string) => {
		if(!params.fullEntity){
			return;
		}

		const query = LGQuery.deleteFrom(params.fullEntity.name)
							.where(GplFilter.field('id', '', id))
							.select('id');

		setDataLoading(true);

		const apiType = params.fullEntity.owner == EntityOwner.Admin ? 'admin' : 'user';
       	RunQuery(client, apiType, query)
		.then((response) => {
			refresh();
		}).catch(() => {
			setDataLoading(false);
		});

    }, [client, params.fullEntity, refresh, setDataLoading]);

    const deleteBulk = useCallback((ids: Array<string>) => {
        
    }, []);

    useEffect(() => {
        if(params.fullEntity?.loading ?? true){
            return;
        }

        if(!params.fullEntity){
            return;
        }

        if (params.fullEntity.loading) {
            return;
        }
        
        fetch();
    }, [fetch, dataQuery]);

    // const numberOfPages = useMemo(() => {
    //     if (!dataConnection) {
    //         return 1;
    //     }

    //     if (!dataConnection.totalCount) {
    //         return 1;
    //     }

    //     if (!rowsPerPage) {
    //         return 1;
    //     }

    //     return Math.ceil(dataConnection.totalCount / rowsPerPage);
    // }, [dataConnection, rowsPerPage]);

    return useMemo(() => ({
        state: {
            data: data,
            // schemaLoading: !fullEntity || fullEntity.loading,
            // dataLoading: !fullEntity || fullEntity.loading || dataLoading,
            // entityFields: fullEntity?.fields ?? [],
            // entityEdges: fullEntity?.edges ?? [],
            dataLoading,

            dataConnection,
            pagesCount: -1,
            totalItems: dataConnection?.totalCount ?? 0
        },
        refresh,
        addItem,
        addItems,
        updateItem,
        deleteItem,
        deleteBulk
    }), [data,  dataLoading, dataConnection?.totalCount, addItem, dataConnection, deleteBulk, deleteItem, updateItem]);
}
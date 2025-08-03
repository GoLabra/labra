import { ContentManagerSearchState } from "@/types/content-manager-search-state";
import { UsersStore } from "@/types/content-manager-store-state";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { FullEntity } from "@/types/entity";
import { addNotification } from "@/lib/notifications/store";
import { Edge, EntityOwner, Field } from "@/lib/apollo/graphql.entities";
import { getFiltersFromQuery } from "@/lib/utils/get-filters-from-query";
import { useApolloClient } from "@apollo/client";
import { RunQuery, useLgQuery } from "./use-lg-query";
import { FilterOperators, GplFilter, LGQuery } from "lg-query";

interface UseGridContentManagerStoreparams {
	fullEntity?: FullEntity;
	searchState: ContentManagerSearchState;
}
export const useContentManagerStore = (params: UseGridContentManagerStoreparams): UsersStore => {

    const { name:entityName } = params.fullEntity ?? {};

	const client = useApolloClient();

	const fields = useMemo(() => params.fullEntity?.fields.map(i => i.name), [params.fullEntity?.fields]);
	const edges = useMemo(() => params.fullEntity?.edges.filter(i => i.relationType !== 'ManyToMany')
													.filter(i => i.relationType !== 'OneToMany')
													.filter(i => i.relationType !== 'Many'), [params.fullEntity?.edges]);

	const dataQuery = useMemo(() => {

		if(params.fullEntity?.loading ?? true){
			return null;
		}

		if(!fields || !edges){
			return null;
		}

		let query = LGQuery.from<any>(params.fullEntity!.name)
							.select( ...fields)
		// add edges							
		query = edges.reduce((query: LGQuery<any>, edge) => {
			return query.include(edge.name, q => q.select('id', edge.relatedEntity.displayField.name));
		}, query);

		// add filters
		if(Object.entries(params.searchState.filter).length){
			query = query.where(
				GplFilter.and(
					...Object.entries(params.searchState.filter).map(([key, filter]) => { 
						const restul = GplFilter.field(key, filter.operator as FilterOperators, filter.value) as GplFilter<any>; 
						return restul;
					})
				)
			);
		}
		
		const queryFilter = getFiltersFromQuery(params.searchState.query, params.fullEntity?.fields ?? [])
		if(Object.entries(queryFilter).length){
			query = query.where(
				GplFilter.or(
					...Object.entries(queryFilter).map(([key, filter]) => { 
						return GplFilter.field(key, filter.operator as FilterOperators, filter.value) as GplFilter<any>; 
					})
				)
			);
		}

		// add skip
		query = query.skip((params.searchState.page - 1) * params.searchState.rowsPerPage);
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
			return null;
		}

		if(!dataQuery){
			return null;
		}

		let query = LGQuery.fromConnection<any>(entityName)
			     .select('totalCount');

		// add filters
		query = Object.entries(params.searchState.filter).reduce((query, [key, value]) => { 
			return query.where(GplFilter.field(key, value.operator as FilterOperators, value.value));
		}, query);

		return query;
	}, [entityName, dataQuery]);

	const apiType = params.fullEntity?.owner == EntityOwner.Admin ? 'admin' : 'user';

	const dataResponse = useLgQuery({
		apiType,
		query: useMemo(() => [dataQuery, connectionQuery], [dataQuery, connectionQuery]),
		skip: dataQuery == null
	});
    
    const addItem = useCallback((data: any) => {
		if(!params.fullEntity){
			return;
		}

		const query = LGQuery.create(params.fullEntity.name, data)
							.select('id');

		const apiType = params.fullEntity.owner == EntityOwner.Admin ? 'admin' : 'user';
       	RunQuery(client, apiType, query)
		.then((response) => {
			dataResponse.refetch();
		});
			
    }, [client, params.fullEntity, dataResponse.refetch]);


    const addItems = useCallback((data: any[]) => {
		if(!params.fullEntity){
			return;
		}

		const query = LGQuery.create(params.fullEntity.name, data)
							.select('id');

		const apiType = params.fullEntity.owner == EntityOwner.Admin ? 'admin' : 'user';
       	RunQuery(client, apiType, query)
		.then((response) => {
			dataResponse.refetch();
		});
    }, [client, params.fullEntity, dataResponse.refetch]);

    const updateItem = useCallback((id: string, data: any) => {

		if(!params.fullEntity){
			return;
		}

		const query = LGQuery.update(params.fullEntity.name, data)
							.where(GplFilter.field('id', '', id))
							.select('id');

		const apiType = params.fullEntity.owner == EntityOwner.Admin ? 'admin' : 'user';
       	RunQuery(client, apiType, query)
		.then((response) => {
			dataResponse.refetch();
		});

    }, [client, params.fullEntity, dataResponse.refetch]);

    const deleteItem = useCallback((id: string) => {
		if(!params.fullEntity){
			return;
		}

		const query = LGQuery.deleteFrom(params.fullEntity.name)
							.where(GplFilter.field('id', '', id))
							.select('id');

		const apiType = params.fullEntity.owner == EntityOwner.Admin ? 'admin' : 'user';
       	RunQuery(client, apiType, query)
		.then((response) => {
			dataResponse.refetch();
		});

    }, [client, params.fullEntity, dataResponse.refetch]);

    const deleteBulk = useCallback((ids: Array<string>) => {
        if(!params.fullEntity){
			return;
		}

		const query = LGQuery.deleteFrom(params.fullEntity.name)
							.where(GplFilter.field('id', 'In', ids));

		const apiType = params.fullEntity.owner == EntityOwner.Admin ? 'admin' : 'user';
       	RunQuery(client, apiType, query)
		.then((response) => {
			dataResponse.refetch();
		});

    }, [client, params.fullEntity, dataResponse.refetch]);

	const onlyData = useMemo(() => ({
		data: dataQuery?.getResultData(dataResponse.data) ?? [],
 		dataConnection: connectionQuery?.getResultData(dataResponse.data) ?? [],
	}), [dataResponse.data]);

	const numberOfPages = useMemo(() => {
        if (!onlyData.dataConnection) {
            return 1;
        }

        if (!onlyData.dataConnection.totalCount) {
            return 1;
        }

        if (!params.searchState.rowsPerPage) {
            return 1;
        }

        return Math.ceil(onlyData.dataConnection.totalCount / params.searchState.rowsPerPage);
    }, [onlyData.dataConnection, params.searchState.rowsPerPage]);

    return useMemo(() => ({
        state: {
            data: onlyData.data,
            dataLoading: dataResponse.loading,
            pagesCount: numberOfPages,
            totalItems: onlyData.dataConnection?.totalCount ?? 0
        },
        refresh: () => dataResponse.refetch(),
        addItem,
        addItems,
        updateItem,
        deleteItem,
        deleteBulk
    }), [dataResponse.refetch, dataResponse.loading, onlyData, addItem, deleteBulk, deleteItem, updateItem, numberOfPages]);
}
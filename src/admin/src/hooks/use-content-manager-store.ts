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
import { RunQuery, useLgQuery } from "./use-lg-query";

interface UseGridContentManagerStoreprops {
	fullEntity?: FullEntity;
	searchState: ContentManagerSearchState;
}
export const useContentManagerStore = (props: UseGridContentManagerStoreprops): UsersStore => {

    const { name:entityName } = props.fullEntity ?? {};

	const client = useApolloClient();

	const fields = useMemo(() => props.fullEntity?.fields.map(i => i.name), [props.fullEntity?.fields]);
	const edges = useMemo(() => props.fullEntity?.edges.filter(i => i.relationType !== 'ManyToMany')
													.filter(i => i.relationType !== 'ManyToOne')
													.filter(i => i.relationType !== 'Many'), [props.fullEntity?.edges]);

	const dataQuery = useMemo(() => {

		if(props.fullEntity?.loading ?? true){
			return null;
		}

		if(!fields || !edges){
			return null;
		}

		let query = LGQuery.from<any>(props.fullEntity!.name)
							.select( ...fields)
		// add edges							
		query = edges.reduce((query: LGQuery<any>, edge) => {
			return query.include(edge.name, q => q.select('id', edge.relatedEntity.displayField.name));
		}, query);

		// add filters
		query = Object.entries(props.searchState.filter).reduce((query, [key, value]) => { 
			return query.where(GplFilter.field(key, value.operator, value.value));
		}, query);

		// add skip
		query = query.skip((props.searchState.page - 1) * props.searchState.rowsPerPage);
		// add first
		query = query.first(props.searchState.rowsPerPage);
		// add order
		if(props.searchState.sortBy) {
			if(props.searchState.order == 'asc'){
				query = query.orderByAscending(props.searchState.sortBy);
			} else {
				query = query.orderByDescending(props.searchState.sortBy);
			}
		}
		
		return query;	
	}, [entityName, fields, edges, props.searchState]);

	const connectionQuery = useMemo(() => {
		if(!entityName){
			return null!;
		}

		let query = LGQuery.fromConnection<any>(entityName)
			     .select('totalCount');

		// add filters
		query = Object.entries(props.searchState.filter).reduce((query, [key, value]) => { 
			return query.where(GplFilter.field(key, value.operator, value.value));
		}, query);

		return query;
	}, [entityName]);

	const apiType = props.fullEntity?.owner == EntityOwner.Admin ? 'admin' : 'user';

	const dataResponse = useLgQuery({
		apiType,
		query: useMemo(() => [dataQuery, connectionQuery], [dataQuery, connectionQuery]),
		skip: dataQuery == null
	});
    
    const addItem = useCallback((data: any) => {
		if(!props.fullEntity){
			return;
		}

		const query = LGQuery.create(props.fullEntity.name, data)
							.select('id');

		const apiType = props.fullEntity.owner == EntityOwner.Admin ? 'admin' : 'user';
       	RunQuery(client, apiType, query)
		.then((response) => {
			dataResponse.refetch();
		});
			
    }, [client, props.fullEntity, dataResponse.refetch]);


    const addItems = useCallback((data: any[]) => {
		if(!props.fullEntity){
			return;
		}

		const query = LGQuery.create(props.fullEntity.name, data)
							.select('id');

		const apiType = props.fullEntity.owner == EntityOwner.Admin ? 'admin' : 'user';
       	RunQuery(client, apiType, query)
		.then((response) => {
			dataResponse.refetch();
		});
    }, [client, props.fullEntity, dataResponse.refetch]);

    const updateItem = useCallback((id: string, data: any) => {

		if(!props.fullEntity){
			return;
		}

		const query = LGQuery.update(props.fullEntity.name, data)
							.where(GplFilter.field('id', '', id))
							.select('id');

		const apiType = props.fullEntity.owner == EntityOwner.Admin ? 'admin' : 'user';
       	RunQuery(client, apiType, query)
		.then((response) => {
			dataResponse.refetch();
		});

    }, [client, props.fullEntity, dataResponse.refetch]);

    const deleteItem = useCallback((id: string) => {
		if(!props.fullEntity){
			return;
		}

		const query = LGQuery.deleteFrom(props.fullEntity.name)
							.where(GplFilter.field('id', '', id))
							.select('id');

		const apiType = props.fullEntity.owner == EntityOwner.Admin ? 'admin' : 'user';
       	RunQuery(client, apiType, query)
		.then((response) => {
			dataResponse.refetch();
		});

    }, [client, props.fullEntity, dataResponse.refetch]);

    const deleteBulk = useCallback((ids: Array<string>) => {
        if(!props.fullEntity){
			return;
		}

		const query = LGQuery.deleteFrom(props.fullEntity.name)
							.where(GplFilter.field('id', 'In', ids));

		const apiType = props.fullEntity.owner == EntityOwner.Admin ? 'admin' : 'user';
       	RunQuery(client, apiType, query)
		.then((response) => {
			dataResponse.refetch();
		});

    }, [client, props.fullEntity, dataResponse.refetch]);

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

        if (!props.searchState.rowsPerPage) {
            return 1;
        }

        return Math.ceil(onlyData.dataConnection.totalCount / props.searchState.rowsPerPage);
    }, [onlyData.dataConnection, props.searchState.rowsPerPage]);

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
import { ContentManagerSearchState } from "@/types/content-manager-search-state";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { FullEntity } from "@/types/entity";
import { Edge, EntityOwner, Field } from "@/lib/apollo/graphql.entities";
import { useApolloClient } from "@apollo/client";
import { RunQuery } from "./use-lg-query";
import { FilterOperators, GplFilter, LGQuery } from "lg-query";

interface UseLookupGridContentManagerStoreParams {
	fullEntity: FullEntity | null;
	searchState: ContentManagerSearchState;
}
export const useLookupContentManagerStore = (props: UseLookupGridContentManagerStoreParams) => {

    const { name:entityName } = props.fullEntity ?? {};

    //const contentManagerStoreRequest = useContentManagerStoreRequest({entityName, entityOwner: owner ?? EntityOwner.User});
	const client = useApolloClient();
    const [dataLoading, setDataLoading] = useState<boolean>(false);
    const [data, setData] = useState<[]>([]);
    const [dataConnection, setDataConnection] = useState<{ totalCount: number } | null>(null);

	const fields = useMemo(() => props.fullEntity?.fields.map(i => i.name), [props.fullEntity?.fields]);
	const edges = useMemo(() => props.fullEntity?.edges.filter(i => i.relationType !== 'ManyToMany')
													.filter(i => i.relationType !== 'ManyToOne')
													.filter(i => i.relationType !== 'Many'), [props.fullEntity?.edges]);

	const dataQuery = useMemo(() => {


		if(props.fullEntity?.loading ?? true) {
			return null;
		}

		if(!fields || !edges){
			return null;
		}

		let query = LGQuery.from<any>(props.fullEntity!.name)
							.select('id', props.fullEntity!.displayField!.name);

		// add filters
		query = Object.entries(props.searchState.filter).reduce((query, [key, value]) => { 
			return query.where(GplFilter.field(key, value.operator as FilterOperators, value.value)); 
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
			return query.where(GplFilter.field(key, value.operator as FilterOperators, value.value));
		}, query);

		return query;
	}, [entityName]);

    const fetch = useCallback(() => {

		if(!props.fullEntity){
			return;
		}

		if(!dataQuery){
			return;
		}
		
        setDataLoading(true);

		const apiType = props.fullEntity.owner == EntityOwner.Admin ? 'admin' : 'user';

		RunQuery(client, apiType, dataQuery, connectionQuery).then((response: any) => {
			setData(dataQuery.getResultData(response.data) as []); 
			setDataConnection(connectionQuery.getResultData(response.data));
		}).finally(() => {
			setDataLoading(false);
		});

    }, [client, props.fullEntity, dataQuery, connectionQuery, setData, setDataLoading]);

    const refresh = useCallback(() => {
        fetch();
    }, [fetch]);

    useEffect(() => {
        if(props.fullEntity?.loading ?? true){
            return;
        }

        if(!props.fullEntity){
            return;
        }

        if (props.fullEntity.loading) {
            return;
        }
        
        fetch();
    }, [fetch, dataQuery]);

    return useMemo(() => ({
        state: {
            data: data,
            dataLoading,
            dataConnection,
            pagesCount: -1,
            totalItems: dataConnection?.totalCount ?? 0
        },
        refresh,
    }), [data,  dataLoading, dataConnection?.totalCount, dataConnection]);
}
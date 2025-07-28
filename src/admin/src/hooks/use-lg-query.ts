import { ADMIN_CONTEXT } from "@/lib/apollo/apolloWrapper";
import { LGQuery } from "@/lib/apollo/builders/LabraGqlApiBuilder/LGQuery";
import { ILGQuery } from "@/lib/apollo/builders/LabraGqlApiBuilder/types/types";
import { ApolloClient, gql, useApolloClient, useQuery } from "@apollo/client";
import { useEffect, useMemo, useState } from "react";


interface UseLgQueryParams {
	apiType: ApiType;
	query: Array<ILGQuery | null>;
	skip?: boolean;
}
export const useLgQuery = (props: UseLgQueryParams) => {

	const gplQuery = useMemo(() => {
		const notNullQuery = props.query.filter(i => !!i);

		if(!notNullQuery.length){
			return null;
		}

		const gqlQuery = LGQuery.merge(...notNullQuery);

		return gqlQuery;
	}, [props.query]);


	const context = props.apiType == 'admin' ? ADMIN_CONTEXT : {};

	 // Use the hook
	const data = useQuery( 
		gql(gplQuery?.query ?? `query { __typename }`), 
		{
			notifyOnNetworkStatusChange: true,
			variables: gplQuery?.variables,
			fetchPolicy: "network-only",
			skip: gplQuery == null,
			context
		}
	);

	return data;
}



export type ApiType = 'admin' | 'user';

export const RunQuery = (client: ApolloClient<object>, apiType: ApiType, ...query: Array<ILGQuery | null>):Promise<any> => { 
	const promise = new Promise((resolve, reject) => {

		const notNullQuery = query.filter(i => !!i);

		if(!notNullQuery.length){
			reject("No query to run");
		}

		const gqlQuery = LGQuery.merge(...notNullQuery);

		const context = apiType == 'admin' ? ADMIN_CONTEXT : {};
		
		if(notNullQuery[0].isMutation == false){
			client.query({ query: gql(gqlQuery.query), variables: gqlQuery.variables, fetchPolicy: "network-only", context })
				.then((response) => {
					resolve(response);
				}).catch((error) => {
					reject(error);
				});
		} else {
			client.mutate({ mutation: gql(gqlQuery.query), variables: gqlQuery.variables, fetchPolicy: "network-only", context }) 
				.then((response) => {
					resolve(response);
				}).catch((error) => {
					reject(error);
				});
		}
	});

	return promise;
}
import { ApolloClient, gql, useApolloClient } from "@apollo/client";
import { ADMIN_CONTEXT } from "@/lib/apollo/apolloWrapper";
import { FullEntity } from "@/types/entity";
import { LGQuery } from "@/lib/apollo/builders/LabraGqlApiBuilder/LGQuery";
import { ILGQuery } from "@/lib/apollo/builders/LabraGqlApiBuilder/types/types";



// interface UseContentManagerStoreRequestParams {
// 	entityOwner: EntityOwner;
//     entityName: string
// }
// export const useContentManagerStoreRequest = (props: UseContentManagerStoreRequestParams) => {

//     const { entityName } = props;
//     const client = useApolloClient();

// 	const apiContext = useMemo(() => props.entityOwner == EntityOwner.Admin ? ADMIN_CONTEXT : {}, [props.entityOwner]);

//     const fetchData = useCallback((query: LGQuery<any>) => {

//         const promise = new Promise((resolve, reject) => {

// 			// const gqlQuery = query.getGraphQlQuery();
//             // client.query({ query: gql(gqlQuery.query), variables: gqlQuery.variables, fetchPolicy: "network-only", context: apiContext })
// 			// .then((response) => {
//             //     // resolve({
//             //     //     data: response.data[operationName],
//             //     //     connection: response.data[builder.getOperationConnectionName()]
//             //     // });
//             // }).catch((error) => {
//             //     reject(error);
//             // });
//         });

//         return promise;
//     }, [client, entityName, apiContext]);

//     const addItem = useCallback((entityName: string, data: any) => {

//         var builder = new GqlDataCREATEMutationBuilder().addField('id');
//         var mutation = builder.addEntityName(entityName)
//             .build(data);

//         if (!mutation) {
//             return Promise.reject("No mutation to add");
//         }

//         const operationName = builder.getOperationName();
//         if (!operationName) {
//             return Promise.reject("No operationName to add");
//         }
//         return client.mutate({ mutation: gql(mutation.query), variables: mutation.variables, context: apiContext });
//     }, [client, apiContext]);

//     const updateItem = useCallback((entityName: string, id: string, data: any) => {

// 		debugger;

//         var builder = new GqlDataUPDATEMutationBuilder().addField('id');
//         var mutation = builder.addEntityName(entityName)
//             .build({ id, data });

//         if (!mutation) {
//             return;
//         }

//         const operationName = builder.getOperationName();
//         if (!operationName) {
//             return;
//         }
//         return client.mutate({ mutation: gql(mutation.query), variables: mutation.variables, context: apiContext });
//     }, [client, apiContext]);

//     const deleteItem = useCallback((entityName: string, id: string) => {
//         var builder = new GqlDataDELETEMutationBuilder().addField('id');
//         var mutation = builder.addEntityName(entityName)
//             .build(id);

//         if (!mutation) {
//             return;
//         }

//         const operationName = builder.getOperationName();
//         if (!operationName) {
//             return;
//         }
//         return client.mutate({ mutation: gql(mutation.query), variables: mutation.variables, context: apiContext });
//     }, [client, apiContext]);

//     const deleteBulk = useCallback((entityName: string, ids: Array<string>) => {
//         var builder = new GqlDataDELETEBulkMutationBuilder();
//         var mutation = builder.addEntityName(entityName)
//             .build(ids);

//         if (!mutation) {
//             return;
//         }

//         const operationName = builder.getOperationName();
//         if (!operationName) {
//             return;
//         }
//         return client.mutate({ mutation: gql(mutation.query), variables: mutation.variables, context: apiContext });
//     }, [client, apiContext]);

//     return useMemo(() => ({
//         fetchData,
//         addItem,
//         updateItem,
//         deleteItem,
//         deleteBulk
//     }), [fetchData, addItem, updateItem, deleteItem, deleteBulk]);
// }
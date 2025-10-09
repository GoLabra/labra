// import { useEntities } from "@/hooks/use-entities";
// import { Edge, Field, useGetEntityFirstLevelSchemaQuery } from "@/lib/apollo/graphql";
// import { EntityStoreSchema } from "@/types/content-manager-store-entity-schema";
// import { useMemo } from "react";

// interface useContentManagerStoreSchemaProps {
//     entityName: string;
// }
// export const useContentManagerStoreSchemaFirstLevel = (props: useContentManagerStoreSchemaProps): EntityStoreSchema => {
//     const { entityName } = props;
//     const { data: schema, error, loading: schemaLoading, refetch } = useGetEntityFirstLevelSchemaQuery({ variables: { name: entityName } });

//     return useMemo(() => ({
//         entityName,
//         displayField: schema?.entity?.displayField,
//         fields: (schema?.entity?.fields?.filter(i => i.name != 'createdAt' && i.name != 'updatedAt') ?? []) as Field[], //TODO: skip createdAt, updateAt from form schema ?? [],
//         edges: (schema?.entity?.edges?.filter(i => i.name.startsWith('ref') == false) ?? []) as Edge[], //TODO
//         isError: !!error || !schema?.entity,
//         loading: schemaLoading,
//     }), [entityName, error, schema?.entity, schemaLoading]);
// }

// export const useContentManagerStoreSchemaDisplayField = (props: useContentManagerStoreSchemaProps): EntityStoreSchema => {

//     const { entityName } = props;
//     const graphEntities = useEntities();
//     const entity = useMemo(() => graphEntities?.entities.find(i => i.name == entityName), [entityName, graphEntities?.entities]);

//     return useMemo(() => ({
//         entityName: entity?.name ?? '',
//         displayField: entity?.displayField,
//         fields: [{
//             caption: "Id",
//             name: "id",
//             type: "ID",
//             __typename: "Field"
//         }, 
//         entity?.displayField] as Field[],
//         edges: [] as Edge[],
//         isError: !!graphEntities?.error,
//         loading: graphEntities?.childrenLoading,
//     }), [entity?.name, entity?.displayField, graphEntities?.childrenLoading, graphEntities?.error]);
// }
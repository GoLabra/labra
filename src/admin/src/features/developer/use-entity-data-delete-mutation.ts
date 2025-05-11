import { useFullEntity } from "@/hooks/use-entities";
import { GqlDataDELETEMutationBuilder } from "@/lib/apollo/builders/gqlMutationBuilder";
import { useMemo } from "react";


export const useEntityDataDeleteMutation = (entityName: string) => {

    var query = useMemo(() => {
        return getEntityDataDeleteMutationQuery(entityName);
    }, [entityName]);

    return query;
}

export const getEntityDataDeleteMutationQuery = (entityName: string) => {
    return new GqlDataDELETEMutationBuilder()
        .addEntityName(entityName)
        .addField('id')
        .build('[id]')
}
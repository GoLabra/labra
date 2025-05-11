import { useFullEntity } from "@/hooks/use-entities";
import {  GqlDataUPDATEMutationBuilder } from "@/lib/apollo/builders/gqlMutationBuilder";
import { useMemo } from "react";

export const useEntityDataUpdateMutation = (entityName: string) => {

    // const fullEntity = useFullEntity({ entityName });

    // if (!fullEntity) {
    //     return null;
    // }

    var query = useMemo(() => {
        return getEntityDataUpdateMutationQuery(entityName);
    }, [entityName]);

    return query;
}

export const getEntityDataUpdateMutationQuery = (entityName: string) => {
    return new GqlDataUPDATEMutationBuilder()
        .addEntityName(entityName)
        .addField('id')
        .build({
            id: '[id]',
            data: {}
        });
}
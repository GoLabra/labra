import { useFullEntity } from "@/hooks/use-entities";
import { GqlDataCREATEMutationBuilder, GqlDataDELETEMutationBuilder } from "@/lib/apollo/builders/gqlMutationBuilder";
import { useMemo } from "react";

export const useEntityDataNewMutation = (entityName: string) => {

    // const fullEntity = useFullEntity({ entityName });

    // if (!fullEntity) {
    //     return null;
    // }

    var query = useMemo(() => {
        return getEntityDataNewMutationQuery(entityName);
    }, [entityName]);

    return query;
}

export const getEntityDataNewMutationQuery = (entityName: string) => {
    return new GqlDataCREATEMutationBuilder()
        .addEntityName(entityName)
        .addField('id')
        .build({

        });
}